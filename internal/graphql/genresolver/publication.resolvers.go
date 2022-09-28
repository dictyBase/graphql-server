package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Rank is the resolver for the rank field.
func (r *authorResolver) Rank(ctx context.Context, obj *publication.Author) (*string, error) {
	panic("not implemented")
}

// Authors is the resolver for the authors field.
func (r *publicationResolver) Authors(ctx context.Context, obj *models.Publication) ([]*publication.Author, error) {
	panic("not implemented")
}

// Author returns generated.AuthorResolver implementation.
func (r *Resolver) Author() generated.AuthorResolver { return &authorResolver{r} }

// Publication returns generated.PublicationResolver implementation.
func (r *Resolver) Publication() generated.PublicationResolver { return &publicationResolver{r} }

type authorResolver struct{ *Resolver }
type publicationResolver struct{ *Resolver }
