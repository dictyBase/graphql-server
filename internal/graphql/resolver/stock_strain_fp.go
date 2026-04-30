package resolver

import (
	"context"

	A "github.com/IBM/fp-go/v2/array"
	Eq "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
)

// ── shared input ─────────────────────────────────────────────────────────────

// listStrainsInput holds only caller-supplied values; it is immutable throughout
// the pipeline.
type listStrainsInput struct {
	stockClient pb.StockServiceClient
	annoClient  anno.TaggedAnnotationServiceClient
	gctx        context.Context
	cursor      *int
	limit       *int
	filter      *models.StrainListFilter
}

// ── branch 1: stock-service context ─────────────────────────────────────────

type stockStrainsContext struct {
	input          listStrainsInput
	resolvedCursor int64
	resolvedLimit  int64
	filterQuery    string
	collection     *pb.StrainCollection
}

// ── branch 2: bacterial annotation context ───────────────────────────────────

type bacterialStrainsContext struct {
	input          listStrainsInput
	resolvedCursor int64
	resolvedLimit  int64
	annotations    *anno.TaggedAnnotationCollection
	uniqueIDs      []string
	strainList     *pb.StrainList
}

// ── shared result helpers ────────────────────────────────────────────────────

func onStrainListError(err error) T.Tuple2[error, *models.StrainListWithCursor] {
	return T.MakeTuple2(err, &models.StrainListWithCursor{})
}

func onStrainListSuccess(
	data *models.StrainListWithCursor,
) T.Tuple2[error, *models.StrainListWithCursor] {
	return T.MakeTuple2[error](nil, data)
}

// ── predicate: is this a bacterial request? ──────────────────────────────────

// filterStrainTypeEq is a strict equality instance for models.StrainType.
// Centralising the comparison here ensures all strain-type checks use the
// same semantics and avoids scattered bare == comparisons.
var filterStrainTypeEq = Eq.FromStrictEquals[models.StrainType]()

// isBacterialStrainType returns true when st equals models.StrainTypeBacterial.
// It is the innermost predicate in the chain; all other predicates below are
// built by contrammapping this one outward toward listStrainsInput.
func isBacterialStrainType(st models.StrainType) bool {
	return filterStrainTypeEq.Equals(st, models.StrainTypeBacterial)
}

// strainTypeFromFilter extracts the StrainType field from a StrainListFilter.
// It is the contramap function that lifts isBacterialStrainType from
// Predicate[StrainType] to Predicate[*StrainListFilter].
func strainTypeFromFilter(f *models.StrainListFilter) models.StrainType {
	return f.StrainType
}

// filterFromInput extracts the filter pointer from a listStrainsInput.
// It is the contramap function that lifts isBacterialListFilter from
// Predicate[*StrainListFilter] to Predicate[listStrainsInput].
func filterFromInput(input listStrainsInput) *models.StrainListFilter {
	return input.filter
}

// isBacterialListFilter is a Predicate[*models.StrainListFilter] that checks
// whether the filter's StrainType is bacterial. Built by contrammapping
// isBacterialStrainType through strainTypeFromFilter.
// Note: this predicate is only safe to call on a non-nil pointer; nil safety
// is enforced by isBacterialFilter via P.IsNonZero and short-circuit &&.
var isBacterialListFilter = F.Pipe1(
	isBacterialStrainType,
	P.ContraMap(strainTypeFromFilter),
)

// isBacterialFilter is a Predicate[listStrainsInput] that returns true when
// the input carries a non-nil filter whose StrainType is bacterial.
//
// Pipeline:
//   - P.IsNonZero guards against a nil filter pointer; Go's short-circuit &&
//     inside P.And ensures isBacterialListFilter is never called on nil.
//   - P.And(isBacterialListFilter) applies the strain-type check only when
//     the pointer is confirmed non-nil.
//   - P.ContraMap(filterFromInput) lifts the whole predicate from
//     Predicate[*StrainListFilter] to Predicate[listStrainsInput].
var isBacterialFilter = F.Pipe2(
	P.IsNonZero[*models.StrainListFilter](),
	P.And(isBacterialListFilter),
	P.ContraMap(filterFromInput),
)

// ── stock-service pipeline ───────────────────────────────────────────────────

func runStockPipeline(
	input listStrainsInput,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return F.Pipe4(
		IOE.Of[error](stockStrainsContext{input: input}),
		IOE.Map[error](computeStockParams),
		IOE.Chain(buildStockFilterQuery),
		IOE.Chain(fetchStrainCollection),
		IOE.Map[error](buildStockStrainResult),
	)
}

var computeStockParams = func(ctx stockStrainsContext) stockStrainsContext {
	ctx.resolvedCursor = resolverutils.GetCursorFP(ctx.input.cursor)
	ctx.resolvedLimit = resolverutils.GetLimitFP(ctx.input.limit)
	return ctx
}

// buildStockFilterQuery lifts StrainFilterToQueryFP (which returns (string,error))
// into IOEither and attaches the result to the context.
func buildStockFilterQuery(
	ctx stockStrainsContext,
) IOE.IOEither[error, stockStrainsContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (string, error) {
			return resolverutils.StrainFilterToQueryFP(ctx.input.filter)
		}),
		IOE.Map[error](func(q string) stockStrainsContext {
			ctx.filterQuery = q
			return ctx
		}),
	)
}

func fetchStrainCollection(
	ctx stockStrainsContext,
) IOE.IOEither[error, stockStrainsContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.StrainCollection, error) {
			return ctx.input.stockClient.ListStrains(
				ctx.input.gctx,
				&pb.StockParameters{
					Cursor: ctx.resolvedCursor,
					Limit:  ctx.resolvedLimit,
					Filter: ctx.filterQuery,
				},
			)
		}),
		IOE.Map[error](func(coll *pb.StrainCollection) stockStrainsContext {
			ctx.collection = coll
			return ctx
		}),
	)
}

var convertStrainCollectionItem = func(item *pb.StrainCollection_Data) *models.Strain {
	return stock.ConvertToStrainModel(item.Id, item.Attributes)
}

var buildStockStrainResult = func(ctx stockStrainsContext) *models.StrainListWithCursor {
	meta := ctx.collection.Meta
	lmt := int(meta.Limit)
	return &models.StrainListWithCursor{
		Strains:        A.Map(convertStrainCollectionItem)(ctx.collection.Data),
		NextCursor:     int(meta.NextCursor),
		PreviousCursor: int(ctx.resolvedCursor),
		Limit:          &lmt,
		TotalCount:     int(meta.Total),
	}
}

// ── bacterial annotation pipeline ────────────────────────────────────────────

func runBacterialPipeline(
	input listStrainsInput,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return F.Pipe4(
		IOE.Of[error](bacterialStrainsContext{input: input}),
		IOE.Map[error](computeBacterialParams),
		IOE.Chain(fetchBacterialAnnotations),
		IOE.Map[error](deduplicateBacterialIDs),
		IOE.Chain(F.Ternary(
			isEmptyBacterialIDs,
			emptyBacterialResult,
			fetchAndBuildBacterialResult,
		)),
	)
}

// isEmptyBacterialIDs checks whether the deduplicated ID list is empty.
// When true, the pipeline short-circuits to an empty result without
// making a gRPC call.
var isEmptyBacterialIDs = F.Pipe2(
	P.IsNonZero[int](),
	P.Not,
	P.ContraMap(func(c bacterialStrainsContext) int {
		return len(c.uniqueIDs)
	}),
)

func emptyBacterialResult(
	_ bacterialStrainsContext,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return IOE.Of[error](&models.StrainListWithCursor{})
}

func fetchAndBuildBacterialResult(
	ctx bacterialStrainsContext,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*pb.StrainList, error) {
			return ctx.input.stockClient.ListStrainsByIds(
				ctx.input.gctx,
				&pb.StockIdList{Id: ctx.uniqueIDs},
			)
		}),
		IOE.Map[error](func(sl *pb.StrainList) bacterialStrainsContext {
			ctx.strainList = sl
			return ctx
		}),
		IOE.Map[error](func(c bacterialStrainsContext) *models.StrainListWithCursor {
			lmt := int(c.resolvedLimit)
			nextCursor := int(F.Pipe1(
				O.FromNillable(c.annotations.Meta),
				O.Fold(
					F.Constant[int64](0),
					func(m *anno.Meta) int64 { return m.NextCursor },
				),
			))
			return &models.StrainListWithCursor{
				Strains: A.Map(func(item *pb.StrainList_Data) *models.Strain {
					return stock.ConvertToStrainModel(item.Id, item.Attributes)
				})(c.strainList.Data),
				NextCursor:     nextCursor,
				PreviousCursor: int(c.resolvedCursor),
				Limit:          &lmt,
				TotalCount:     len(c.strainList.Data),
			}
		}),
	)
}

func computeBacterialParams(
	ctx bacterialStrainsContext,
) bacterialStrainsContext {
	ctx.resolvedCursor = resolverutils.GetCursorFP(ctx.input.cursor)
	ctx.resolvedLimit = resolverutils.GetLimitFP(ctx.input.limit)
	return ctx
}

func fetchBacterialAnnotations(
	ctx bacterialStrainsContext,
) IOE.IOEither[error, bacterialStrainsContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*anno.TaggedAnnotationCollection, error) {
			return ctx.input.annoClient.ListAnnotations(
				ctx.input.gctx,
				&anno.ListParameters{
					Cursor: ctx.resolvedCursor,
					Limit:  ctx.resolvedLimit,
					Filter: resolverutils.BacterialAnnotationFilter(),
				},
			)
		}),
		IOE.Map[error](func(
			coll *anno.TaggedAnnotationCollection,
		) bacterialStrainsContext {
			ctx.annotations = coll
			return ctx
		}),
	)
}

// deduplicateBacterialIDs extracts the EntryId from each annotation, removes
// duplicates, and sorts the resulting list. This ensures we only query the
// stock service once per unique strain ID and that the IDs are in a consistent
// order for testing and caching.
func deduplicateBacterialIDs(
	ctx bacterialStrainsContext,
) bacterialStrainsContext {
	ctx.uniqueIDs = F.Pipe3(
		ctx.annotations.Data,
		A.Map(func(d *anno.TaggedAnnotationCollection_Data) string {
			return d.Attributes.EntryId
		}),
		A.StrictUniq[string],
		A.Sort(S.Ord),
	)
	return ctx
}
