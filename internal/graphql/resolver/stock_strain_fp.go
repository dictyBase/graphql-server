package resolver

import (
	"context"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	Pa "github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/record"
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

var isBacterialFilter = func(input listStrainsInput) bool {
	return F.Pipe2(
		input.filter,
		O.FromNillable[models.StrainListFilter],
		O.Fold(
			F.Constant(false),
			func(f *models.StrainListFilter) bool {
				return f.StrainType == models.StrainTypeBacterial
			},
		),
	)
}

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
		IOE.Chain(shortCircuitOrFetchStrains),
	)
}

// shortCircuitOrFetchStrains returns an empty StrainListWithCursor immediately
// when uniqueIDs is empty, bypassing the gRPC call entirely. For non-empty ID
// lists it delegates to fetchAndBuildBacterialResult.
var shortCircuitOrFetchStrains = func(
	ctx bacterialStrainsContext,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return F.Ternary(
		func(c bacterialStrainsContext) bool { return len(c.uniqueIDs) == 0 },
		func(c bacterialStrainsContext) IOE.IOEither[error, *models.StrainListWithCursor] {
			return IOE.Of[error](&models.StrainListWithCursor{})
		},
		fetchAndBuildBacterialResult,
	)(ctx)
}

func fetchAndBuildBacterialResult(
	ctx bacterialStrainsContext,
) IOE.IOEither[error, *models.StrainListWithCursor] {
	return F.Pipe1(
		fetchStrainsByIDs(ctx),
		IOE.Map[error](buildBacterialStrainResult),
	)
}

var computeBacterialParams = func(
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
		IOE.Map[error](func(coll *anno.TaggedAnnotationCollection) bacterialStrainsContext {
			ctx.annotations = coll
			return ctx
		}),
	)
}

// deduplicateBacterialIDs extracts unique strain IDs from annotation data using
// R.FromEntries (builds map[string]struct{} keyed on EntryId) then R.Keys.
// The resulting slice is alphabetically sorted with A.Sort(S.Ord) for
// deterministic ordering. No imperative set/loop — purely functional.
var deduplicateBacterialIDs = func(
	ctx bacterialStrainsContext,
) bacterialStrainsContext {
	ctx.uniqueIDs = F.Pipe2(
		R.FromEntries[string, struct{}](
			A.Map(func(d *anno.TaggedAnnotationCollection_Data) R.Entry[string, struct{}] {
				return Pa.MakePair(d.Attributes.EntryId, struct{}{})
			})(ctx.annotations.Data),
		),
		R.Keys[string, struct{}],
		A.Sort(S.Ord),
	)
	return ctx
}

func fetchStrainsByIDs(
	ctx bacterialStrainsContext,
) IOE.IOEither[error, bacterialStrainsContext] {
	return F.Pipe1(
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
	)
}

var convertStrainListItem = func(item *pb.StrainList_Data) *models.Strain {
	return stock.ConvertToStrainModel(item.Id, item.Attributes)
}

var buildBacterialStrainResult = func(
	ctx bacterialStrainsContext,
) *models.StrainListWithCursor {
	lmt := int(ctx.resolvedLimit)
	nextCursor := int(F.Pipe1(
		O.FromNillable(ctx.annotations.Meta),
		O.Fold(
			F.Constant[int64](0),
			func(m *anno.Meta) int64 { return m.NextCursor },
		),
	))
	return &models.StrainListWithCursor{
		Strains:        A.Map(convertStrainListItem)(ctx.strainList.Data),
		NextCursor:     nextCursor,
		PreviousCursor: int(ctx.resolvedCursor),
		Limit:          &lmt,
		TotalCount:     len(ctx.strainList.Data),
	}
}
