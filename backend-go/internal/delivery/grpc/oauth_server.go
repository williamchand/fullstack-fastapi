package grpc

import (
	"context"
	"errors"

	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/salonapp/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type oAuthServer struct {
	salonappv1.UnimplementedOAuthServiceServer
	oauth *services.OAuthService
}

func NewOAuthServer(oauth *services.OAuthService) salonappv1.OAuthServiceServer {
	return &oAuthServer{
		oauth: oauth,
	}
}

//
// ────────────────────────────────────────────────────────────────
//   1. Get OAuth URL
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) GetOAuthURL(ctx context.Context, req *salonappv1.GetOAuthURLRequest) (*salonappv1.GetOAuthURLResponse, error) {

	url, state, err := s.oauth.GetAuthURL(ctx, req.Provider, req.RedirectUri)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate OAuth URL")
	}

	return &salonappv1.GetOAuthURLResponse{
		Url:   url,
		State: state,
	}, nil
}

//
// ────────────────────────────────────────────────────────────────
//   2. Handle OAuth Callback
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) HandleOAuthCallback(ctx context.Context, req *salonappv1.HandleOAuthCallbackRequest) (*salonappv1.HandleOAuthCallbackResponse, error) {

	user, tokenInfo, err := s.oauth.HandleCallback(ctx, req.Provider, req.Code, req.State)
	if err != nil {
		if errors.Is(err, services.ErrInvalidOAuthCode) {
			return nil, status.Error(codes.InvalidArgument, "invalid authorization code")
		}
		if errors.Is(err, services.ErrOAuthUnauthorized) {
			return nil, status.Error(codes.Unauthenticated, "authentication failed")
		}
		return nil, status.Error(codes.Internal, "failed to handle OAuth callback")
	}

	return &salonappv1.HandleOAuthCallbackResponse{
		User:         s.userToProto(user),
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
		ExpiresAt:    timestamppb.New(tokenInfo.ExpiresAt),
		IsNewUser:    tokenInfo.IsNewUser,
	}, nil
}

//
// ────────────────────────────────────────────────────────────────
//   3. Link OAuth Account
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) LinkOAuthAccount(ctx context.Context, req *salonappv1.LinkOAuthAccountRequest) (*salonappv1.LinkOAuthAccountResponse, error) {

	account, err := s.oauth.LinkAccount(ctx, req.Provider, req.UserId, req.Code, req.State)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to link OAuth account")
	}

	return &salonappv1.LinkOAuthAccountResponse{
		Account: s.oauthAccountToProto(account),
	}, nil
}

//
// ────────────────────────────────────────────────────────────────
//   4. Unlink OAuth Account
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) UnlinkOAuthAccount(ctx context.Context, req *salonappv1.UnlinkOAuthAccountRequest) (*salonappv1.UnlinkOAuthAccountResponse, error) {

	err := s.oauth.UnlinkAccount(ctx, req.Provider, req.UserId)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to unlink OAuth account")
	}

	return &salonappv1.UnlinkOAuthAccountResponse{
		Success: true,
	}, nil
}

//
// ────────────────────────────────────────────────────────────────
//   5. Get User OAuth Accounts
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) GetUserOAuthAccounts(ctx context.Context, req *salonappv1.GetUserOAuthAccountsRequest) (*salonappv1.GetUserOAuthAccountsResponse, error) {

	accounts, err := s.oauth.GetUserAccounts(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch user OAuth accounts")
	}

	resp := &salonappv1.GetUserOAuthAccountsResponse{}
	for _, acc := range accounts {
		resp.Accounts = append(resp.Accounts, s.oauthAccountToProto(acc))
	}

	return resp, nil
}

//
// ────────────────────────────────────────────────────────────────
//   Mapping: Domain → Proto
// ────────────────────────────────────────────────────────────────
//

func (s *oAuthServer) userToProto(u *entities.User) *salonappv1.User {
	p := &salonappv1.User{
		Id:              u.ID.String(),
		Email:           u.Email,
		FullName:        fromPtr(u.FullName),
		PhoneNumber:     fromPtr(u.PhoneNumber),
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
		IsPhoneVerified: u.IsPhoneVerified,
		CreatedAt:       timestamppb.New(u.CreatedAt),
		UpdatedAt:       timestamppb.New(u.UpdatedAt),
	}

	if u.LastLoginAt != nil {
		p.LastLoginAt = timestamppb.New(*u.LastLoginAt)
	}

	return p
}

func (s *oAuthServer) oauthAccountToProto(a *entities.OAuthAccount) *salonappv1.OAuthAccount {
	return &salonappv1.OAuthAccount{
		Id:         a.ID,
		Provider:   a.Provider,
		ProviderId: a.ProviderID,
		UserId:     a.UserID,
		Email:      a.Email,
		Name:       a.Name,
		Picture:    a.Picture,
		LinkedAt:   timestamppb.New(a.LinkedAt),
		LastUsedAt: timestamppb.New(a.LastUsedAt),
	}
}
