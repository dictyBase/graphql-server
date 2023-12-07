package resolver

import (
	"context"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

func (m *MutationResolver) Login(
	ctx context.Context,
	input *models.LoginInput,
) (*pb.Auth, error) {
	return &pb.Auth{}, nil
}
func (m *MutationResolver) Logout(ctx context.Context) (*models.Logout, error) {
	return &models.Logout{}, nil
}

func (qrs *QueryResolver) GetRefreshToken(
	ctx context.Context,
	token string,
) (*pb.Auth, error) {
	return &pb.Auth{}, nil
}
