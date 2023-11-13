package user

import (
	"context"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/sirupsen/logrus"
)

type RoleResolver struct {
	Client *authentication.LogtoClient
	Logger *logrus.Entry
}

func (rrs *RoleResolver) ID(ctx context.Context, obj *pb.Role) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}
func (rrs *RoleResolver) Role(ctx context.Context, obj *pb.Role) (string, error) {
	return obj.Data.Attributes.Role, nil
}
func (rrs *RoleResolver) Description(ctx context.Context, obj *pb.Role) (string, error) {
	return obj.Data.Attributes.Description, nil
}
func (rrs *RoleResolver) CreatedAt(ctx context.Context, obj *pb.Role) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}
func (rrs *RoleResolver) UpdatedAt(ctx context.Context, obj *pb.Role) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}
func (rrs *RoleResolver) Permissions(ctx context.Context, obj *pb.Role) ([]*pb.Permission, error) {
	permissions := []*pb.Permission{}
	rp, err := rrs.Client.GetRelatedPermissions(ctx, &jsonapi.RelationshipRequest{Id: obj.Data.Id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		rrs.Logger.Error(err)
		return permissions, err
	}
	for _, n := range rp.Data {
		item := &pb.Permission{
			Data: &pb.PermissionData{
				Type: "permission",
				Id:   n.Id,
				Attributes: &pb.PermissionAttributes{
					Permission:  n.Attributes.Permission,
					Description: n.Attributes.Description,
					CreatedAt:   n.Attributes.CreatedAt,
					UpdatedAt:   n.Attributes.UpdatedAt,
					Resource:    n.Attributes.Resource,
				},
			},
		}
		permissions = append(permissions, item)
	}
	rrs.Logger.Infof("successfully retrieved list of %d permissions for role ID %d", len(permissions), obj.Data.Id)
	return permissions, nil
}
