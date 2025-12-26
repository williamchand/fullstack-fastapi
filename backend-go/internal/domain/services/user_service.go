package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	cfg              *config.Config
	userRepo         repositories.UserRepository
	oauthRepo        repositories.OAuthRepository
	txManager        repositories.TransactionManager
	jwtRepo          repositories.JWTRepository
	emailTplRepo     repositories.EmailTemplateRepository
	verificationRepo repositories.VerificationCodeRepository
	smtpSender       repositories.Sender
	wahaClient       repositories.WahaClient
}

func NewUserService(
	cfg *config.Config,
	userRepo repositories.UserRepository,
	oauthRepo repositories.OAuthRepository,
	txManager repositories.TransactionManager,
	jwtRepo repositories.JWTRepository,
	emailTplRepo repositories.EmailTemplateRepository,
	verificationRepo repositories.VerificationCodeRepository,
	smtpSender repositories.Sender,
	wahaClient repositories.WahaClient,
) *UserService {
	return &UserService{
		cfg:              cfg,
		userRepo:         userRepo,
		oauthRepo:        oauthRepo,
		txManager:        txManager,
		jwtRepo:          jwtRepo,
		emailTplRepo:     emailTplRepo,
		verificationRepo: verificationRepo,
		smtpSender:       smtpSender,
		wahaClient:       wahaClient,
	}
}

var ErrInvalidState = fmt.Errorf("invalid state")
var ErrInvalidToken = fmt.Errorf("invalid token")
var ErrWeakPassword = fmt.Errorf("weak password")
var ErrInvalidPreviousPassword = fmt.Errorf("invalid previous password")
var ErrUnauthorized = fmt.Errorf("unauthorized")

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) ListUsers(ctx context.Context, offset, limit int32) ([]*entities.User, int, error) {
	return s.userRepo.ListUsers(ctx, offset, limit)
}

func (s *UserService) CreateUser(ctx context.Context, email, password, fullName string, roles []entities.RoleEnum, isActive bool) (*entities.User, error) {
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
		IsActive:        isActive,
		IsEmailVerified: false,
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

func (s *UserService) UpdateProfile(ctx context.Context, id string, fullName *string, password *string, previousPassword *string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || existingUser == nil {
		return nil, ErrUserNotFound
	}
	var hashed *string
	if password != nil && *password != "" {
		// require and validate previous password before changing, unless superuser
		if !util.HasRole(existingUser, string(entities.RoleSuperuser)) {
			if previousPassword == nil || *previousPassword == "" {
				return nil, ErrInvalidPreviousPassword
			}
			if _, err := s.ValidatePassword(ctx, existingUser.Email, *previousPassword); err != nil {
				return nil, ErrInvalidPreviousPassword
			}
		}

		hp, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hs := string(hp)
		hashed = &hs
	}
	u := &entities.User{
		ID:             existingUser.ID,
		FullName:       fullName,
		HashedPassword: hashed,
	}
	updated, err := s.userRepo.UpdateProfile(ctx, u.ID, u.FullName, u.HashedPassword)
	if err != nil {
		return nil, err
	}
	updated.Roles = existingUser.Roles
	return updated, nil
}

func (s *UserService) AdminUpdateUser(ctx context.Context, adminID string, targetUserID string, fullName *string, password *string, roles []string, isActive *bool) (*entities.User, error) {
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	admin, err := s.userRepo.GetByID(ctx, adminUUID)
	if err != nil || admin == nil {
		return nil, ErrUserNotFound
	}
	if !util.HasRole(admin, string(entities.RoleSuperuser)) {
		return nil, ErrUnauthorized
	}
	targetUUID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	target, err := s.userRepo.GetByID(ctx, targetUUID)
	if err != nil || target == nil {
		return nil, ErrUserNotFound
	}
	var hashed *string
	if password != nil && *password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hs := string(hp)
		hashed = &hs
	}
	newRoles := target.Roles
	if len(roles) > 0 {
		newRoles = roles
	}
	newActive := target.IsActive
	if isActive != nil {
		newActive = *isActive
	}
	u := &entities.User{
		ID:             target.ID,
		FullName:       fullName,
		HashedPassword: hashed,
		Roles:          newRoles,
		IsActive:       newActive,
	}
	updated, err := s.userRepo.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *UserService) AddPhoneNumber(ctx context.Context, id string, phone string, region string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return ErrUserNotFound
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	if !user.IsEmailVerified {
		return ErrInvalidState
	}
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ErrInvalidState
	}
	normalized, ok := util.NormalizeE164(phone, region)
	if !ok {
		return ErrInvalidState
	}
	if existing, _ := s.userRepo.GetByPhone(ctx, normalized); existing != nil && existing.ID != user.ID {
		return ErrUserExists
	}
	code, err := generateOTP(6)
	if err != nil {
		return err
	}
	v := &entities.VerificationCode{
		UserID:        &user.ID,
		Code:          code,
		Type:          entities.VerificationTypePhone,
		ExpiresAt:     time.Now().Add(10 * time.Minute),
		ExtraMetadata: map[string]any{"purpose": entities.VerificationPurposeAddPhone, "new_phone": normalized},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return err
	}
	if s.wahaClient != nil {
		tpl, tplErr := s.emailTplRepo.GetByName(ctx, entities.EmailTemplateVerificationPhone)
		if tplErr == nil {
			body, _ := util.FillTextTemplate(tpl.Body, map[string]string{"code": code})
			go func() {
				if err := s.wahaClient.SendText(context.Background(), normalized, body); err != nil {
					log.Println(fmt.Errorf("failed to send WhatsApp OTP: %w", err))
				}
			}()
		}
	}
	return nil
}

func (s *UserService) VerifyAddPhone(ctx context.Context, id string, code string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return ErrUserNotFound
	}
	user, err := s.userRepo.GetByID(ctx, userID)
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
	newPhone, _ := v.ExtraMetadata["new_phone"].(string)
	if newPhone == "" {
		return ErrInvalidState
	}
	if existing, _ := s.userRepo.GetByPhone(ctx, newPhone); existing != nil && existing.ID != user.ID {
		return ErrUserExists
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return err
	}
	if _, err := s.userRepo.UpdatePhone(ctx, user.ID, newPhone); err != nil {
		return err
	}
	return s.userRepo.SetPhoneVerified(ctx, user.ID)
}

func (s *UserService) AddEmail(ctx context.Context, id string, email string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return ErrUserNotFound
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	if !user.IsPhoneVerified {
		return ErrInvalidState
	}
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return ErrInvalidState
	}
	if existing, _ := s.userRepo.GetByEmail(ctx, email); existing != nil && existing.ID != user.ID {
		return ErrUserExists
	}
	code, err := generateOTP(6)
	if err != nil {
		return err
	}
	v := &entities.VerificationCode{
		UserID:        &user.ID,
		Code:          code,
		Type:          entities.VerificationTypeEmail,
		ExpiresAt:     time.Now().Add(15 * time.Minute),
		ExtraMetadata: map[string]any{"purpose": entities.VerificationPurposeAddEmail, "new_email": email},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return err
	}
	if s.smtpSender != nil {
		tpl, tplErr := s.emailTplRepo.GetByName(ctx, entities.EmailTemplateVerificationEmail)
		if tplErr == nil {
			body, _ := util.FillTextTemplate(tpl.Body, map[string]string{"code": code})
			go func() { _ = s.smtpSender.Send(entities.Message{To: email, Subject: tpl.Subject, Body: body}) }()
		}
	}
	return nil
}

func (s *UserService) VerifyAddEmail(ctx context.Context, id string, code string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return ErrUserNotFound
	}
	user, err := s.userRepo.GetByID(ctx, userID)
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
	newEmail, _ := v.ExtraMetadata["new_email"].(string)
	if newEmail == "" {
		return ErrInvalidState
	}
	newEmail = strings.TrimSpace(strings.ToLower(newEmail))
	if existing, _ := s.userRepo.GetByEmail(ctx, newEmail); existing != nil && existing.ID != user.ID {
		return ErrUserExists
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return err
	}
	if _, err := s.userRepo.UpdateEmail(ctx, user.ID, newEmail); err != nil {
		return err
	}
	return s.userRepo.SetEmailVerified(ctx, user.ID)
}

func (s *UserService) ValidatePassword(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
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
		UserID:        &user.ID,
		Code:          code,
		Type:          entities.VerificationTypeEmail,
		ExpiresAt:     expires,
		ExtraMetadata: map[string]any{"purpose": entities.VerificationPurposeEmailVerification},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	tpl, err := s.emailTplRepo.GetByName(ctx, entities.EmailTemplateVerificationEmail)
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
	go func() {
		if err := s.smtpSender.Send(msg); err != nil {
			log.Println(fmt.Errorf("failed to send email: %w", err))
		}
	}()
	return nil
}

// RequestPasswordReset creates a password reset token, stores it, and sends email
func (s *UserService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// Generate hash token
	token := util.GenerateSecureToken(32)

	expires := time.Now().Add(60 * time.Minute)
	v := &entities.VerificationCode{
		UserID:        &user.ID,
		Code:          token,
		Type:          entities.VerificationTypePasswordReset,
		ExpiresAt:     expires,
		ExtraMetadata: map[string]any{"purpose": entities.VerificationPurposePasswordReset},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	tpl, err := s.emailTplRepo.GetByName(ctx, entities.EmailTemplatePasswordReset)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}
	link := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.BaseURL, token)
	body, err := util.FillTextTemplate(tpl.Body, map[string]string{"link": link})
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}
	msg := entities.Message{To: user.Email, Subject: tpl.Subject, Body: body}
	go func() {
		if err := s.smtpSender.Send(msg); err != nil {
			log.Println(fmt.Errorf("failed to send email: %w", err))
		}
	}()
	return nil
}

// ResetPassword validates token and updates user's password
func (s *UserService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}
	v, err := s.verificationRepo.GetByCodeOnly(ctx, entities.VerificationTypePasswordReset, token)
	if err != nil || v == nil {
		return ErrInvalidOrExpiredCode
	}
	if v.UsedAt != nil || time.Now().After(v.ExpiresAt) {
		return ErrInvalidOrExpiredCode
	}
	user, err := s.userRepo.GetByID(ctx, *v.UserID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	hp, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	hs := string(hp)
	if _, err := s.userRepo.UpdateProfile(ctx, user.ID, nil, &hs); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if err := s.verificationRepo.MarkUsed(ctx, v.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}
	return nil
}

// RequestPhoneOTP generates and stores OTP for a user with given phone number
func (s *UserService) RequestPhoneOTP(ctx context.Context, phone string, region string) error {
	normalized, ok := util.NormalizeE164(phone, region)
	if !ok {
		return ErrInvalidState
	}
	user, err := s.userRepo.GetByPhone(ctx, normalized)
	if err != nil || user == nil {
		return ErrUserNotFound
	}
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}
	expires := time.Now().Add(10 * time.Minute)
	v := &entities.VerificationCode{
		UserID:        &user.ID,
		Code:          code,
		Type:          entities.VerificationTypePhone,
		ExpiresAt:     expires,
		ExtraMetadata: map[string]any{"purpose": entities.VerificationPurposePhoneOTP},
	}
	if err := s.verificationRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}
	// Send OTP via WhatsApp using WAHA client
	if s.wahaClient != nil {
		tpl, tplErr := s.emailTplRepo.GetByName(ctx, entities.EmailTemplateVerificationPhone)
		if tplErr == nil {
			body, _ := util.FillTextTemplate(tpl.Body, map[string]string{"code": code})
			go func() {
				if err := s.wahaClient.SendText(context.Background(), normalized, body); err != nil {
					log.Println(fmt.Errorf("failed to send WhatsApp OTP: %w", err))
				}
			}()
		}
	}
	return nil
}

var ErrInvalidOrExpiredCode = fmt.Errorf("invalid or expired code")

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
func (s *UserService) LoginWithPhone(ctx context.Context, phone, code string, region string) (*entities.TokenPair, error) {
	normalized, ok := util.NormalizeE164(phone, region)
	if !ok {
		return nil, ErrInvalidState
	}
	user, err := s.userRepo.GetByPhone(ctx, normalized)
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
	// Mark phone as verified after successful OTP validation
	if err := s.userRepo.SetPhoneVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to set phone verified: %w", err)
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

// RegisterPhoneUser creates a verification code for phone registration
func (s *UserService) RegisterPhoneUser(ctx context.Context, phone, fullName, region string) (string, error) {
	normalized, ok := util.NormalizeE164(phone, region)
	if !ok {
		return "", ErrInvalidState
	}
	existing, err := s.userRepo.GetByPhone(ctx, normalized)
	if err == nil && existing != nil {
		return "", ErrUserExists
	}

	// Generate hash token
	token := util.GenerateSecureToken(32)

	extraMetadata := map[string]any{
		"phone":    normalized,
		"fullName": fullName,
		"region":   region,
	}

	_, err = s.verificationRepo.CreateNoUser(ctx, token, entities.VerificationTypePhoneRegistration, extraMetadata, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	// Send OTP for verification
	_ = s.RequestPhoneOTP(ctx, normalized, region)

	return token, nil
}

// VerifyRegisterPhoneUser verifies the token and OTP, then creates the user
func (s *UserService) VerifyRegisterPhoneUser(ctx context.Context, token, otpCode string) (*entities.TokenPair, error) {
	vc, err := s.verificationRepo.GetByCodeOnly(ctx, entities.VerificationTypePhoneRegistration, token)
	if err != nil {
		return nil, ErrInvalidOrExpiredCode
	}
	if vc.UsedAt != nil || time.Now().After(vc.ExpiresAt) {
		return nil, ErrInvalidOrExpiredCode
	}

	phone, ok := vc.ExtraMetadata["phone"].(string)
	if !ok {
		return nil, ErrInvalidState
	}

	// Validate OTP
	vcOTP, err := s.verificationRepo.GetByCodeOnly(ctx, entities.VerificationTypePhone, otpCode)
	if err != nil {
		return nil, ErrInvalidOrExpiredCode
	}
	if vcOTP.UsedAt != nil || time.Now().After(vcOTP.ExpiresAt) {
		return nil, ErrInvalidOrExpiredCode
	}
	phoneFromOTP, ok := vcOTP.ExtraMetadata["phone"].(string)
	if !ok || phoneFromOTP != phone {
		return nil, ErrInvalidOrExpiredCode
	}

	fullName, ok := vc.ExtraMetadata["fullName"].(string)
	if !ok {
		return nil, ErrInvalidState
	}

	// Check if phone already exists
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
		IsPhoneVerified: true, // Since verified via OTP
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)
		verificationRepoTx := s.verificationRepo.WithTx(tx)
		var errTx error
		user, errTx = userRepoTx.Create(ctx, user)
		if errTx != nil {
			return errTx
		}
		errTx = userRepoTx.SetUserRoles(ctx, user.ID, []entities.RoleEnum{entities.RoleCustomer})
		if errTx != nil {
			return errTx
		}
		errTx = verificationRepoTx.MarkUsed(ctx, vc.ID)
		if errTx != nil {
			return errTx
		}
		errTx = verificationRepoTx.MarkUsed(ctx, vcOTP.ID)
		if errTx != nil {
			return errTx
		}
		return nil
	})
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
