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

func (mrs *MutationResolver) CreateRole(
	ctx context.Context,
	input *models.CreateRoleInput,
) (*pb.Role, error) {
	n, err := mrs.GetRoleClient(registry.ROLE).
		CreateRole(ctx, &pb.CreateRoleRequest{
			Data: &pb.CreateRoleRequest_Data{
				Type: "role",
				Attributes: &pb.RoleAttributes{
					Role:        input.Role,
					Description: input.Description,
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf("successfully created new role with ID %d", n.Data.Id)
	return n, nil
}

func (mrs *MutationResolver) CreateRolePermissionRelationship(
	ctx context.Context,
	roleID string,
	permissionID string,
) (*pb.Role, error) {
	rid, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"error in parsing string %s to int %s",
			roleID,
			err,
		)
	}
	pid, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"error in parsing string %s to int %s",
			permissionID,
			err,
		)
	}
	rr, err := mrs.GetRoleClient(registry.ROLE).
		CreatePermissionRelationship(ctx, &jsonapi.DataCollection{
			Id: rid,
			Data: []*jsonapi.Data{
				{
					Type: "permission",
					Id:   pid,
				},
			}})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf(
		"successfully created role ID %d relationship permission with ID %d %s",
		rid,
		pid,
		rr,
	)
	g, err := mrs.GetRoleClient(registry.ROLE).
		GetRole(ctx, &jsonapi.GetRequest{Id: rid})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	return g, nil
}

func (mrs *MutationResolver) UpdateRole(
	ctx context.Context,
	id string,
	input *models.UpdateRoleInput,
) (*pb.Role, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	n, err := mrs.GetRoleClient(registry.ROLE).
		UpdateRole(ctx, &pb.UpdateRoleRequest{
			Id: i,
			Data: &pb.UpdateRoleRequest_Data{
				Id:   i,
				Type: "role",
				Attributes: &pb.RoleAttributes{
					Role:        input.Role,
					Description: input.Description,
					UpdatedAt:   aphgrpc.TimestampProto(time.Now()),
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	o, err := mrs.GetRoleClient(registry.ROLE).
		GetRole(ctx, &jsonapi.GetRequest{Id: n.Data.Id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf("successfully updated role with ID %d", o.Data.Id)
	return o, nil
}

func (mrs *MutationResolver) DeleteRole(
	ctx context.Context,
	id string,
) (*models.DeleteRole, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	if _, err := mrs.GetRoleClient(registry.ROLE).DeleteRole(ctx, &jsonapi.DeleteRequest{Id: i}); err != nil {
		mrs.Logger.Error(err)
		return &models.DeleteRole{
			Success: false,
		}, err
	}
	mrs.Logger.Debugf("successfully deleted role with ID %s", id)
	return &models.DeleteRole{
		Success: true,
	}, nil
}

func (qrs *QueryResolver) Role(
	ctx context.Context,
	id string,
) (*pb.Role, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	g, err := qrs.GetRoleClient(registry.ROLE).
		GetRole(ctx, &jsonapi.GetRequest{Id: i})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	qrs.Logger.Debugf("successfully found role with ID %s", id)
	return g, nil
}
func (qrs *QueryResolver) ListRoles(ctx context.Context) ([]*pb.Role, error) {
	roles := []*pb.Role{}
	l, err := qrs.GetRoleClient(registry.ROLE).
		ListRoles(ctx, &jsonapi.SimpleListRequest{})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	for _, n := range l.Data {
		item := &pb.Role{
			Data: &pb.RoleData{
				Type: "role",
				Id:   n.Id,
				Attributes: &pb.RoleAttributes{
					Role:        n.Attributes.Role,
					Description: n.Attributes.Description,
					CreatedAt:   n.Attributes.CreatedAt,
					UpdatedAt:   n.Attributes.UpdatedAt,
				},
			},
		}
		roles = append(roles, item)
	}
	qrs.Logger.Debugf("successfully provided list of %d roles", len(roles))
	return roles, nil
}
