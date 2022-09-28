package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input *models.LoginInput) (*auth.Auth, error) {
	panic("not implemented")
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (*models.Logout, error) {
	panic("not implemented")
}

// CreateContent is the resolver for the createContent field.
func (r *mutationResolver) CreateContent(ctx context.Context, input *models.CreateContentInput) (*content.Content, error) {
	panic("not implemented")
}

// UpdateContent is the resolver for the updateContent field.
func (r *mutationResolver) UpdateContent(ctx context.Context, input *models.UpdateContentInput) (*content.Content, error) {
	panic("not implemented")
}

// DeleteContent is the resolver for the deleteContent field.
func (r *mutationResolver) DeleteContent(ctx context.Context, id string) (*models.DeleteContent, error) {
	panic("not implemented")
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *models.CreateOrderInput) (*order.Order, error) {
	panic("not implemented")
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, id string, input *models.UpdateOrderInput) (*order.Order, error) {
	panic("not implemented")
}

// CreateStrain is the resolver for the createStrain field.
func (r *mutationResolver) CreateStrain(ctx context.Context, input *models.CreateStrainInput) (*models.Strain, error) {
	panic("not implemented")
}

// CreatePlasmid is the resolver for the createPlasmid field.
func (r *mutationResolver) CreatePlasmid(ctx context.Context, input *models.CreatePlasmidInput) (*models.Plasmid, error) {
	panic("not implemented")
}

// UpdateStrain is the resolver for the updateStrain field.
func (r *mutationResolver) UpdateStrain(ctx context.Context, id string, input *models.UpdateStrainInput) (*models.Strain, error) {
	panic("not implemented")
}

// UpdatePlasmid is the resolver for the updatePlasmid field.
func (r *mutationResolver) UpdatePlasmid(ctx context.Context, id string, input *models.UpdatePlasmidInput) (*models.Plasmid, error) {
	panic("not implemented")
}

// DeleteStock is the resolver for the deleteStock field.
func (r *mutationResolver) DeleteStock(ctx context.Context, id string) (*models.DeleteStock, error) {
	panic("not implemented")
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *models.CreateUserInput) (*user.User, error) {
	panic("not implemented")
}

// CreateUserRoleRelationship is the resolver for the createUserRoleRelationship field.
func (r *mutationResolver) CreateUserRoleRelationship(ctx context.Context, userID string, roleID string) (*user.User, error) {
	panic("not implemented")
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*user.User, error) {
	panic("not implemented")
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*models.DeleteUser, error) {
	panic("not implemented")
}

// CreateRole is the resolver for the createRole field.
func (r *mutationResolver) CreateRole(ctx context.Context, input *models.CreateRoleInput) (*user.Role, error) {
	panic("not implemented")
}

// CreateRolePermissionRelationship is the resolver for the createRolePermissionRelationship field.
func (r *mutationResolver) CreateRolePermissionRelationship(ctx context.Context, roleID string, permissionID string) (*user.Role, error) {
	panic("not implemented")
}

// UpdateRole is the resolver for the updateRole field.
func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input *models.UpdateRoleInput) (*user.Role, error) {
	panic("not implemented")
}

// DeleteRole is the resolver for the deleteRole field.
func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*models.DeleteRole, error) {
	panic("not implemented")
}

// CreatePermission is the resolver for the createPermission field.
func (r *mutationResolver) CreatePermission(ctx context.Context, input *models.CreatePermissionInput) (*user.Permission, error) {
	panic("not implemented")
}

// UpdatePermission is the resolver for the updatePermission field.
func (r *mutationResolver) UpdatePermission(ctx context.Context, id string, input *models.UpdatePermissionInput) (*user.Permission, error) {
	panic("not implemented")
}

// DeletePermission is the resolver for the deletePermission field.
func (r *mutationResolver) DeletePermission(ctx context.Context, id string) (*models.DeletePermission, error) {
	panic("not implemented")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
