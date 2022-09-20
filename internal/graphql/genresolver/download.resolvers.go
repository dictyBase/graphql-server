package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Downloads is the resolver for the downloads field.
func (r *organismResolver) Downloads(ctx context.Context, obj *models.Organism) ([]*models.Download, error) {
	panic("not implemented")
}

// Organism returns generated.OrganismResolver implementation.
func (r *Resolver) Organism() generated.OrganismResolver { return &organismResolver{r} }

type organismResolver struct{ *Resolver }
