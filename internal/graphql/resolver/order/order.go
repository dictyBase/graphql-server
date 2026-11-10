package order

import (
	"context"
	"fmt"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	stockPB "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/sirupsen/logrus"
)

type OrderResolver struct {
	Client      pb.OrderServiceClient
	StockClient stockPB.StockServiceClient
	UserClient  authentication.LogtoClient
	Logger      *logrus.Entry
}

func (ord *OrderResolver) ID(ctx context.Context, obj *pb.Order) (string, error) {
	return obj.Data.Id, nil
}

func (ord *OrderResolver) CreatedAt(
	ctx context.Context,
	obj *pb.Order,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}

func (ord *OrderResolver) UpdatedAt(
	ctx context.Context,
	obj *pb.Order,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}

func (ord *OrderResolver) Courier(
	ctx context.Context,
	obj *pb.Order,
) (*string, error) {
	return &obj.Data.Attributes.Courier, nil
}

func (ord *OrderResolver) CourierAccount(
	ctx context.Context,
	obj *pb.Order,
) (*string, error) {
	return &obj.Data.Attributes.CourierAccount, nil
}

func (ord *OrderResolver) Comments(
	ctx context.Context,
	obj *pb.Order,
) (*string, error) {
	return &obj.Data.Attributes.Comments, nil
}

func (ord *OrderResolver) Payment(
	ctx context.Context,
	obj *pb.Order,
) (*string, error) {
	return &obj.Data.Attributes.Payment, nil
}

func (ord *OrderResolver) PurchaseOrderNum(
	ctx context.Context,
	obj *pb.Order,
) (*string, error) {
	return &obj.Data.Attributes.PurchaseOrderNum, nil
}

func (ord *OrderResolver) Status(
	ctx context.Context,
	obj *pb.Order,
) (*models.StatusEnum, error) {
	status := obj.Data.Attributes.Status
	switch status {
	case pb.OrderStatus_IN_PREPARATION:
		s := models.StatusEnumInPreparation
		return &s, nil
	case pb.OrderStatus_GROWING:
		s := models.StatusEnumGrowing
		return &s, nil
	case pb.OrderStatus_CANCELLED:
		s := models.StatusEnumCancelled
		return &s, nil
	case pb.OrderStatus_SHIPPED:
		s := models.StatusEnumShipped
		return &s, nil
	default:
		return nil, fmt.Errorf("incompatible order status")
	}
}

func (ord *OrderResolver) Consumer(
	ctx context.Context,
	obj *pb.Order,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.Data.Attributes.Consumer),
		resolverutils.LookupClientLens.Set(ord.UserClient),
		resolverutils.LookupLoggerLens.Set(ord.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

func (ord *OrderResolver) Payer(
	ctx context.Context,
	obj *pb.Order,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.Data.Attributes.Payer),
		resolverutils.LookupClientLens.Set(ord.UserClient),
		resolverutils.LookupLoggerLens.Set(ord.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

func (ord *OrderResolver) Purchaser(
	ctx context.Context,
	obj *pb.Order,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.Data.Attributes.Purchaser),
		resolverutils.LookupClientLens.Set(ord.UserClient),
		resolverutils.LookupLoggerLens.Set(ord.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

// ConsumerInfo resolves the shipper profile embedded in the order.
func (ord *OrderResolver) ConsumerInfo(
	ctx context.Context,
	obj *pb.Order,
) (*pb.UserInfo, error) {
	return obj.Data.Attributes.ConsumerInfo, nil
}

// PayerInfo resolves the payer profile embedded in the order.
func (ord *OrderResolver) PayerInfo(
	ctx context.Context,
	obj *pb.Order,
) (*pb.UserInfo, error) {
	return obj.Data.Attributes.PayerInfo, nil
}

func (ord *OrderResolver) Items(
	ctx context.Context,
	obj *pb.Order,
) ([]models.Stock, error) {
	stocks := []models.Stock{}
	for _, id := range obj.Data.Attributes.Items {
		if id[:3] == "DBS" {
			gs, err := ord.StockClient.GetStrain(ctx, &stockPB.StockId{Id: id})
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				ord.Logger.Error(err)
				return stocks, err
			}
			stocks = append(
				stocks,
				stock.ConvertToStrainModel(id, gs.Data.Attributes),
			)
		}
		if id[:3] == "DBP" {
			gp, err := ord.StockClient.GetPlasmid(ctx, &stockPB.StockId{Id: id})
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				ord.Logger.Error(err)
				return stocks, err
			}
			stocks = append(
				stocks,
				stock.ConvertToPlasmidModel(id, gp.Data.Attributes),
			)
		}
	}
	return stocks, nil
}
