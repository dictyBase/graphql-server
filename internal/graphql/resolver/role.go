package resolver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dictyBase/apihelpers/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/user"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

func (m *MutationResolver) CreateRole(ctx context.Context, input *models.CreateRoleInput) (*pb.Role, error) {
	n, err := m.GetRoleClient(registry.ROLE).CreateRole(context.Background(), &pb.CreateRoleRequest{
		Data: &pb.CreateRoleRequest_Data{
			Type: "role",
			Attributes: &pb.RoleAttributes{
				Role:        input.Role,
				Description: input.Description,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new role %s", err)
	}
	m.Logger.Debugf("successfully created new role with ID %d", n.Data.Id)
	return n, nil
}
func (m *MutationResolver) CreateRolePermissionRelationship(ctx context.Context, roleId string, permissionId string) (*pb.Role, error) {
	rid, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", roleId, err)
	}
	pid, err := strconv.ParseInt(permissionId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", permissionId, err)
	}
	rr, err := m.GetRoleClient(registry.ROLE).CreatePermissionRelationship(ctx, &jsonapi.DataCollection{
		Id: rid,
		Data: []*jsonapi.Data{
			{
				Type: "permission",
				Id:   pid,
			},
		}})
	if err != nil {
		return nil, fmt.Errorf("error in creating permission relationship with role %s", err)
	}
	m.Logger.Debugf("successfully created role ID %d relationship permission with ID %d %s", rid, pid, rr)
	g, err := m.GetRoleClient(registry.ROLE).GetRole(ctx, &jsonapi.GetRequest{Id: rid})
	if err != nil {
		return nil, fmt.Errorf("error in getting role by ID %d: %s", rid, err)
	}
	return g, nil
}
func (m *MutationResolver) UpdateRole(ctx context.Context, id string, input *models.UpdateRoleInput) (*pb.Role, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	n, err := m.GetRoleClient(registry.ROLE).UpdateRole(context.Background(), &pb.UpdateRoleRequest{
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
		return nil, fmt.Errorf("error updating role %d: %s", n.Data.Id, err)
	}
	o, err := m.GetRoleClient(registry.ROLE).GetRole(context.Background(), &jsonapi.GetRequest{Id: i})
	if err != nil {
		return nil, fmt.Errorf("error fetching recently updated role: %s", err)
	}
	m.Logger.Debugf("successfully updated role with ID %d", n.Data.Id)
	return o, nil
}
func (m *MutationResolver) DeleteRole(ctx context.Context, id string) (*models.DeleteRole, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	if _, err := m.GetRoleClient(registry.ROLE).DeleteRole(context.Background(), &jsonapi.DeleteRequest{Id: i}); err != nil {
		return &models.DeleteRole{
			Success: false,
		}, fmt.Errorf("error deleting role with ID %s: %s", id, err)
	}
	m.Logger.Debugf("successfully deleted role with ID %s", id)
	return &models.DeleteRole{
		Success: true,
	}, nil
}
func (q *QueryResolver) Role(ctx context.Context, id string) (*pb.Role, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	g, err := q.GetRoleClient(registry.ROLE).GetRole(ctx, &jsonapi.GetRequest{Id: i})
	if err != nil {
		return nil, fmt.Errorf("error in getting role by ID %d: %s", i, err)
	}
	q.Logger.Debugf("successfully found role with ID %s", id)
	return g, nil
}
func (q *QueryResolver) ListRoles(ctx context.Context) ([]pb.Role, error) {
	roles := []pb.Role{}
	l, err := q.GetRoleClient(registry.ROLE).ListRoles(ctx, &jsonapi.SimpleListRequest{})
	if err != nil {
		return nil, fmt.Errorf("error in listing roles %s", err)
	}
	for _, n := range l.Data {
		item := pb.Role{
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
	q.Logger.Debugf("successfully provided list of %d roles", len(roles))
	return roles, nil
}
