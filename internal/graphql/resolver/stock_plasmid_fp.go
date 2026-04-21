package resolver

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	Pa "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
)

type listPlasmidsContext struct {
	// inputs (set by caller)
	client pb.StockServiceClient
	gctx   context.Context
	cursor *int
	limit  *int
	filter *models.PlasmidListFilter
	// computed by pipeline steps
	resolvedCursor cursor
	resolvedLimit  limit
	filterQuery    filterQuery
	collection     *pb.PlasmidCollection
}

// filterValidationPair threads the stable context alongside the evolving filter
// through the inner validation sub-pipe, keeping all validators univariate.
type filterValidationPair = Pa.Pair[listPlasmidsContext, *models.PlasmidListFilter]

type (
	cursor      = int64
	limit       = int64
	filterQuery = string
)

var (
	ptrString = func(s string) *string { return &s }

	formatPlasmidFieldQuery = F.Curry2(
		func(format string, value *string) string {
			return S.Format[string](format)(*value)
		},
	)

	compactOptionStrings = A.FilterMap(O.Fold(
		F.Constant(O.None[string]()),
		O.Some[string],
	))

	isNilPlasmidIDFilter = F.Pipe1(
		P.IsZero[*string](),
		P.ContraMap(func(pair filterValidationPair) *string {
			return Pa.Tail(pair).ID
		}),
	)

	checkNilPlasmidInStockFilter = E.FromPredicate(
		F.Pipe1(
			P.IsZero[*bool](),
			P.ContraMap(func(pair filterValidationPair) *bool {
				return Pa.Tail(pair).InStock
			}),
		),
		func(pair filterValidationPair) error {
			return fmt.Errorf(
				"plasmid list filter %v: in_stock filter is not yet supported in stock query conversion",
				Pa.Tail(pair),
			)
		},
	)

	convertPlasmidDataItem = func(item *pb.PlasmidCollection_Data) *models.Plasmid {
		return &models.Plasmid{
			ID:              item.Id,
			CreatedAt:       aphgrpc.ProtoTimeStamp(item.Attributes.CreatedAt),
			UpdatedAt:       aphgrpc.ProtoTimeStamp(item.Attributes.UpdatedAt),
			Summary:         &item.Attributes.Summary,
			EditableSummary: &item.Attributes.EditableSummary,
			Dbxrefs:         A.Map(ptrString)(item.Attributes.Dbxrefs),
			ImageMap:        &item.Attributes.ImageMap,
			Sequence:        &item.Attributes.Sequence,
			Name:            item.Attributes.Name,
			LazyStock: models.LazyStock{
				CreatedBy:    item.Attributes.CreatedBy,
				UpdatedBy:    item.Attributes.UpdatedBy,
				Depositor:    item.Attributes.Depositor,
				Genes:        item.Attributes.Genes,
				Publications: item.Attributes.Publications,
			},
		}
	}

	validateAndBuildPlasmidFilter = func(pair filterValidationPair) E.Either[error, listPlasmidsContext] {
		filter := Pa.Tail(pair)
		ctx := Pa.Head(pair)
		return F.Pipe2(
			filter.PlasmidType,
			E.FromPredicate(
				func(plasmidType models.PlasmidType) bool {
					return plasmidType.IsValid()
				},
				func(plasmidType models.PlasmidType) error {
					return fmt.Errorf("invalid plasmid type %s", plasmidType.String())
				},
			),
			E.Map[error](func(plasmidType models.PlasmidType) listPlasmidsContext {
				ctx.filterQuery = buildFilterQuery(filter, plasmidType)
				return ctx
			}),
		)
	}
)

func toEither[ERR, A any](ioe IOE.IOEither[ERR, A]) E.Either[ERR, A] {
	return ioe()
}

func onPlasmidListError(
	err error,
) T.Tuple2[error, *models.PlasmidListWithCursor] {
	return T.MakeTuple2(err, &models.PlasmidListWithCursor{})
}

func onPlasmidListSuccess(
	data *models.PlasmidListWithCursor,
) T.Tuple2[error, *models.PlasmidListWithCursor] {
	return T.MakeTuple2[error](nil, data)
}

func buildListPlasmidFilterQuery(
	ctx listPlasmidsContext,
) IOE.IOEither[error, listPlasmidsContext] {
	return F.Pipe7(
		ctx.filter,
		O.FromNillable[models.PlasmidListFilter],
		O.Map(func(filter *models.PlasmidListFilter) filterValidationPair {
			return Pa.MakePair(ctx, filter)
		}),
		O.GetOrElse(F.Constant(Pa.MakePair(ctx, &models.PlasmidListFilter{
			PlasmidType: models.PlasmidTypeAll,
		}))),
		E.FromPredicate(isNilPlasmidIDFilter, func(pair filterValidationPair) error {
			return fmt.Errorf(
				"plasmid list filter %v: id filter is not yet supported in stock query conversion",
				Pa.Tail(pair),
			)
		}),
		E.Chain(checkNilPlasmidInStockFilter),
		E.Chain(validateAndBuildPlasmidFilter),
		IOE.FromEither[error, listPlasmidsContext],
	)
}

func computeListPlasmidParams(
	ctx listPlasmidsContext,
) listPlasmidsContext {
	ctx.resolvedCursor = resolverutils.GetCursorFP(ctx.cursor)
	ctx.resolvedLimit = resolverutils.GetLimitFP(ctx.limit)
	return ctx
}

func fetchListPlasmidCollection(
	ctx listPlasmidsContext,
) IOE.IOEither[error, listPlasmidsContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.PlasmidCollection, error) {
			return ctx.client.ListPlasmids(
				ctx.gctx,
				&pb.StockParameters{
					Cursor: ctx.resolvedCursor,
					Limit:  ctx.resolvedLimit,
					Filter: ctx.filterQuery,
				},
			)
		}),
		IOE.Map[error](func(coll *pb.PlasmidCollection) listPlasmidsContext {
			ctx.collection = coll
			return ctx
		}),
	)
}

func extractListPlasmidResult(
	ctx listPlasmidsContext,
) *models.PlasmidListWithCursor {
	meta := ctx.collection.Meta
	lmt := int(meta.Limit)
	return &models.PlasmidListWithCursor{
		Plasmids:       A.Map(convertPlasmidDataItem)(ctx.collection.Data),
		NextCursor:     int(meta.NextCursor),
		PreviousCursor: int(ctx.resolvedCursor),
		Limit:          &lmt,
		TotalCount:     int(meta.Total),
	}
}

func buildFilterQuery(
	filter *models.PlasmidListFilter,
	plasmidType models.PlasmidType,
) filterQuery {
	return F.Pipe2(
		[]O.Option[string]{
			F.Pipe1(
				O.FromNillable(filter.Summary),
				O.Map(formatPlasmidFieldQuery("summary=~%s")),
			),
			F.Pipe1(
				O.FromNillable(filter.Name),
				O.Map(formatPlasmidFieldQuery("plasmid_name===%s")),
			),
			R.Lookup[string](plasmidType)(map[models.PlasmidType]string{
				models.PlasmidTypeRegular: fmt.Sprintf(
					"ontology==%s;tag==%s",
					registry.DictyPlasmidPropOntology,
					registry.RegularPlasmidTag,
				),
				models.PlasmidTypeGoldenBraid: fmt.Sprintf(
					"ontology==%s;tag==%s",
					registry.DictyPlasmidPropOntology,
					registry.GoldenBraidPlasmidTag,
				),
			}),
		},
		compactOptionStrings,
		A.Intercalate(S.Monoid)(";"),
	)
}
