package stock

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	T "github.com/IBM/fp-go/v2/tuple"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
)

type stockPubsState struct {
	ctx        context.Context
	redis      repository.Repository
	stockID    string
	logger     *logrus.Entry
	ids        []string
	resolvable []string
	skipped    []string
}

var (
	stockPubsCtxLens = L.MakeLens(
		func(s stockPubsState) context.Context { return s.ctx },
		func(s stockPubsState, v context.Context) stockPubsState {
			s.ctx = v
			return s
		},
	)
	stockPubsRedisLens = L.MakeLens(
		func(s stockPubsState) repository.Repository { return s.redis },
		func(s stockPubsState, v repository.Repository) stockPubsState {
			s.redis = v
			return s
		},
	)
	stockPubsIDLens = L.MakeLens(
		func(s stockPubsState) string { return s.stockID },
		func(s stockPubsState, v string) stockPubsState {
			s.stockID = v
			return s
		},
	)
	stockPubsLoggerLens = L.MakeLens(
		func(s stockPubsState) *logrus.Entry { return s.logger },
		func(s stockPubsState, v *logrus.Entry) stockPubsState {
			s.logger = v
			return s
		},
	)
	stockPubsRawIDsLens = L.MakeLens(
		func(s stockPubsState) []string { return s.ids },
		func(s stockPubsState, v []string) stockPubsState {
			s.ids = v
			return s
		},
	)
	stockPubsResolvableLens = L.MakeLens(
		func(s stockPubsState) []string { return s.resolvable },
		func(s stockPubsState, v []string) stockPubsState {
			s.resolvable = v
			return s
		},
	)
	stockPubsSkippedLens = L.MakeLens(
		func(s stockPubsState) []string { return s.skipped },
		func(s stockPubsState, v []string) stockPubsState {
			s.skipped = v
			return s
		},
	)

	phenoNameLens = L.MakeLensRef(
		func(m *models.Phenotype) string { return m.Phenotype },
		func(m *models.Phenotype, v string) *models.Phenotype {
			m.Phenotype = v
			return m
		},
	)
	phenoEnvLens = L.MakeLensRef(
		func(m *models.Phenotype) *string { return m.Environment },
		func(m *models.Phenotype, v *string) *models.Phenotype {
			m.Environment = v
			return m
		},
	)
	phenoAssayLens = L.MakeLensRef(
		func(m *models.Phenotype) *string { return m.Assay },
		func(m *models.Phenotype, v *string) *models.Phenotype {
			m.Assay = v
			return m
		},
	)
	phenoNoteLens = L.MakeLensRef(
		func(m *models.Phenotype) *string { return m.Note },
		func(m *models.Phenotype, v *string) *models.Phenotype {
			m.Note = v
			return m
		},
	)
	phenoPubLens = L.MakeLensRef(
		func(m *models.Phenotype) *models.Publication {
			return m.Publication
		},
		func(
			m *models.Phenotype,
			v *models.Publication,
		) *models.Phenotype {
			m.Publication = v
			return m
		},
	)

	noopIO              = IO.Of(F.VOID)
	nilPublication      *models.Publication
	filterResolvableIDs = A.Filter(fetch.IsResolvableIdentifier)
	filterSkippedIDs    = A.Filter(P.Not(fetch.IsResolvableIdentifier))
	deriveResolvableIDs = F.Flow2(
		stockPubsRawIDsLens.Get,
		filterResolvableIDs,
	)
	deriveSkippedIDs = F.Flow2(
		stockPubsRawIDsLens.Get,
		filterSkippedIDs,
	)
	hasSkippedStockPubs = F.Pipe1(
		predarrays.IsNonEmpty[string](),
		P.ContraMap(stockPubsSkippedLens.Get),
	)
	warnSkippedStockPubs = F.Pipe1(
		hasSkippedStockPubs,
		P.Fold(
			F.Constant1[stockPubsState](noopIO),
			logSkippedStockPubsIO,
		),
	)
	nonNilPublications = F.Flow2(
		O.FromPredicate(predarrays.IsNonEmpty[*models.Publication]()),
		O.GetOrElse(emptyPublications),
	)
	nonNilPhenotypes = F.Flow2(
		O.FromPredicate(predarrays.IsNonEmpty[*models.Phenotype]()),
		O.GetOrElse(emptyPhenotypes),
	)

	isPhenoOnto = F.Pipe1(
		predord.StrEq(registry.PhenoOntology),
		P.ContraMap(annoOntology),
	)
	isEnvOnto = F.Pipe1(
		predord.StrEq(registry.EnvOntology),
		P.ContraMap(annoOntology),
	)
	isAssayOnto = F.Pipe1(
		predord.StrEq(registry.AssayOntology),
		P.ContraMap(annoOntology),
	)
	isDictyAnnoOnto = F.Pipe1(
		predord.StrEq(registry.DictyAnnoOntology),
		P.ContraMap(annoOntology),
	)
	isLiteratureTag = F.Pipe1(
		predord.StrEq(registry.LiteratureTag),
		P.ContraMap(annoTag),
	)
	isNoteTag = F.Pipe1(
		predord.StrEq(registry.NoteTag),
		P.ContraMap(annoTag),
	)
	isLiteratureAnno = F.Pipe1(
		isDictyAnnoOnto,
		P.And(isLiteratureTag),
	)
	isNoteAnno = F.Pipe1(
		isDictyAnnoOnto,
		P.And(isNoteTag),
	)
	isResolvableAnnoValue = F.Pipe1(
		fetch.IsResolvableIdentifier,
		P.ContraMap(annoValue),
	)
	isResolvableLiteratureAnno = F.Pipe1(
		isLiteratureAnno,
		P.And(isResolvableAnnoValue),
	)
)

func emptyPublications() []*models.Publication {
	return make([]*models.Publication, 0)
}

func emptyPhenotypes() []*models.Phenotype {
	return make([]*models.Phenotype, 0)
}

func newStockPubsState(
	ctx context.Context,
	redis repository.Repository,
	stockID string,
	logger *logrus.Entry,
	ids []string,
) stockPubsState {
	return F.Pipe5(
		stockPubsState{},
		stockPubsCtxLens.Set(ctx),
		stockPubsRedisLens.Set(redis),
		stockPubsIDLens.Set(stockID),
		stockPubsLoggerLens.Set(logger),
		stockPubsRawIDsLens.Set(ids),
	)
}

func logSkippedStockPubsIO(s stockPubsState) IO.IO[IO.Void] {
	return IO.FromImpure(func() {
		stockPubsLoggerLens.Get(s).Warnf(
			"skipping unresolvable publication identifiers %v for stock %s",
			stockPubsSkippedLens.Get(s),
			stockPubsIDLens.Get(s),
		)
	})
}

func wrapStockPubsError(s stockPubsState, err error) error {
	return fmt.Errorf(
		"error in fetching publication for stock %s: %w",
		stockPubsIDLens.Get(s),
		err,
	)
}

func recordStockPubsError(
	s stockPubsState,
	err error,
) IOE.IOEither[error, IO.Void] {
	return F.Pipe1(
		IO.FromImpure(func() {
			errorutils.AddGQLError(s.ctx, err)
			s.logger.Error(err)
		}),
		IOE.FromIO[error, IO.Void],
	)
}

func onPublicationsError(err error) T.Tuple2[error, []*models.Publication] {
	return T.MakeTuple2(err, emptyPublications())
}

func onPublicationsSuccess(
	pubs []*models.Publication,
) T.Tuple2[error, []*models.Publication] {
	return T.MakeTuple2[error](nil, pubs)
}

// fetchStockPublications resolves the publications field for a stock item.
// It skips stored identifiers that are neither a decimal PMID nor a DOI,
// logs a warning when identifiers are skipped, and fails only on a real
// fetch error.
func fetchStockPublications(
	s stockPubsState,
) ([]*models.Publication, error) {
	redis := stockPubsRedisLens.Get(s)
	result := F.Pipe10(
		IOE.Of[error](s),
		IOE.Let[error](stockPubsResolvableLens.Set, deriveResolvableIDs),
		IOE.Let[error](stockPubsSkippedLens.Set, deriveSkippedIDs),
		IOE.ChainFirstIOK[error](warnSkippedStockPubs),
		IOE.Map[error](stockPubsResolvableLens.Get),
		IOE.Chain(
			IOE.TraverseArraySeq(
				F.Bind1st(fetch.FetchPublicationIOE, redis),
			),
		),
		IOE.Map[error](nonNilPublications),
		IOE.MapLeft[[]*models.Publication](
			F.Bind1st(wrapStockPubsError, s),
		),
		IOE.TapLeft[[]*models.Publication](
			F.Bind1st(recordStockPubsError, s),
		),
		ioeutils.ToEither[error, []*models.Publication],
		E.Fold(onPublicationsError, onPublicationsSuccess),
	)
	return result.F2, result.F1
}

func annoOntology(anno *annotation.TaggedAnnotationGroup_Data) string {
	return anno.GetAttributes().GetOntology()
}

func annoTag(anno *annotation.TaggedAnnotationGroup_Data) string {
	return anno.GetAttributes().GetTag()
}

func annoValue(anno *annotation.TaggedAnnotationGroup_Data) string {
	return anno.GetAttributes().GetValue()
}

func setPhenoTag(
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	tag := annoTag(anno)
	return phenoNameLens.Set(tag)
}

func setEnvTag(
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	tag := annoTag(anno)
	return phenoEnvLens.Set(&tag)
}

func setAssayTag(
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	tag := annoTag(anno)
	return phenoAssayLens.Set(&tag)
}

func setNoteValue(
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	val := annoValue(anno)
	return phenoNoteLens.Set(&val)
}

func setFetchedPhenoPub(
	s stockPubsState,
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	pub := F.Pipe3(
		fetch.FetchPublicationIOE(
			stockPubsRedisLens.Get(s),
			annoValue(anno),
		),
		IOE.TapLeft[*models.Publication](
			F.Bind1st(recordStockPubsError, s),
		),
		ioeutils.ToEither[error, *models.Publication],
		E.Fold(
			F.Constant1[error](nilPublication),
			F.Identity[*models.Publication],
		),
	)
	return phenoPubLens.Set(pub)
}

func skipUnresolvablePhenoPub(
	s stockPubsState,
	anno *annotation.TaggedAnnotationGroup_Data,
) func(*models.Phenotype) *models.Phenotype {
	stockPubsLoggerLens.Get(s).Warnf(
		"skipping unresolvable phenotype publication identifier %q",
		annoValue(anno),
	)
	return F.Identity[*models.Phenotype]
}

func phenotypeAnnoCases(
	s stockPubsState,
) []O.Kleisli[*annotation.TaggedAnnotationGroup_Data, func(*models.Phenotype) *models.Phenotype] {
	return []O.Kleisli[*annotation.TaggedAnnotationGroup_Data, func(*models.Phenotype) *models.Phenotype]{
		MO.Case(isPhenoOnto, setPhenoTag),
		MO.Case(isEnvOnto, setEnvTag),
		MO.Case(isAssayOnto, setAssayTag),
		MO.Case(isNoteAnno, setNoteValue),
		MO.Case(
			isResolvableLiteratureAnno,
			F.Bind1st(setFetchedPhenoPub, s),
		),
		MO.Case(
			isLiteratureAnno,
			F.Bind1st(skipUnresolvablePhenoPub, s),
		),
	}
}

func applyPhenotypeStep(
	m *models.Phenotype,
	step func(*models.Phenotype) *models.Phenotype,
) *models.Phenotype {
	return step(m)
}

func buildPhenotypeFromGroup(
	s stockPubsState,
	item *annotation.TaggedAnnotationGroupCollection_Data,
) *models.Phenotype {
	return F.Pipe2(
		item.GetGroup().GetData(),
		A.Map(
			MO.Match(
				F.Identity[*models.Phenotype],
				phenotypeAnnoCases(s),
			),
		),
		A.Reduce(applyPhenotypeStep, &models.Phenotype{}),
	)
}

func getPhenotypes(
	ctx context.Context,
	r *StrainResolver,
	data []*annotation.TaggedAnnotationGroupCollection_Data,
) []*models.Phenotype {
	state := newStockPubsState(
		ctx,
		r.Registry.GetRedisRepository(cache.RedisKey),
		"",
		r.Logger,
		nil,
	)
	return F.Pipe2(
		data,
		A.Map(F.Bind1st(buildPhenotypeFromGroup, state)),
		nonNilPhenotypes,
	)
}
