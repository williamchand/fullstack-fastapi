package grpc

import (
	"context"
	"errors"

	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	salonappv1.UnimplementedUserServiceServer
	userService *services.UserService
}

func NewUserServer(userService *services.UserService) salonappv1.UserServiceServer {
	return &userServer{
		userService: userService,
	}
}

func (s *userServer) GetUser(ctx context.Context, req *emptypb.Empty) (*salonappv1.GetUserResponse, error) {
	user := util.UserFromContext(ctx)
	user, err := s.userService.GetUserByID(ctx, user.ID.String())
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &salonappv1.GetUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *salonappv1.ListUsersRequest) (*salonappv1.ListUsersResponse, error) {
	// Check if user is superuser
	user := util.UserFromContext(ctx)
	if !util.HasRole(user, string(entities.RoleSuperuser)) {
		return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
	}

	users, total, err := s.userService.ListUsers(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	protoUsers := make([]*salonappv1.User, len(users))
	for i, u := range users {
		protoUsers[i] = s.userToProto(u)
	}

	return &salonappv1.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
	}, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *salonappv1.CreateUserRequest) (*salonappv1.CreateUserResponse, error) {
	user := util.UserFromContext(ctx)
	isSuperuser := util.HasRole(user, string(entities.RoleSuperuser))

	roles := []entities.RoleEnum{}
	if isSuperuser && len(req.Roles) > 0 {
		for _, r := range req.Roles {
			roles = append(roles, entities.RoleEnum(r))
		}
	} else {
		roles = []entities.RoleEnum{entities.RoleCustomer}
	}

	isActive := false
	if isSuperuser {
		isActive = req.IsActive
	}

	userEntity, err := s.userService.CreateUser(ctx, req.Email, req.Password, req.FullName, roles, isActive)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &salonappv1.CreateUserResponse{
		User: s.userToProto(userEntity),
	}, nil
}

func (s *userServer) RefreshToken(ctx context.Context, req *salonappv1.RefreshTokenRequest) (*salonappv1.RefreshTokenResponse, error) {
	tokenPair, err := s.userService.RefreshToken(
		ctx,
		req.RefreshToken,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidRefreshToken):
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		default:
			return nil, status.Error(codes.Internal, "failed to refresh token")
		}
	}
	return &salonappv1.RefreshTokenResponse{
		AccessToken: tokenPair.Token,
		ExpiresAt:   timestamppb.New(tokenPair.ExpiresAt),
	}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *salonappv1.UpdateUserRequest) (*salonappv1.UpdateUserResponse, error) {
	user := util.UserFromContext(ctx)

	user, err := s.userService.UpdateProfile(ctx, user.ID.String(), req.FullName, req.Password, req.PreviousPassword)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, services.ErrInvalidPreviousPassword) {
			return nil, status.Error(codes.InvalidArgument, "invalid previous password")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &salonappv1.UpdateUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) AdminUpdateUser(ctx context.Context, req *salonappv1.AdminUpdateUserRequest) (*salonappv1.AdminUpdateUserResponse, error) {
	admin := util.UserFromContext(ctx)
	user, err := s.userService.AdminUpdateUser(ctx, admin.ID.String(), req.UserId, req.FullName, req.Password, req.Roles, req.IsActive)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, services.ErrUnauthorized) {
			return nil, status.Error(codes.PermissionDenied, "unauthorized")
		}
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	return &salonappv1.AdminUpdateUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) AddPhoneNumber(ctx context.Context, req *salonappv1.AddPhoneNumberRequest) (*salonappv1.AddPhoneNumberResponse, error) {
	user := util.UserFromContext(ctx)
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}
	if err := s.userService.AddPhoneNumber(ctx, user.ID.String(), req.PhoneNumber, req.Region); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidState):
			return nil, status.Error(codes.FailedPrecondition, "email must be verified")
		case errors.Is(err, services.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "phone already in use")
		default:
			return nil, status.Error(codes.Internal, "failed to add phone number")
		}
	}
	return &salonappv1.AddPhoneNumberResponse{Success: true, Message: "otp sent"}, nil
}

func (s *userServer) VerifyAddPhoneOTP(ctx context.Context, req *salonappv1.VerifyAddPhoneOTPRequest) (*salonappv1.VerifyAddPhoneOTPResponse, error) {
	user := util.UserFromContext(ctx)
	if req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_code is required")
	}
	if err := s.userService.VerifyAddPhone(ctx, user.ID.String(), req.OtpCode); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.InvalidArgument, "invalid or expired code")
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, services.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "phone already in use")
		default:
			return nil, status.Error(codes.Internal, "failed to verify phone")
		}
	}
	return &salonappv1.VerifyAddPhoneOTPResponse{Success: true, Message: "phone updated and verified"}, nil
}

func (s *userServer) AddEmail(ctx context.Context, req *salonappv1.AddEmailRequest) (*salonappv1.AddEmailResponse, error) {
	user := util.UserFromContext(ctx)
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if err := s.userService.AddEmail(ctx, user.ID.String(), req.Email); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidState):
			return nil, status.Error(codes.FailedPrecondition, "phone must be verified")
		case errors.Is(err, services.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "email already in use")
		default:
			return nil, status.Error(codes.Internal, "failed to add email")
		}
	}
	return &salonappv1.AddEmailResponse{Success: true, Message: "verification email sent"}, nil
}

func (s *userServer) VerifyAddEmailOTP(ctx context.Context, req *salonappv1.VerifyAddEmailOTPRequest) (*salonappv1.VerifyAddEmailOTPResponse, error) {
	user := util.UserFromContext(ctx)
	if req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_code is required")
	}
	if err := s.userService.VerifyAddEmail(ctx, user.ID.String(), req.OtpCode); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.InvalidArgument, "invalid or expired code")
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, services.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "email already in use")
		default:
			return nil, status.Error(codes.Internal, "failed to verify email")
		}
	}
	return &salonappv1.VerifyAddEmailOTPResponse{Success: true, Message: "email updated and verified"}, nil
}

func (s *userServer) LoginUser(ctx context.Context, req *salonappv1.LoginUserRequest) (*salonappv1.LoginUserResponse, error) {
	// Validate required params
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	tokenPair, err := s.userService.Login(
		ctx,
		req.Username,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		case errors.Is(err, services.ErrUserNotActive) || errors.Is(err, services.ErrInvalidEmailNotVerified):
			return nil, status.Error(codes.Unauthenticated, "user is not active")
		default:
			return nil, status.Error(codes.Internal, "failed to login user")
		}
	}

	// Map to proto response
	return &salonappv1.LoginUserResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresAt:        timestamppb.New(tokenPair.ExpiresAt),
		RefreshExpiresAt: timestamppb.New(tokenPair.RefreshExpiresAt),
		TokenType:        "bearer",
	}, nil
}

func (s *userServer) ResendEmailVerification(ctx context.Context, req *salonappv1.ResendEmailVerificationRequest) (*salonappv1.ResendEmailVerificationResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if err := s.userService.SendEmailVerification(ctx, req.Email); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to send verification email")
	}
	return &salonappv1.ResendEmailVerificationResponse{Success: true, Message: "verification email sent"}, nil
}

func (s *userServer) RecoverPassword(ctx context.Context, req *salonappv1.RecoverPasswordRequest) (*salonappv1.RecoverPasswordResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if err := s.userService.RequestPasswordReset(ctx, req.Email); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "failed to send recovery email")
		}
	}
	return &salonappv1.RecoverPasswordResponse{Success: true, Message: "password recovery email sent"}, nil
}

func (s *userServer) ResetPassword(ctx context.Context, req *salonappv1.ResetPasswordRequest) (*salonappv1.ResetPasswordResponse, error) {
	if req.Token == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "token and new_password are required")
	}
	if err := s.userService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, services.ErrWeakPassword):
			return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters")
		case errors.Is(err, services.ErrInvalidOrExpiredCode), errors.Is(err, services.ErrInvalidToken):
			return nil, status.Error(codes.InvalidArgument, "Invalid token")
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "failed to reset password")
		}
	}
	return &salonappv1.ResetPasswordResponse{Success: true, Message: "password reset successful"}, nil
}

func (s *userServer) RequestPhoneOTP(ctx context.Context, req *salonappv1.RequestPhoneOTPRequest) (*salonappv1.RequestPhoneOTPResponse, error) {
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}
	if err := s.userService.RequestPhoneOTP(ctx, req.PhoneNumber, req.Region); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to generate otp")
	}
	return &salonappv1.RequestPhoneOTPResponse{Success: true, Message: "otp generated"}, nil
}

func (s *userServer) VerifyPhoneOTP(ctx context.Context, req *salonappv1.VerifyPhoneOTPRequest) (*salonappv1.VerifyPhoneOTPResponse, error) {
	if req.VerificationToken == "" {
		return nil, status.Error(codes.InvalidArgument, "verification_token is required")
	}
	if req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_code is required")
	}
	pair, err := s.userService.VerifyRegisterPhoneUser(ctx, req.VerificationToken, req.OtpCode)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.InvalidArgument, "invalid or expired token or otp")
		case errors.Is(err, services.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to verify registration")
		}
	}
	return &salonappv1.VerifyPhoneOTPResponse{
		Success:          true,
		Message:          "phone number verified",
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		ExpiresAt:        timestamppb.New(pair.ExpiresAt),
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		TokenType:        "bearer",
	}, nil
}

func (s *userServer) VerifyEmailOTP(ctx context.Context, req *salonappv1.VerifyEmailOTPRequest) (*salonappv1.VerifyEmailOTPResponse, error) {
	if req.Email == "" || req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "email and otp_code are required")
	}
	if err := s.userService.VerifyEmailOTP(ctx, req.Email, req.OtpCode); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.InvalidArgument, "invalid or expired code")
		default:
			return nil, status.Error(codes.Internal, "failed to verify email")
		}
	}
	return &salonappv1.VerifyEmailOTPResponse{Success: true, Message: "email verified"}, nil
}
func (s *userServer) LoginWithPhone(ctx context.Context, req *salonappv1.LoginWithPhoneRequest) (*salonappv1.LoginWithPhoneResponse, error) {
	if req.PhoneNumber == "" || req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number and otp_code are required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}
	pair, err := s.userService.LoginWithPhone(ctx, req.PhoneNumber, req.OtpCode, req.Region)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.Unauthenticated, "invalid or expired code")
		default:
			return nil, status.Error(codes.Internal, "failed to login with phone")
		}
	}
	return &salonappv1.LoginWithPhoneResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		ExpiresAt:        timestamppb.New(pair.ExpiresAt),
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		TokenType:        "bearer",
	}, nil
}

func (s *userServer) RegisterPhoneUser(ctx context.Context, req *salonappv1.RegisterPhoneUserRequest) (*salonappv1.RegisterPhoneUserResponse, error) {
	if req.PhoneNumber == "" || req.FullName == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number and full_name are required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}
	token, err := s.userService.RegisterPhoneUser(ctx, req.PhoneNumber, req.FullName, req.Region)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &salonappv1.RegisterPhoneUserResponse{VerificationToken: token}, nil
}

func (s *userServer) userToProto(user *entities.User) *salonappv1.User {
	protoUser := &salonappv1.User{
		Id:              user.ID.String(),
		Email:           user.Email,
		PhoneNumber:     fromPtr(user.PhoneNumber),
		FullName:        fromPtr(user.FullName),
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
		Roles:           user.Roles,
	}

	if user.PhoneNumber != nil {
		protoUser.PhoneNumber = *user.PhoneNumber
	}

	protoUser.Roles = user.Roles

	return protoUser
}
