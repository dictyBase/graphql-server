package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
)

// ID is the resolver for the id field.
func (r *permissionResolver) ID(ctx context.Context, obj *user.Permission) (string, error) {
	panic("not implemented")
}

// Permission is the resolver for the permission field.
func (r *permissionResolver) Permission(ctx context.Context, obj *user.Permission) (string, error) {
	panic("not implemented")
}

// Description is the resolver for the description field.
func (r *permissionResolver) Description(ctx context.Context, obj *user.Permission) (string, error) {
	panic("not implemented")
}

// CreatedAt is the resolver for the created_at field.
func (r *permissionResolver) CreatedAt(ctx context.Context, obj *user.Permission) (*time.Time, error) {
	panic("not implemented")
}

// UpdatedAt is the resolver for the updated_at field.
func (r *permissionResolver) UpdatedAt(ctx context.Context, obj *user.Permission) (*time.Time, error) {
	panic("not implemented")
}

// Resource is the resolver for the resource field.
func (r *permissionResolver) Resource(ctx context.Context, obj *user.Permission) (*string, error) {
	panic("not implemented")
}

// ID is the resolver for the id field.
func (r *roleResolver) ID(ctx context.Context, obj *user.Role) (string, error) {
	panic("not implemented")
}

// Role is the resolver for the role field.
func (r *roleResolver) Role(ctx context.Context, obj *user.Role) (string, error) {
	panic("not implemented")
}

// Description is the resolver for the description field.
func (r *roleResolver) Description(ctx context.Context, obj *user.Role) (string, error) {
	panic("not implemented")
}

// CreatedAt is the resolver for the created_at field.
func (r *roleResolver) CreatedAt(ctx context.Context, obj *user.Role) (*time.Time, error) {
	panic("not implemented")
}

// UpdatedAt is the resolver for the updated_at field.
func (r *roleResolver) UpdatedAt(ctx context.Context, obj *user.Role) (*time.Time, error) {
	panic("not implemented")
}

// Permissions is the resolver for the permissions field.
func (r *roleResolver) Permissions(ctx context.Context, obj *user.Role) ([]*user.Permission, error) {
	panic("not implemented")
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *user.User) (string, error) {
	panic("not implemented")
}

// FirstName is the resolver for the first_name field.
func (r *userResolver) FirstName(ctx context.Context, obj *user.User) (string, error) {
	panic("not implemented")
}

// LastName is the resolver for the last_name field.
func (r *userResolver) LastName(ctx context.Context, obj *user.User) (string, error) {
	panic("not implemented")
}

// Email is the resolver for the email field.
func (r *userResolver) Email(ctx context.Context, obj *user.User) (string, error) {
	panic("not implemented")
}

// Organization is the resolver for the organization field.
func (r *userResolver) Organization(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// GroupName is the resolver for the group_name field.
func (r *userResolver) GroupName(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// FirstAddress is the resolver for the first_address field.
func (r *userResolver) FirstAddress(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// SecondAddress is the resolver for the second_address field.
func (r *userResolver) SecondAddress(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// City is the resolver for the city field.
func (r *userResolver) City(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// State is the resolver for the state field.
func (r *userResolver) State(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// Zipcode is the resolver for the zipcode field.
func (r *userResolver) Zipcode(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// Country is the resolver for the country field.
func (r *userResolver) Country(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// Phone is the resolver for the phone field.
func (r *userResolver) Phone(ctx context.Context, obj *user.User) (*string, error) {
	panic("not implemented")
}

// IsActive is the resolver for the is_active field.
func (r *userResolver) IsActive(ctx context.Context, obj *user.User) (bool, error) {
	panic("not implemented")
}

// CreatedAt is the resolver for the created_at field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *user.User) (*time.Time, error) {
	panic("not implemented")
}

// UpdatedAt is the resolver for the updated_at field.
func (r *userResolver) UpdatedAt(ctx context.Context, obj *user.User) (*time.Time, error) {
	panic("not implemented")
}

// Roles is the resolver for the roles field.
func (r *userResolver) Roles(ctx context.Context, obj *user.User) ([]*user.Role, error) {
	panic("not implemented")
}

// Permission returns generated.PermissionResolver implementation.
func (r *Resolver) Permission() generated.PermissionResolver { return &permissionResolver{r} }

// Role returns generated.RoleResolver implementation.
func (r *Resolver) Role() generated.RoleResolver { return &roleResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type permissionResolver struct{ *Resolver }
type roleResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
