package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Goas is the resolver for the goas field.
func (r *geneResolver) Goas(ctx context.Context, obj *models.Gene) ([]*models.GOAnnotation, error) {
	panic("not implemented")
}

type geneResolver struct{ *Resolver }
