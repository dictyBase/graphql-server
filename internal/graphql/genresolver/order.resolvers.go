package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// ID is the resolver for the id field.
func (r *orderResolver) ID(ctx context.Context, obj *order.Order) (string, error) {
	panic("not implemented")
}

// CreatedAt is the resolver for the created_at field.
func (r *orderResolver) CreatedAt(ctx context.Context, obj *order.Order) (*time.Time, error) {
	panic("not implemented")
}

// UpdatedAt is the resolver for the updated_at field.
func (r *orderResolver) UpdatedAt(ctx context.Context, obj *order.Order) (*time.Time, error) {
	panic("not implemented")
}

// Courier is the resolver for the courier field.
func (r *orderResolver) Courier(ctx context.Context, obj *order.Order) (*string, error) {
	panic("not implemented")
}

// CourierAccount is the resolver for the courier_account field.
func (r *orderResolver) CourierAccount(ctx context.Context, obj *order.Order) (*string, error) {
	panic("not implemented")
}

// Comments is the resolver for the comments field.
func (r *orderResolver) Comments(ctx context.Context, obj *order.Order) (*string, error) {
	panic("not implemented")
}

// Payment is the resolver for the payment field.
func (r *orderResolver) Payment(ctx context.Context, obj *order.Order) (*string, error) {
	panic("not implemented")
}

// PurchaseOrderNum is the resolver for the purchase_order_num field.
func (r *orderResolver) PurchaseOrderNum(ctx context.Context, obj *order.Order) (*string, error) {
	panic("not implemented")
}

// Status is the resolver for the status field.
func (r *orderResolver) Status(ctx context.Context, obj *order.Order) (*models.StatusEnum, error) {
	panic("not implemented")
}

// Consumer is the resolver for the consumer field.
func (r *orderResolver) Consumer(ctx context.Context, obj *order.Order) (*user.User, error) {
	panic("not implemented")
}

// Payer is the resolver for the payer field.
func (r *orderResolver) Payer(ctx context.Context, obj *order.Order) (*user.User, error) {
	panic("not implemented")
}

// Purchaser is the resolver for the purchaser field.
func (r *orderResolver) Purchaser(ctx context.Context, obj *order.Order) (*user.User, error) {
	panic("not implemented")
}

// Items is the resolver for the items field.
func (r *orderResolver) Items(ctx context.Context, obj *order.Order) ([]models.Stock, error) {
	panic("not implemented")
}

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

type orderResolver struct{ *Resolver }
