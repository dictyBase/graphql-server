package stock

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var numRgx = regexp.MustCompile("[0-9]+")

type PlasmidResolver struct {
	Client           pb.StockServiceClient
	UserClient       authentication.LogtoClient
	AnnotationClient annotation.TaggedAnnotationServiceClient
	Registry         registry.Registry
	Logger           *logrus.Entry
}

func (prs *PlasmidResolver) CreatedBy(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	u, err := userByEmail(ctx, obj.CreatedBy, prs.UserClient, prs.Logger)
	if err != nil {
		prs.Logger.Error(err)
		return nil, err
	}
	return u, nil
}

func (prs *PlasmidResolver) UpdatedBy(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	u, err := userByEmail(ctx, obj.UpdatedBy, prs.UserClient, prs.Logger)
	if err != nil {
		prs.Logger.Error(err)
		return nil, err
	}
	return u, nil
}

func (prs *PlasmidResolver) Depositor(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	if len(obj.Depositor) == 0 {
		return fakeUser(), nil
	}
	u, err := userByEmail(ctx, obj.Depositor, prs.UserClient, prs.Logger)
	if err != nil {
		prs.Logger.Error(err)
		return nil, err
	}
	return u, nil
}

func (prs *PlasmidResolver) Genes(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Gene, error) {
	g := []*models.Gene{}
	redis := prs.Registry.GetRedisRepository(cache.RedisKey)
	for _, v := range obj.Genes {
		gene, err := cache.GetGeneFromCache(ctx, redis, v)
		if err != nil {
			prs.Logger.Error(err)
			continue
		}
		g = append(g, gene)
	}
	return g, nil
}

func (prs *PlasmidResolver) Publications(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Publication, error) {
	pubs := make([]*models.Publication, 0)
	for _, id := range obj.Publications {
		pub, err := fetch.FetchPublication(
			ctx, prs.Registry.GetRedisRepository(cache.RedisKey),
			prs.Registry.GetAPIEndpoint(registry.PUBLICATION), id,
		)
		if err != nil {
			errorutils.AddGQLError(ctx, err)
			prs.Logger.Error(err)
			return pubs, fmt.Errorf(
				"error in fetching publication %s for plasmid %s",
				id,
				err,
			)
		}
		pubs = append(pubs, pub)
	}
	return pubs, nil
}

func (prs *PlasmidResolver) InStock(
	ctx context.Context,
	obj *models.Plasmid,
) (bool, error) {
	id := obj.ID
	_, err := prs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.PlasmidInvTag,
			Ontology: registry.PlasmidInvOnto,
			EntryId:  id,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		errorutils.AddGQLError(ctx, err)
		prs.Logger.Errorf("error getting %s from inventory: %v", id, err)
		return false, err
	}
	return true, nil
}

/*
* Note: none of the below have been implemented yet.
 */
func (prs *PlasmidResolver) Keywords(
	ctx context.Context,
	obj *models.Plasmid,
) ([]string, error) {
	return []string{""}, nil
}

func (prs *PlasmidResolver) GenbankAccession(
	ctx context.Context,
	obj *models.Plasmid,
) (*string, error) {
	s := ""
	return &s, nil
}

func ConvertToPlasmidModel(
	id string,
	attr *pb.PlasmidAttributes,
) *models.Plasmid {
	return &models.Plasmid{
		ID:              id,
		CreatedAt:       aphgrpc.ProtoTimeStamp(attr.CreatedAt),
		UpdatedAt:       aphgrpc.ProtoTimeStamp(attr.UpdatedAt),
		Summary:         &attr.Summary,
		EditableSummary: &attr.EditableSummary,
		Dbxrefs:         sliceConverter(attr.Dbxrefs),
		ImageMap:        &attr.ImageMap,
		Sequence:        &attr.Sequence,
		Name:            attr.Name,
		LazyStock: models.LazyStock{
			CreatedBy:    attr.CreatedBy,
			UpdatedBy:    attr.UpdatedBy,
			Depositor:    attr.Depositor,
			Genes:        attr.Genes,
			Publications: attr.Publications,
		},
	}
}

func userByEmail(
	ctx context.Context,
	email string,
	client authentication.LogtoClient,
	logger *logrus.Entry,
) (*user.User, error) {
	userResp, err := client.UserWithEmail(email)
	if err != nil {
		userErr := fmt.Errorf("unable to retreieve user %s", err)
		errorutils.AddGQLError(ctx, userErr)
		logger.Error(userErr)
		return nil, userErr
	}
	matches := numRgx.FindAllString(userResp.ID, -1)
	if len(matches) == 0 {
		nonumErr := fmt.Errorf(
			"cannot convert user id to number %s",
			userResp.ID,
		)
		errorutils.AddGQLError(ctx, nonumErr)
		logger.Error(nonumErr)
		return nil, nonumErr
	}
	userID, err := strconv.ParseInt(strings.Join(matches, ""), 10, 64)
	if err != nil {
		parseErr := fmt.Errorf("unable to convert user id to integer %s", err)
		errorutils.AddGQLError(ctx, parseErr)
		logger.Error(parseErr)
		return nil, parseErr
	}
	logger.Debugf("successfully found user with id %s", email)
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

func fakeUser() *user.User {
	return &user.User{
		Data: &user.UserData{
			Type: "user",
			Id:   2375,
			Attributes: &user.UserAttributes{
				FirstName:    "Dicty",
				LastName:     "Stock Center",
				Email:        "dictystocks@northwestern.edu",
				Organization: "Northwestern University",
				IsActive:     true,
			},
		},
	}
}
