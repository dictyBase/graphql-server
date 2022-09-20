package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
)

// ID is the resolver for the id field.
func (r *contentResolver) ID(ctx context.Context, obj *content.Content) (string, error) {
	panic("not implemented")
}

// Name is the resolver for the name field.
func (r *contentResolver) Name(ctx context.Context, obj *content.Content) (string, error) {
	panic("not implemented")
}

// Slug is the resolver for the slug field.
func (r *contentResolver) Slug(ctx context.Context, obj *content.Content) (string, error) {
	panic("not implemented")
}

// CreatedBy is the resolver for the created_by field.
func (r *contentResolver) CreatedBy(ctx context.Context, obj *content.Content) (*user.User, error) {
	panic("not implemented")
}

// UpdatedBy is the resolver for the updated_by field.
func (r *contentResolver) UpdatedBy(ctx context.Context, obj *content.Content) (*user.User, error) {
	panic("not implemented")
}

// CreatedAt is the resolver for the created_at field.
func (r *contentResolver) CreatedAt(ctx context.Context, obj *content.Content) (*time.Time, error) {
	panic("not implemented")
}

// UpdatedAt is the resolver for the updated_at field.
func (r *contentResolver) UpdatedAt(ctx context.Context, obj *content.Content) (*time.Time, error) {
	panic("not implemented")
}

// Content is the resolver for the content field.
func (r *contentResolver) Content(ctx context.Context, obj *content.Content) (string, error) {
	panic("not implemented")
}

// Namespace is the resolver for the namespace field.
func (r *contentResolver) Namespace(ctx context.Context, obj *content.Content) (string, error) {
	panic("not implemented")
}

// Content returns generated.ContentResolver implementation.
func (r *Resolver) Content() generated.ContentResolver { return &contentResolver{r} }

type contentResolver struct{ *Resolver }
