package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo         repositories.UserRepository
	oauthRepo        repositories.OAuthRepository
	txManager        repositories.TransactionManager
	jwtRepo          repositories.JWTRepository
	emailTplRepo     repositories.EmailTemplateRepository
	verificationRepo repositories.VerificationCodeRepository
	smtpSender       repositories.Sender
}

func NewUserService(
	userRepo repositories.UserRepository,
	oauthRepo repositories.OAuthRepository,
	txManager repositories.TransactionManager,
	jwtRepo repositories.JWTRepository,
	emailTplRepo repositories.EmailTemplateRepository,
	verificationRepo repositories.VerificationCodeRepository,
	smtpSender repositories.Sender,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		oauthRepo:        oauthRepo,
		txManager:        txManager,
		jwtRepo:          jwtRepo,
		emailTplRepo:     emailTplRepo,
		verificationRepo: verificationRepo,
		smtpSender:       smtpSender,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) CreateUser(ctx context.Context, email, password, fullName, phoneNumber string, roles []entities.RoleEnum, isEmailVerified bool) (*entities.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	user := &entities.User{
		Email:           email,
		HashedPassword:  &hashedPasswordStr,
		FullName:        &fullName,
		IsActive:        true,
		IsEmailVerified: isEmailVerified,
	}

	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)

		user, err = userRepoTx.Create(ctx, user)
		if err != nil {
			return err
		}

		err = userRepoTx.SetUserRoles(ctx, user.ID, roles)
		if err != nil {
			return err
		}
		return nil
	})

	err = s.SendEmailVerification(ctx, email)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		user.Roles = append(user.Roles, string(role))
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, email, password, fullName, phoneNumber string, roles []entities.RoleEnum) (*entities.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser == nil {
		return nil, ErrUserNotFound
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:          existingUser.ID,
		Email:       email,
		FullName:    &fullName,
		PhoneNumber: &phoneNumber,
		IsActive:    true,
	}
	hashedPasswordStr := string(hashedPassword)
	if password != "" {
		user.HashedPassword = &hashedPasswordStr
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)

		user, err = userRepoTx.Update(ctx, user)
		if err != nil {
			return err
		}

		err = userRepoTx.SetUserRoles(ctx, user.ID, roles)
		if err != nil {
			return err
		}
		return nil
	})

	for _, role := range roles {
		user.Roles = append(user.Roles, string(role))
	}

	return user, nil
}

func (s *UserService) ValidatePassword(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}
	if !user.IsEmailVerified {
		return nil, ErrInvalidEmailNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context,
	username string,
	password string,
) (*entities.TokenPair, error) {
	user, err := s.ValidatePassword(ctx, username, password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtRepo.GenerateToken(user.ID, user.Email, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtRepo.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return &entities.TokenPair{
		User:             user,
		AccessToken:      accessToken.Token,
		RefreshToken:     refreshToken.Token,
		ExpiresAt:        accessToken.ExpiresAt,
		RefreshExpiresAt: refreshToken.ExpiresAt,
		IsNewUser:        false,
	}, nil
}

func (s *UserService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*entities.TokenResult, error) {
	newRefreshToken, err := s.jwtRepo.RefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	return newRefreshToken, nil
}

// Generate a numeric OTP code of given length
func generateOTP(length int) (string, error) {
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte(byte('0' + n.Int64()))
	}
	return b.String(), nil
}

// SendEmailVerification generates a verification code and sends an email using template
func (s *UserService) SendEmailVerification(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	expires := time.Now().Add(15 * time.Minute)
	v := &entities.VerificationCode{
		UserID:        user.ID,
		Code:          code,
		Type:          entities.VerificationTypeEmail,
		ExpiresAt:     expires,
		ExtraMetadata: map[string]any{"purpose": "email_verification"},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	tpl, err := s.emailTplRepo.GetByName(ctx, "verification_email")
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}
	fieldMap := map[string]string{
		"code": code,
	}
	body, err := util.FillTextTemplate(tpl.Body, fieldMap)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}
	msg := entities.Message{To: user.Email, Subject: tpl.Subject, Body: body}
	if err := s.smtpSender.Send(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// RequestPhoneOTP generates and stores OTP for a user with given phone number
func (s *UserService) RequestPhoneOTP(ctx context.Context, phone string) error {
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}
	expires := time.Now().Add(10 * time.Minute)
	v := &entities.VerificationCode{
		UserID:        user.ID,
		Code:          code,
		Type:          entities.VerificationTypePhone,
		ExpiresAt:     expires,
		ExtraMetadata: map[string]any{"purpose": "phone_otp"},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}
	// Optional: send notification via email as fallback if email exists
	if user.Email != "" && s.smtpSender != nil {
		tpl, tplErr := s.emailTplRepo.GetByName(ctx, "verification_email")
		if tplErr == nil {
			body := strings.ReplaceAll(tpl.Body, "{{code}}", code)
			_ = s.smtpSender.Send(entities.Message{To: user.Email, Subject: tpl.Subject, Body: body})
		}
	}
	return nil
}

var ErrInvalidOrExpiredCode = fmt.Errorf("invalid or expired code")

// VerifyPhoneOTP verifies the OTP and marks phone as verified
func (s *UserService) VerifyPhoneOTP(ctx context.Context, phone, code string) error {
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	v, err := s.verificationRepo.GetByCode(ctx, user.ID, entities.VerificationTypePhone, code)
	if err != nil || v == nil {
		return ErrInvalidOrExpiredCode
	}
	if v.UsedAt != nil || time.Now().After(v.ExpiresAt) {
		return ErrInvalidOrExpiredCode
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return fmt.Errorf("failed to mark code used: %w", err)
	}
	if err := s.userRepo.SetPhoneVerified(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to set phone verified: %w", err)
	}
	return nil
}

// VerifyEmailOTP verifies the OTP and marks email as verified
func (s *UserService) VerifyEmailOTP(ctx context.Context, email, code string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	v, err := s.verificationRepo.GetByCode(ctx, user.ID, entities.VerificationTypeEmail, code)
	if err != nil || v == nil {
		return ErrInvalidOrExpiredCode
	}
	if v.UsedAt != nil || time.Now().After(v.ExpiresAt) {
		return ErrInvalidOrExpiredCode
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return fmt.Errorf("failed to mark code used: %w", err)
	}
	if err := s.userRepo.SetEmailVerified(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}
	return nil
}

// LoginWithPhone verifies OTP and returns token pair
func (s *UserService) LoginWithPhone(ctx context.Context, phone, code string) (*entities.TokenPair, error) {
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	v, err := s.verificationRepo.GetByCode(ctx, user.ID, entities.VerificationTypePhone, code)
	if err != nil || v == nil {
		return nil, ErrInvalidOrExpiredCode
	}
	if v.UsedAt != nil || time.Now().After(v.ExpiresAt) {
		return nil, ErrInvalidOrExpiredCode
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return nil, fmt.Errorf("failed to mark code used: %w", err)
	}
	accessToken, err := s.jwtRepo.GenerateToken(user.ID, user.Email, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwtRepo.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return &entities.TokenPair{
		User:             user,
		AccessToken:      accessToken.Token,
		RefreshToken:     refreshToken.Token,
		ExpiresAt:        accessToken.ExpiresAt,
		RefreshExpiresAt: refreshToken.ExpiresAt,
		IsNewUser:        false,
	}, nil
}

// RegisterPhoneUser creates a user by phone and requests OTP
func (s *UserService) RegisterPhoneUser(ctx context.Context, phone, fullName string) (*entities.User, error) {
	existing, err := s.userRepo.GetByPhone(ctx, phone)
	if err == nil && existing != nil {
		return nil, ErrUserExists
	}

	user := &entities.User{
		Email:           "",
		PhoneNumber:     &phone,
		FullName:        &fullName,
		IsActive:        true,
		IsEmailVerified: false,
		IsPhoneVerified: false,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)
		var errTx error
		user, errTx = userRepoTx.Create(ctx, user)
		if errTx != nil {
			return errTx
		}
		errTx = userRepoTx.SetUserRoles(ctx, user.ID, []entities.RoleEnum{entities.RoleCustomer})
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Immediately issue OTP for phone verification
	_ = s.RequestPhoneOTP(ctx, phone)
	return user, nil
}
