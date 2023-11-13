package content

import (
	"context"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/sirupsen/logrus"
)

type ContentResolver struct {
	Client     pb.ContentServiceClient
	UserClient *authentication.LogtoClient
	Logger     *logrus.Entry
}

func (rcs *ContentResolver) ID(ctx context.Context, obj *pb.Content) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}
func (rcs *ContentResolver) Name(ctx context.Context, obj *pb.Content) (string, error) {
	return obj.Data.Attributes.Name, nil
}
func (rcs *ContentResolver) Slug(ctx context.Context, obj *pb.Content) (string, error) {
	return obj.Data.Attributes.Slug, nil
}
func (rcs *ContentResolver) CreatedBy(ctx context.Context, obj *pb.Content) (*user.User, error) {
	user := user.User{}
	id := obj.Data.Attributes.CreatedBy
	g, err := rcs.UserClient.GetUser(ctx, &jsonapi.GetRequest{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		rcs.Logger.Error(err)
		return &user, err
	}
	rcs.Logger.Debugf("successfully found user with id %d", id)
	return g, nil
}
func (rcs *ContentResolver) UpdatedBy(ctx context.Context, obj *pb.Content) (*user.User, error) {
	user := user.User{}
	id := obj.Data.Attributes.UpdatedBy
	g, err := rcs.UserClient.GetUser(ctx, &jsonapi.GetRequest{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		rcs.Logger.Error(err)
		return &user, err
	}
	rcs.Logger.Debugf("successfully found user with id %d", id)
	return g, nil
}
func (rcs *ContentResolver) CreatedAt(ctx context.Context, obj *pb.Content) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}
func (rcs *ContentResolver) UpdatedAt(ctx context.Context, obj *pb.Content) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}
func (rcs *ContentResolver) Content(ctx context.Context, obj *pb.Content) (string, error) {
	return obj.Data.Attributes.Content, nil
}
func (rcs *ContentResolver) Namespace(ctx context.Context, obj *pb.Content) (string, error) {
	return obj.Data.Attributes.Namespace, nil
}
