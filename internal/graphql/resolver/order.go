package resolver

import (
	"context"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// CreateOrder creates a new stock order.
func (m *MutationResolver) CreateOrder(
	ctx context.Context,
	input *models.CreateOrderInput,
) (*pb.Order, error) {
	attr := &pb.NewOrderAttributes{}
	if input.Comments != nil {
		attr.Comments = *input.Comments
	}
	attr.Consumer = input.Consumer
	attr.Courier = input.Courier
	attr.CourierAccount = input.CourierAccount
	attr.Items = input.Items
	attr.Payer = input.Payer
	attr.Payment = input.Payment
	attr.PurchaseOrderNum = *input.PurchaseOrderNum
	attr.Purchaser = input.Purchaser
	attr.Status = statusConverter(input.Status)
	o, err := m.GetOrderClient(registry.ORDER).CreateOrder(ctx, &pb.NewOrder{
		Data: &pb.NewOrder_Data{
			Type:       "order",
			Attributes: attr,
		},
	})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		m.Logger.Error(err)
		return nil, err
	}
	m.Logger.Debugf("successfully created new order with ID %s", o.Data.Id)
	return o, nil
}

// statusConverter converts the enum status string to protocol buffer int32 value
func statusConverter(e models.StatusEnum) pb.OrderStatus {
	var status pb.OrderStatus
	switch e {
	case "IN_PREPARATION":
		status = pb.OrderStatus_IN_PREPARATION
	case "GROWING":
		status = pb.OrderStatus_GROWING
	case "CANCELLED":
		status = pb.OrderStatus_CANCELLED
	case "SHIPPED":
		status = pb.OrderStatus_SHIPPED
	}
	return status
}

// convertPtrToStr converts a slice of string pointers to a slice of strings
func convertPtrToStr(items []*string) []string {
	var sl []string
	for _, n := range items {
		sl = append(sl, *n)
	}
	return sl
}

// UpdateOrder updates an existing stock order.
func (m *MutationResolver) UpdateOrder(
	ctx context.Context,
	id string,
	input *models.UpdateOrderInput,
) (*pb.Order, error) {
	g, err := m.GetOrderClient(registry.ORDER).
		GetOrder(ctx, &pb.OrderId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		m.Logger.Error(err)
		return nil, err
	}
	attr := &pb.OrderUpdateAttributes{}
	norm := normalizeUpdateOrderAttr(input)
	err = mapstructure.Decode(norm, attr)
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}
	if input.Status != nil {
		attr.Status = statusConverter(*input.Status)
	} else {
		attr.Status = g.Data.Attributes.Status
	}
	o, err := m.GetOrderClient(registry.ORDER).UpdateOrder(ctx, &pb.OrderUpdate{
		Data: &pb.OrderUpdate_Data{
			Type:       "order",
			Id:         id,
			Attributes: attr,
		},
	})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		m.Logger.Error(err)
		return nil, err
	}
	m.Logger.Debugf("successfully updated order with ID %s", o.Data.Id)
	return o, nil
}

func normalizeUpdateOrderAttr(
	attr *models.UpdateOrderInput,
) map[string]interface{} {
	fields := structs.Fields(attr)
	newAttr := make(map[string]interface{})
	for _, k := range fields {
		if k.Name() == "Status" {
			newAttr["Status"] = statusConverter(*attr.Status)
		} else if !k.IsZero() {
			newAttr[k.Name()] = k.Value()
		}
	}
	return newAttr
}

// Order retrieves an individual order by ID.
func (qrs *QueryResolver) Order(
	ctx context.Context,
	id string,
) (*pb.Order, error) {
	g, err := qrs.GetOrderClient(registry.ORDER).
		GetOrder(ctx, &pb.OrderId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	qrs.Logger.Debugf("successfully found order with id %s", id)
	return g, nil
}

// ListOrders retrieves all orders in the database.
func (qrs *QueryResolver) ListOrders(
	ctx context.Context,
	cursor *int,
	limit *int,
	filter *string,
) (*models.OrderListWithCursor, error) {
	c := resolverutils.GetCursor(cursor)
	l := resolverutils.GetLimit(limit)
	f := resolverutils.GetFilter(filter)
	list, err := qrs.GetOrderClient(registry.ORDER).
		ListOrders(ctx, &pb.ListParameters{Cursor: c, Limit: l, Filter: f})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	orders := []*pb.Order{}
	for _, n := range list.Data {
		item := &pb.Order{
			Data: &pb.Order_Data{
				Type:       n.Type,
				Id:         n.Id,
				Attributes: n.Attributes,
			},
		}
		orders = append(orders, item)
	}
	lm := int(l)
	return &models.OrderListWithCursor{
		Orders:         orders,
		NextCursor:     int(list.Meta.NextCursor),
		PreviousCursor: int(c),
		Limit:          &lm,
		TotalCount:     len(orders),
	}, nil
}
