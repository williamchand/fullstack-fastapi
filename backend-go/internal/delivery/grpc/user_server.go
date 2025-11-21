package grpc

import (
	"context"
	"errors"

	userv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	userv1.UnimplementedUserServiceServer
	userService *services.UserService
}

func NewUserServer(userService *services.UserService) userv1.UserServiceServer {
	return &userServer{
		userService: userService,
	}
}

func (s *userServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &userv1.GetUserResponse{
		User: s.userToProto(user),
	}, nil
}

func (s *userServer) userToProto(user *entities.User) *userv1.User {
	protoUser := &userv1.User{
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

func fromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
