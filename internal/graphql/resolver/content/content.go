package content

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/sirupsen/logrus"
)

type ContentResolver struct {
	Client     pb.ContentServiceClient
	UserClient authentication.LogtoClient
	Logger     *logrus.Entry
}

func (rcs *ContentResolver) ID(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}

func (rcs *ContentResolver) Name(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Name, nil
}

func (rcs *ContentResolver) Slug(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Slug, nil
}

func (rcs *ContentResolver) CreatedBy(
	ctx context.Context,
	obj *pb.Content,
) (*user.User, error) {
	email := obj.Data.Attributes.CreatedBy
	userResp, err := rcs.UserClient.UserWithEmail(email)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("unable to retreieve user %s", err),
		)
		rcs.Logger.Error(err)
		return nil, err
	}
	userID, err := strconv.ParseInt(userResp.ID, 10, 64)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("unable to convert user id to integer %s", err),
		)
		rcs.Logger.Error(err)
		return nil, err
	}
	rcs.Logger.Debugf("successfully found user with id %s", email)
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

func (rcs *ContentResolver) UpdatedBy(
	ctx context.Context,
	obj *pb.Content,
) (*user.User, error) {
	email := obj.Data.Attributes.UpdatedBy
	userResp, err := rcs.UserClient.UserWithEmail(email)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("unable to retreieve user %s", err),
		)
		rcs.Logger.Error(err)
		return nil, err
	}
	userID, err := strconv.ParseInt(userResp.ID, 10, 64)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("unable to convert user id to integer %s", err),
		)
		rcs.Logger.Error(err)
		return nil, err
	}
	rcs.Logger.Debugf("successfully found user with id %s", email)
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

func (rcs *ContentResolver) CreatedAt(
	ctx context.Context,
	obj *pb.Content,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}

func (rcs *ContentResolver) UpdatedAt(
	ctx context.Context,
	obj *pb.Content,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}

func (rcs *ContentResolver) Content(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Content, nil
}

func (rcs *ContentResolver) Namespace(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Namespace, nil
}
