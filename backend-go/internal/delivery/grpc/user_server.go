package grpc

import (
	"context"
	"errors"

	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	genprotov1.UnimplementedUserServiceServer
	userService *services.UserService
}

func NewUserServer(userService *services.UserService) genprotov1.UserServiceServer {
	return &userServer{
		userService: userService,
	}
}

func (s *userServer) GetUser(ctx context.Context, req *genprotov1.GetUserRequest) (*genprotov1.GetUserResponse, error) {
	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &genprotov1.GetUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *genprotov1.CreateUserRequest) (*genprotov1.CreateUserResponse, error) {
	user, err := s.userService.CreateUser(ctx, req.Email, req.Password, req.FullName, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &genprotov1.CreateUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *genprotov1.UpdateUserRequest) (*genprotov1.UpdateUserResponse, error) {
	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &genprotov1.GetUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) userToProto(user *entities.User) *genprotov1.User {
	protoUser := &genprotov1.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		FullName:  fromPtr(user.FullName),
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	if user.PhoneNumber != nil {
		protoUser.PhoneNumber = *user.PhoneNumber
	}

	// Add roles
	for _, role := range user.Roles {
		protoUser.Roles = append(protoUser.Roles, role.Name)
	}

	return protoUser
}
