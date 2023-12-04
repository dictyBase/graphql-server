package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoleResolver struct {
	Client *authentication.LogtoClient
	Logger *logrus.Entry
}

func (rrs *RoleResolver) ID(ctx context.Context, obj *pb.Role) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}

func (rrs *RoleResolver) Role(
	ctx context.Context,
	obj *pb.Role,
) (string, error) {
	return obj.Data.Attributes.Role, nil
}

func (rrs *RoleResolver) Description(
	ctx context.Context,
	obj *pb.Role,
) (string, error) {
	return obj.Data.Attributes.Description, nil
}

func (rrs *RoleResolver) CreatedAt(
	ctx context.Context,
	obj *pb.Role,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}

func (rrs *RoleResolver) UpdatedAt(
	ctx context.Context,
	obj *pb.Role,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}

func (rrs *RoleResolver) Permissions(
	ctx context.Context,
	obj *pb.Role,
) ([]*pb.Permission, error) {
	permissions := []*pb.Permission{}
	roleId := strconv.FormatInt(obj.Data.Id, 10)
	permResp, err := rrs.Client.Permissions(roleId)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("error in fetching permissions %s", err),
		)
		rrs.Logger.Error(err)
		return permissions, err
	}
	for _, perm := range permResp {
		permId, err := strconv.ParseInt(perm.ID, 10, 64)
		if err != nil {
			errorutils.AddGQLError(
				ctx,
				fmt.Errorf(
					"error in converting permission id to integer %s",
					err,
				),
			)
			rrs.Logger.Error(err)
			return permissions, err
		}
		item := &pb.Permission{
			Data: &pb.PermissionData{
				Type: "permission",
				Id:   permId,
				Attributes: &pb.PermissionAttributes{
					Permission:  perm.Name,
					Description: perm.Description,
					Resource:    perm.Resource.Name,
					CreatedAt:   timestamppb.Now(),
					UpdatedAt:   timestamppb.Now(),
				},
			},
		}
		permissions = append(permissions, item)
	}
	rrs.Logger.Infof(
		"successfully retrieved list of %d permissions for role ID %d",
		len(permissions),
		obj.Data.Id,
	)

	return permissions, nil
}
