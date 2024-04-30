package order

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	stockPB "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/sirupsen/logrus"
)

var numRgx = regexp.MustCompile("[0-9]+")

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
	email := obj.Data.Attributes.Consumer
	g, err := ord.userByEmail(ctx, email)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		ord.Logger.Error(err)
		return nil, err
	}
	ord.Logger.Debugf("successfully found user with email %s", email)
	return g, nil
}

func (ord *OrderResolver) Payer(
	ctx context.Context,
	obj *pb.Order,
) (*user.User, error) {
	email := obj.Data.Attributes.Payer
	g, err := ord.userByEmail(ctx, email)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		ord.Logger.Error(err)
		return nil, err
	}
	ord.Logger.Debugf("successfully found user with email %s", email)
	return g, nil
}

func (ord *OrderResolver) Purchaser(
	ctx context.Context,
	obj *pb.Order,
) (*user.User, error) {
	email := obj.Data.Attributes.Purchaser
	g, err := ord.userByEmail(ctx, email)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		ord.Logger.Error(err)
		return nil, err
	}
	ord.Logger.Debugf("successfully found user with email %s", email)
	return g, nil
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

func (ord *OrderResolver) userByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	userResp, err := ord.UserClient.UserWithEmail(email)
	if err != nil {
		userErr := fmt.Errorf("unable to retreieve user %s", err)
		errorutils.AddGQLError(ctx, userErr)
		ord.Logger.Error(userErr)
		return nil, userErr
	}
	matches := numRgx.FindAllString(userResp.ID, -1)
	if len(matches) == 0 {
		nonumErr := fmt.Errorf(
			"cannot convert user id to number %s",
			userResp.ID,
		)
		errorutils.AddGQLError(ctx, nonumErr)
		ord.Logger.Error(nonumErr)
		return nil, nonumErr
	}
	userID, err := strconv.ParseInt(strings.Join(matches, ""), 10, 64)
	if err != nil {
		parseErr := fmt.Errorf("unable to convert user id to integer %s", err)
		errorutils.AddGQLError(ctx, parseErr)
		ord.Logger.Error(parseErr)
		return nil, parseErr
	}
	ord.Logger.Debugf("successfully found user with id %s", email)
	return &user.User{
		Data: &user.UserData{
			Type: "user",
			Id:   userID,
			Attributes: &user.UserAttributes{
				FirstName:    userResp.Username,
				LastName:     userResp.Name,
				Email:        userResp.PrimaryEmail,
				Organization: userResp.CustomData.Institution,
				FirstAddress: userResp.CustomData.Address,
				City:         userResp.CustomData.City,
				State:        userResp.CustomData.State,
				Zipcode:      userResp.CustomData.Zipcode,
				Country:      userResp.CustomData.Country,
				Phone:        userResp.PrimaryPhone,
				IsActive:     true,
			},
		},
	}, nil
}
