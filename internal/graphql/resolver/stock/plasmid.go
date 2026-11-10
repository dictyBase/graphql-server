package stock

import (
	"context"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/dictyBase/aphgrpc"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/collection"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
)

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
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.CreatedBy),
		resolverutils.LookupClientLens.Set(prs.UserClient),
		resolverutils.LookupLoggerLens.Set(prs.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

func (prs *PlasmidResolver) UpdatedBy(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.UpdatedBy),
		resolverutils.LookupClientLens.Set(prs.UserClient),
		resolverutils.LookupLoggerLens.Set(prs.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

func (prs *PlasmidResolver) Depositor(
	ctx context.Context,
	obj *models.Plasmid,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(obj.Depositor),
		resolverutils.LookupClientLens.Set(prs.UserClient),
		resolverutils.LookupLoggerLens.Set(prs.Logger),
		resolverutils.StockDepositor,
	)
	return res.F2, res.F1
}

func (prs *PlasmidResolver) Genes(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Gene, error) {
	return collection.Map(obj.Genes, func(g string) *models.Gene {
		return &models.Gene{ID: g, Name: g}
	}), nil
}

func (prs *PlasmidResolver) Publications(
	ctx context.Context,
	obj *models.Plasmid,
) ([]*models.Publication, error) {
	res := F.Pipe6(
		stockPubsState{},
		stockPubsCtxLens.Set(ctx),
		stockPubsRedisLens.Set(prs.Registry.GetRedisRepository(cache.RedisKey)),
		stockPubsIDLens.Set(obj.ID),
		stockPubsLoggerLens.Set(prs.Logger),
		stockPubsRawIDsLens.Set(obj.Publications),
		fetchStockPublications,
	)
	return res.F2, res.F1
}

func (prs *PlasmidResolver) InStock(
	ctx context.Context,
	obj *models.Plasmid,
) (bool, error) {
	result := F.Pipe2(
		inStockForEntity(
			ctx,
			prs.AnnotationClient,
			obj.ID,
			models.StockEntityTypePlasmid,
		),
		ioeutils.ToEither[error, bool],
		E.Fold(onInventoryError, onInventorySuccess),
	)

	if result.F1 != nil {
		errorutils.AddGQLError(ctx, result.F1)
		prs.Logger.Error(result.F1)
		return false, result.F1
	}
	return result.F2, nil
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
