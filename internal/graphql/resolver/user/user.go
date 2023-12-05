package user

import (
	"context"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/sirupsen/logrus"
)

type UserResolver struct {
	Client authentication.LogtoClient
	Logger *logrus.Entry
}

func (urs *UserResolver) ID(ctx context.Context, obj *pb.User) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}

func (urs *UserResolver) FirstName(
	ctx context.Context,
	obj *pb.User,
) (string, error) {
	return obj.Data.Attributes.FirstName, nil
}

func (urs *UserResolver) LastName(
	ctx context.Context,
	obj *pb.User,
) (string, error) {
	return obj.Data.Attributes.LastName, nil
}

func (urs *UserResolver) Email(
	ctx context.Context,
	obj *pb.User,
) (string, error) {
	return obj.Data.Attributes.Email, nil
}

func (urs *UserResolver) Organization(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.Organization, nil
}

func (urs *UserResolver) GroupName(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.GroupName, nil
}

func (urs *UserResolver) FirstAddress(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.FirstAddress, nil
}

func (urs *UserResolver) SecondAddress(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.SecondAddress, nil
}

func (urs *UserResolver) City(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.City, nil
}

func (urs *UserResolver) State(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.State, nil
}

func (urs *UserResolver) Zipcode(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.Zipcode, nil
}

func (urs *UserResolver) Country(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.Country, nil
}

func (urs *UserResolver) Phone(
	ctx context.Context,
	obj *pb.User,
) (*string, error) {
	return &obj.Data.Attributes.Phone, nil
}

func (urs *UserResolver) IsActive(
	ctx context.Context,
	obj *pb.User,
) (bool, error) {
	return obj.Data.Attributes.IsActive, nil
}

func (urs *UserResolver) CreatedAt(
	ctx context.Context,
	obj *pb.User,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}

func (urs *UserResolver) UpdatedAt(
	ctx context.Context,
	obj *pb.User,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}

func (urs *UserResolver) Roles(
	ctx context.Context,
	obj *pb.User,
) ([]*pb.Role, error) {
	roles := []*pb.Role{}
	resp, err := urs.Client.Roles(strconv.FormatInt(obj.Data.Id, 10))
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		urs.Logger.Error(err)
		return roles, err
	}
	for _, rdata := range resp {
		roleId, _ := strconv.ParseInt(rdata.Id, 10, 64)
		item := &pb.Role{
			Data: &pb.RoleData{
				Type: "role",
				Id:   roleId,
				Attributes: &pb.RoleAttributes{
					Role:        rdata.Name,
					Description: rdata.Description,
					CreatedAt:   obj.Data.Attributes.CreatedAt,
					UpdatedAt:   obj.Data.Attributes.UpdatedAt,
				},
			},
		}
		roles = append(roles, item)
	}
	urs.Logger.Infof(
		"successfully retrieved list of %d roles for user ID %d",
		len(roles),
		obj.Data.Id,
	)
	return roles, nil
}
