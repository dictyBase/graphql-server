package resolver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/user"

	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

func (mrs *MutationResolver) CreatePermission(
	ctx context.Context,
	input *models.CreatePermissionInput,
) (*pb.Permission, error) {
	n, err := mrs.GetPermissionClient(registry.PERMISSION).
		CreatePermission(ctx, &pb.CreatePermissionRequest{
			Data: &pb.CreatePermissionRequest_Data{
				Type: "permission",
				Attributes: &pb.PermissionAttributes{
					Permission:  input.Permission,
					Description: input.Description,
					Resource:    input.Resource,
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf("successfully created new permission with ID %d", n.Data.Id)
	return n, nil
}

func (mrs *MutationResolver) UpdatePermission(
	ctx context.Context,
	id string,
	input *models.UpdatePermissionInput,
) (*pb.Permission, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	n, err := mrs.GetPermissionClient(registry.PERMISSION).
		UpdatePermission(ctx, &pb.UpdatePermissionRequest{
			Id: i,
			Data: &pb.UpdatePermissionRequest_Data{
				Id:   i,
				Type: "permission",
				Attributes: &pb.PermissionAttributes{
					Permission:  input.Permission,
					Description: input.Description,
					Resource:    input.Resource,
					UpdatedAt:   aphgrpc.TimestampProto(time.Now()),
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	o, err := mrs.GetPermissionClient(registry.PERMISSION).
		GetPermission(ctx, &jsonapi.GetRequestWithFields{Id: n.Data.Id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf("successfully updated permission with ID %d", o.Data.Id)
	return o, nil
}

func (mrs *MutationResolver) DeletePermission(
	ctx context.Context,
	id string,
) (*models.DeletePermission, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}

	_, err = mrs.GetPermissionClient(registry.PERMISSION).
		DeletePermission(ctx, &jsonapi.DeleteRequest{Id: i})
	if err != nil {
		mrs.Logger.Error(err)
		return &models.DeletePermission{
			Success: false,
		}, err
	}
	mrs.Logger.Debugf("successfully deleted permission with ID %s", id)
	return &models.DeletePermission{
		Success: true,
	}, nil
}

func (qrs *QueryResolver) Permission(
	ctx context.Context,
	id string,
) (*pb.Permission, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	g, err := qrs.GetPermissionClient(registry.PERMISSION).
		GetPermission(ctx, &jsonapi.GetRequestWithFields{Id: i})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	qrs.Logger.Debugf("successfully found permission with ID %s", id)
	return g, nil
}

func (qrs *QueryResolver) ListPermissions(
	ctx context.Context,
) ([]*pb.Permission, error) {
	permissions := []*pb.Permission{}
	l, err := qrs.GetPermissionClient(registry.PERMISSION).
		ListPermissions(ctx, &jsonapi.SimpleListRequest{})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	for _, n := range l.Data {
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
	qrs.Logger.Debugf(
		"successfully provided list of %d permissions",
		len(permissions),
	)
	return permissions, nil
}
