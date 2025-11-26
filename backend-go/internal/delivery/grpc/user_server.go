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
	user, err := s.userService.CreateUser(ctx, req.Email, req.Password, req.FullName, req.PhoneNumber, []entities.RoleEnum{entities.RoleCustomer}, false)
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
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    timestamppb.New(tokenPair.ExpiresAt),
		TokenType:    "bearer",
	}, nil
}

func (s *userServer) userToProto(user *entities.User) *salonappv1.User {
	protoUser := &salonappv1.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		FullName:  fromPtr(user.FullName),
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     user.Roles,
	}

	if user.PhoneNumber != nil {
		protoUser.PhoneNumber = *user.PhoneNumber
	}

	protoUser.Roles = user.Roles

	return protoUser
}
