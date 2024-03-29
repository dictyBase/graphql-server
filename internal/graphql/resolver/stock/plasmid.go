package stock

import (
	"context"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlasmidResolver struct {
	Client           pb.StockServiceClient
	UserClient       user.UserServiceClient
	AnnotationClient annotation.TaggedAnnotationServiceClient
	Registry         registry.Registry
	Logger           *logrus.Entry
}

func (r *PlasmidResolver) CreatedBy(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	u, err := getUserByEmail(ctx, r.UserClient, obj.CreatedBy)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), err
	}
	return u, nil
}

func (r *PlasmidResolver) UpdatedBy(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	u, err := getUserByEmail(ctx, r.UserClient, obj.UpdatedBy)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), err
	}
	return u, nil
}

func (r *PlasmidResolver) Depositor(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	u, err := getUserByEmail(ctx, r.UserClient, obj.Depositor)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), nil
	}
	return u, nil
}

func (r *PlasmidResolver) Genes(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Gene, error) {
	g := []*models.Gene{}
	redis := r.Registry.GetRedisRepository(cache.RedisKey)
	for _, v := range obj.Genes {
		gene, err := cache.GetGeneFromCache(ctx, redis, v)
		if err != nil {
			r.Logger.Error(err)
			continue
		}
		g = append(g, gene)
	}
	return g, nil
}

func (r *PlasmidResolver) Publications(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Publication, error) {
	pubs := make([]*models.Publication, 0)
	for _, id := range obj.Publications {
		endpoint := r.Registry.GetAPIEndpoint(registry.PUBLICATION)
		p, err := fetch.FetchPublication(ctx, endpoint, id)
		if err != nil {
			errorutils.AddGQLError(ctx, err)
			r.Logger.Error(err)
			return pubs, err
		}
		pubs = append(pubs, p)
	}
	return pubs, nil
}

func (r *PlasmidResolver) InStock(
	ctx context.Context,
	obj *models.Plasmid,
) (bool, error) {
	id := obj.ID
	_, err := r.AnnotationClient.GetEntryAnnotation(
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
		r.Logger.Errorf("error getting %s from inventory: %v", id, err)
		return false, err
	}
	return true, nil
}

/*
* Note: none of the below have been implemented yet.
 */
func (r *PlasmidResolver) Keywords(
	ctx context.Context,
	obj *models.Plasmid,
) ([]string, error) {
	return []string{""}, nil
}

func (r *PlasmidResolver) GenbankAccession(
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
