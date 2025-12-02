package grpc

import (
	"context"
	"errors"

	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"

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
	user := auth.UserFromContext(ctx)
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

func (s *userServer) CreateUser(ctx context.Context, req *salonappv1.CreateUserRequest) (*salonappv1.CreateUserResponse, error) {
	user, err := s.userService.CreateUser(ctx, req.Email, req.Password, req.FullName, []entities.RoleEnum{entities.RoleCustomer}, false)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &salonappv1.CreateUserResponse{
		User: s.userToProto(user),
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
	user, err := s.userService.UpdateUser(ctx, req.Id, *req.Email, *req.FullName, *req.PhoneNumber, []entities.RoleEnum{entities.RoleCustomer})
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &salonappv1.UpdateUserResponse{
		User: s.userToProto(user),
	}, nil
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
			return nil, status.Error(codes.PermissionDenied, "user is not active")
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

func (s *userServer) RequestPhoneOTP(ctx context.Context, req *salonappv1.RequestPhoneOTPRequest) (*salonappv1.RequestPhoneOTPResponse, error) {
	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number is required")
	}
	if err := s.userService.RequestPhoneOTP(ctx, req.PhoneNumber); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to generate otp")
	}
	return &salonappv1.RequestPhoneOTPResponse{Success: true, Message: "otp generated"}, nil
}

func (s *userServer) VerifyPhoneOTP(ctx context.Context, req *salonappv1.VerifyPhoneOTPRequest) (*salonappv1.VerifyPhoneOTPResponse, error) {
	if req.PhoneNumber == "" || req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number and otp_code are required")
	}
	if err := s.userService.VerifyPhoneOTP(ctx, req.PhoneNumber, req.OtpCode); err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, services.ErrInvalidOrExpiredCode):
			return nil, status.Error(codes.InvalidArgument, "invalid or expired code")
		default:
			return nil, status.Error(codes.Internal, "failed to verify phone")
		}
	}
	return &salonappv1.VerifyPhoneOTPResponse{Success: true, Message: "phone verified"}, nil
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
	pair, err := s.userService.LoginWithPhone(ctx, req.PhoneNumber, req.OtpCode)
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
	user, err := s.userService.RegisterPhoneUser(ctx, req.PhoneNumber, req.FullName)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &salonappv1.RegisterPhoneUserResponse{User: s.userToProto(user)}, nil
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
