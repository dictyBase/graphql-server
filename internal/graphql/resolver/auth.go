package resolver

import (
	"context"
	"fmt"

	"github.com/dictyBase/aphgrpc"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/graphql-server/internal/app/middleware"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

func (m *MutationResolver) Login(ctx context.Context, input *models.LoginInput) (*pb.Auth, error) {
	a := &pb.Auth{}
	// 1. Call service login method
	l, err := m.GetAuthClient(registry.AUTH).Login(ctx, &pb.NewLogin{
		ClientId:    input.ClientID,
		Scopes:      input.Scopes,
		State:       input.State,
		RedirectUrl: input.RedirectURL,
		Code:        input.Code,
		Provider:    input.Provider,
	})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		m.Logger.Error(err)
		return a, err
	}
	// 2, Set refresh token cookie with response
	arw := middleware.WriterFromContext(ctx)
	arw.RefreshToken = l.RefreshToken
	arw.Identifier = "login"
	// 3. Convert rest of response to Auth model
	a = &pb.Auth{
		Token:        l.Token,
		RefreshToken: arw.RefreshToken,
		User:         l.User,
		Identity:     l.Identity,
	}
	return a, nil
}
func (m *MutationResolver) Logout(ctx context.Context) (*models.Logout, error) {
	// 1. Check for refresh token
	arw := middleware.WriterFromContext(ctx)
	if arw.RefreshToken == "" {
		nerr := aphgrpc.HandleNotFoundError(ctx, fmt.Errorf("refresh token does not exist"))
		errorutils.AddGQLError(ctx, nerr)
		return nil, nerr
	}
	// 2. Call Logout service method
	_, err := m.GetAuthClient(registry.AUTH).Logout(ctx, &pb.NewRefreshToken{
		RefreshToken: arw.RefreshToken,
	})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		m.Logger.Error(err)
		return nil, err
	}
	// 3. Set expired cookie
	arw.Identifier = "logout"
	return &models.Logout{
		Success: true,
	}, nil
}
func (q *QueryResolver) GetRefreshToken(ctx context.Context, token string) (*pb.Auth, error) {
	a := &pb.Auth{}
	// 1. Get the refresh token from the cookie
	arw := middleware.WriterFromContext(ctx)
	if arw.RefreshToken == "" {
		nerr := aphgrpc.HandleNotFoundError(ctx, fmt.Errorf("refresh token does not exist"))
		errorutils.AddGQLError(ctx, nerr)
		return a, nerr
	}
	// 2. Get user and identity information
	a, err := q.GetAuthClient(registry.AUTH).Relogin(ctx, &pb.NewRelogin{
		RefreshToken: arw.RefreshToken,
	})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		q.Logger.Error(err)
		return a, err
	}
	// 3. Set new refresh token cookie from response
	arw.RefreshToken = a.RefreshToken
	arw.Identifier = "login"
	// 4. Return new JWT and user/identity data
	return a, nil
}
