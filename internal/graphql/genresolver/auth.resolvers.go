package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Identity is the resolver for the identity field.
func (r *authResolver) Identity(ctx context.Context, obj *auth.Auth) (*models.Identity, error) {
	panic("not implemented")
}

// Auth returns generated.AuthResolver implementation.
func (r *Resolver) Auth() generated.AuthResolver { return &authResolver{r} }

type authResolver struct{ *Resolver }
