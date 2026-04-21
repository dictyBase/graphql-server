package resolver

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	Pa "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
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
	total       = int64
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

	validateAndBuildPlasmidFilter = F.Curry2(
		func(state listPlasmidParamsTuple, filter *models.PlasmidListFilter) E.Either[error, listPlasmidFilterTuple] {
			return F.Pipe3(
				filter.PlasmidType,
				E.FromPredicate(
					func(plasmidType models.PlasmidType) bool {
						return plasmidType.IsValid()
					},
					func(plasmidType models.PlasmidType) error {
						return fmt.Errorf("invalid plasmid type %s", plasmidType.String())
					},
				),
				E.Map[error](func(plasmidType models.PlasmidType) listPlasmidFilterBuildTuple {
					return T.MakeTuple5(state.F1, state.F2, state.F3, filter, plasmidType)
				}),
				E.Map[error](buildListPlasmidFilterTuple),
			)
		},
	)
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

func computeListPlasmidParams(
	ctx listPlasmidsContext,
) listPlasmidParamsTuple {
	return T.MakeTuple3(
		ctx,
		resolverutils.GetCursorFP(ctx.cursor),
		resolverutils.GetLimitFP(ctx.limit),
	)
}

func buildListPlasmidFilterQuery(
	state listPlasmidParamsTuple,
) IOE.IOEither[error, listPlasmidFilterTuple] {
	return F.Pipe6(
		state.F1.filter,
		O.FromNillable[models.PlasmidListFilter],
		O.GetOrElse(F.Constant(&models.PlasmidListFilter{
			PlasmidType: models.PlasmidTypeAll,
		})),
		E.FromPredicate(
			isNilPlasmidIDFilter,
			func(filter *models.PlasmidListFilter) error {
				return fmt.Errorf(
					"plasmid list filter %v: id filter is not yet supported in stock query conversion",
					filter,
				)
			},
		),
		E.Chain(checkNilPlasmidInStockFilter),
		E.Chain(validateAndBuildPlasmidFilter(state)),
		IOE.FromEither[error, listPlasmidFilterTuple],
	)
}

func fetchListPlasmidCollection(
	state listPlasmidFilterTuple,
) IOE.IOEither[error, listPlasmidCollectionTuple] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.PlasmidCollection, error) {
			return state.F1.client.ListPlasmids(
				state.F1.gctx,
				&pb.StockParameters{
					Cursor: state.F2,
					Limit:  state.F3,
					Filter: state.F4,
				},
			)
		}),
		IOE.Map[error](func(coll *pb.PlasmidCollection) listPlasmidCollectionTuple {
			return T.MakeTuple5(state.F1, state.F2, state.F3, state.F4, coll)
		}),
	)
}

func extractListPlasmidResult(
	state listPlasmidCollectionTuple,
) *models.PlasmidListWithCursor {
	return F.Pipe2(
		T.MakeTuple2(
			state.F5.Data,
			plasmidResultContext{
				limit:      state.F5.Meta.Limit,
				nextCursor: state.F5.Meta.NextCursor,
				total:      state.F5.Meta.Total,
				cursor:     state.F2,
			},
		),
		T.Map2(A.Map(convertPlasmidDataItem), F.Identity[plasmidResultContext]),
		func(tuple plasmidResultTuple) *models.PlasmidListWithCursor {
			lmt := int(tuple.F2.limit)
			return &models.PlasmidListWithCursor{
				Plasmids:       tuple.F1,
				NextCursor:     int(tuple.F2.nextCursor),
				PreviousCursor: int(tuple.F2.cursor),
				Limit:          &lmt,
				TotalCount:     int(tuple.F2.total),
			}
		},
	)
}

func buildListPlasmidFilterTuple(ctx listPlasmidFilterBuildTuple) listPlasmidFilterTuple {
	return T.MakeTuple4(
		ctx.F1,
		ctx.F2,
		ctx.F3,
		F.Pipe2(
			[]O.Option[string]{
				F.Pipe1(
					O.FromNillable(ctx.F4.Summary),
					O.Map(formatPlasmidFieldQuery("summary=~%s")),
				),
				F.Pipe1(
					O.FromNillable(ctx.F4.Name),
					O.Map(formatPlasmidFieldQuery("plasmid_name===%s")),
				),
				R.Lookup[string](ctx.F5)(map[models.PlasmidType]string{
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
		),
	)
}
