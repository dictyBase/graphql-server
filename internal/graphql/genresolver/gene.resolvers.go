package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Goas is the resolver for the goas field.
func (r *geneResolver) Goas(ctx context.Context, obj *models.Gene) ([]*models.GOAnnotation, error) {
	panic("not implemented")
}

// Gene returns generated.GeneResolver implementation.
func (r *Resolver) Gene() generated.GeneResolver { return &geneResolver{r} }

type geneResolver struct{ *Resolver }
