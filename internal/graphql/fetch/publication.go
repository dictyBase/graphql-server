package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	arrutils "github.com/dictyBase/fp-go-loom/array"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/dictyBase/fp-go-loom/strutils"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/dictybase/literature"
)

const (
	RedisKey     = "PUBLICATION_KEY"
	pubMedSource = "MED"
)

var Status = "published"

// nilTime is the zero-valued optional publication date.
var nilTime *time.Time

// pubMedClient is the narrow seam over the NCBI eUtils client.
type pubMedClient interface {
	GetArticle(pmid string) (*literature.Article, error)
}

// europePMCClient is the narrow seam over the EuropePMC client.
type europePMCClient interface {
	GetArticle(pmid string) (*literature.EuropePMCArticle, error)
	GetArticleByDOI(doi string) (*literature.EuropePMCArticle, error)
}

// PublicationParams carries optional NCBI and EuropePMC request identity values.
type PublicationParams struct {
	APIKey string
	Tool   string
	Email  string
}

// Publication resolves PMIDs (NCBI first, EuropePMC fallback) and DOIs
// (EuropePMC) into GraphQL publication models.
type Publication struct {
	pubMed    pubMedClient
	europePMC europePMCClient
}

type clientBuildState struct {
	params    *PublicationParams
	pubMed    *literature.Client
	europePMC *literature.EuropePMCClient
}

type pubFetchState struct {
	identifier string
	repo       repository.Repository
	source     *Publication
	redisKey   string
	cached     bool
	rawCached  string
	pub        *models.Publication
}

var (
	buildPubMedLens = L.MakeLens(
		func(s clientBuildState) *literature.Client { return s.pubMed },
		func(s clientBuildState, v *literature.Client) clientBuildState {
			s.pubMed = v
			return s
		},
	)
	buildEuropeLens = L.MakeLens(
		func(s clientBuildState) *literature.EuropePMCClient {
			return s.europePMC
		},
		func(
			s clientBuildState,
			v *literature.EuropePMCClient,
		) clientBuildState {
			s.europePMC = v
			return s
		},
	)

	identifierLens = L.MakeLens(
		func(s pubFetchState) string { return s.identifier },
		func(s pubFetchState, v string) pubFetchState {
			s.identifier = v
			return s
		},
	)
	repoLens = L.MakeLens(
		func(s pubFetchState) repository.Repository { return s.repo },
		func(s pubFetchState, v repository.Repository) pubFetchState {
			s.repo = v
			return s
		},
	)
	sourceLens = L.MakeLens(
		func(s pubFetchState) *Publication { return s.source },
		func(s pubFetchState, v *Publication) pubFetchState {
			s.source = v
			return s
		},
	)
	redisKeyLens = L.MakeLens(
		func(s pubFetchState) string { return s.redisKey },
		func(s pubFetchState, v string) pubFetchState {
			s.redisKey = v
			return s
		},
	)
	cachedLens = L.MakeLens(
		func(s pubFetchState) bool { return s.cached },
		func(s pubFetchState, v bool) pubFetchState {
			s.cached = v
			return s
		},
	)
	rawCachedLens = L.MakeLens(
		func(s pubFetchState) string { return s.rawCached },
		func(s pubFetchState, v string) pubFetchState {
			s.rawCached = v
			return s
		},
	)
	pubLens = L.MakeLens(
		func(s pubFetchState) *models.Publication { return s.pub },
		func(s pubFetchState, v *models.Publication) pubFetchState {
			s.pub = v
			return s
		},
	)

	trimHTTPSDOI = F.Bind2nd(strings.TrimPrefix, "https://doi.org/")
	trimHTTPDOI  = F.Bind2nd(strings.TrimPrefix, "http://doi.org/")
	trimDOIColon = F.Bind2nd(strings.TrimPrefix, "doi:")

	pmidPattern   = regexp.MustCompile(`^[1-9][0-9]*$`)
	validatePMID  = E.FromPredicate(validPMID, invalidPMIDError)
	isNonZeroTime = P.Not(time.Time.IsZero)

	isDOIStr = F.Pipe1(
		S.Includes("/"),
		P.Or(S.HasPrefix("10.")),
	)
	isDOI = F.Pipe1(
		isDOIStr,
		P.ContraMap(identifierLens.Get),
	)
	resolveArticleFetch = F.Pipe1(
		isDOI,
		P.Fold(fetchByPMID, fetchEuropeByDOI),
	)
	validateIdentifierBranch = F.Pipe1(
		isDOI,
		P.Fold(validateStatePMID, E.Of[error, pubFetchState]),
	)

	hasRepo = F.Pipe1(
		P.IsNonZero[repository.Repository](),
		P.ContraMap(repoLens.Get),
	)
	branchByRepo = F.Pipe1(
		hasRepo,
		P.Fold(resolveArticleFetch, fetchWithCache),
	)
	branchByCache = F.Pipe1(
		cachedLens.Get,
		P.Fold(fetchAndStoreArticle, loadFromCache),
	)
)

func someOption[T any](
	apply func(string) T,
	value string,
) O.Option[T] {
	return F.Pipe2(
		value,
		O.FromPredicate(S.IsNonEmpty),
		O.Map(apply),
	)
}

func pubMedOptions(p *PublicationParams) []literature.Option {
	return F.Pipe1(
		[]O.Option[literature.Option]{
			someOption(literature.WithAPIKey, p.APIKey),
			someOption(literature.WithTool, p.Tool),
			someOption(literature.WithEmail, p.Email),
		},
		arrutils.Compact[literature.Option],
	)
}

func europePMCOptions(p *PublicationParams) []literature.EuropePMCOption {
	return F.Pipe1(
		[]O.Option[literature.EuropePMCOption]{
			someOption(literature.WithEuropePMCEmail, p.Email),
			someOption(literature.WithEuropePMCUserAgent, p.Tool),
		},
		arrutils.Compact[literature.EuropePMCOption],
	)
}

func createPubMedFromState(
	s clientBuildState,
) IOE.IOEither[error, *literature.Client] {
	opts := pubMedOptions(s.params)
	newClient := func() (*literature.Client, error) {
		return literature.New(opts...)
	}
	return F.Pipe1(
		IOE.TryCatchError(newClient),
		IOE.MapLeft[*literature.Client](func(err error) error {
			return fmt.Errorf(
				"error creating literature pubmed client: %w",
				err,
			)
		}),
	)
}

func createEuropeFromState(
	s clientBuildState,
) IOE.IOEither[error, *literature.EuropePMCClient] {
	opts := europePMCOptions(s.params)
	newClient := func() (*literature.EuropePMCClient, error) {
		return literature.NewEuropePMCClient(opts...)
	}
	return F.Pipe1(
		IOE.TryCatchError(newClient),
		IOE.MapLeft[*literature.EuropePMCClient](func(err error) error {
			return fmt.Errorf(
				"error creating literature europepmc client: %w",
				err,
			)
		}),
	)
}

func assemblePublication(s clientBuildState) *Publication {
	return &Publication{
		pubMed:    s.pubMed,
		europePMC: s.europePMC,
	}
}

// NewPublication builds the publication source with both NCBI PubMed and
// EuropePMC clients.
func NewPublication(
	p *PublicationParams,
) IOE.IOEither[error, *Publication] {
	return F.Pipe3(
		IOE.Of[error](clientBuildState{params: p}),
		IOE.Bind(buildPubMedLens.Set, createPubMedFromState),
		IOE.Bind(buildEuropeLens.Set, createEuropeFromState),
		IOE.Map[error](assemblePublication),
	)
}

func validPMID(pmid string) bool {
	return pmidPattern.MatchString(pmid)
}

func invalidPMIDError(pmid string) error {
	return fmt.Errorf(
		"error fetching publication %s: %w",
		pmid,
		literature.NewErrorWithPMID(
			literature.ErrorTypeInvalidInput,
			pmid,
			"pmid must contain decimal digits",
		),
	)
}

func toPBAuthor(idx int, a literature.Author) *pb.Author {
	return &pb.Author{
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Rank:      int64(idx),
	}
}

func toPBEuropeAuthor(idx int, a literature.EuropePMCAuthor) *pb.Author {
	return &pb.Author{
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Initials:  a.Initials,
		Rank:      int64(idx),
	}
}

func pubMedArticleToPublication(a *literature.Article) *models.Publication {
	emptyIssn := ""
	toTimePtr := func(t time.Time) *time.Time { return &t }
	authors := F.Pipe1(a.Authors, A.MapWithIndex(toPBAuthor))
	pubDate := F.Pipe3(
		a.PublishDate,
		O.FromPredicate(isNonZeroTime),
		O.Map(toTimePtr),
		O.GetOrElse(F.Constant(nilTime)),
	)

	return &models.Publication{
		ID:       a.PMID,
		Title:    a.Title,
		Abstract: a.Abstract,
		Journal:  a.Journal,
		PubDate:  pubDate,
		Doi:      &a.DOI,
		Volume:   &a.Volume,
		Issue:    &a.Issue,
		Pages:    &a.Pages,
		Issn:     &emptyIssn,
		PubType:  "",
		Source:   pubMedSource,
		Status:   &Status,
		Authors:  authors,
	}
}

func europeArticleToPublication(
	a *literature.EuropePMCArticle,
) *models.Publication {
	authors := F.Pipe1(a.Authors, A.MapWithIndex(toPBEuropeAuthor))
	pubType := F.Pipe2(
		a.PubTypes,
		A.Head[string],
		O.GetOrElse(F.Constant("")),
	)
	id := F.Pipe2(
		a.PMID,
		O.FromPredicate(S.IsNonEmpty),
		O.GetOrElse(F.Constant(a.ID)),
	)

	return &models.Publication{
		ID:       id,
		Title:    a.Title,
		Abstract: a.Abstract,
		Journal:  a.Journal.Title,
		PubDate:  a.PublishDate,
		Doi:      &a.DOI,
		Volume:   &a.Journal.Volume,
		Issue:    &a.Journal.Issue,
		Pages:    &a.PageInfo,
		Issn:     &a.Journal.ISSN,
		PubType:  pubType,
		Source:   a.Source,
		Status:   &Status,
		Authors:  authors,
	}
}

func validateStatePMID(s pubFetchState) E.Either[error, pubFetchState] {
	return F.Pipe1(
		validatePMID(s.identifier),
		E.Map[error](F.Constant1[string](s)),
	)
}

func callPubMedByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	rawFetch := func() (*literature.Article, error) {
		return s.source.pubMed.GetArticle(s.identifier)
	}
	return F.Pipe2(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*literature.Article](func(err error) error {
			return fmt.Errorf(
				"error fetching publication %s: %w",
				s.identifier,
				err,
			)
		}),
		IOE.Map[error](pubMedArticleToPublication),
	)
}

func callEuropeByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	rawFetch := func() (*literature.EuropePMCArticle, error) {
		return s.source.europePMC.GetArticle(s.identifier)
	}
	return F.Pipe2(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*literature.EuropePMCArticle](func(err error) error {
			return fmt.Errorf(
				"error fetching publication %s: %w",
				s.identifier,
				err,
			)
		}),
		IOE.Map[error](europeArticleToPublication),
	)
}

func fetchByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	fallback := func() IOE.IOEither[error, *models.Publication] {
		return callEuropeByPMID(s)
	}
	return F.Pipe2(
		s,
		callPubMedByPMID,
		IOE.Alt(fallback),
	)
}

func fetchEuropeByDOI(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	rawFetch := func() (*literature.EuropePMCArticle, error) {
		return s.source.europePMC.GetArticleByDOI(s.identifier)
	}
	return F.Pipe2(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*literature.EuropePMCArticle](func(err error) error {
			return fmt.Errorf("error fetching doi %s: %w", s.identifier, err)
		}),
		IOE.Map[error](europeArticleToPublication),
	)
}

func deriveRedisKey(s pubFetchState) string {
	prefix := strutils.JoinStrings(RedisKey, "/")
	return strutils.JoinStrings(prefix, s.identifier)
}

func checkCacheExists(s pubFetchState) IOE.IOEither[error, bool] {
	keyExists := func() (bool, error) {
		return s.repo.Exists(s.redisKey)
	}
	return F.Pipe1(
		IOE.TryCatchError(keyExists),
		IOE.MapLeft[bool](func(err error) error {
			return fmt.Errorf("error in checking for key %w", err)
		}),
	)
}

func unmarshalPublication(
	raw string,
) E.Either[error, *models.Publication] {
	pub := &models.Publication{}
	err := json.Unmarshal([]byte(raw), pub)
	return F.Pipe1(
		E.TryCatchError(pub, err),
		E.MapLeft[*models.Publication](func(err error) error {
			return fmt.Errorf("error in decoding json %w", err)
		}),
	)
}

func loadFromCache(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	rawGet := func() (string, error) {
		return s.repo.Get(s.redisKey)
	}
	return F.Pipe2(
		IOE.TryCatchError(rawGet),
		IOE.MapLeft[string](func(err error) error {
			return fmt.Errorf("error in getting existing key %w", err)
		}),
		IOE.ChainEitherK(unmarshalPublication),
	)
}

func encodeStatePublication(s pubFetchState) E.Either[error, pubFetchState] {
	setRaw := F.Pipe1(s, F.Flip(rawCachedLens.Set))
	b, err := json.Marshal(s.pub)
	return F.Pipe3(
		E.TryCatchError(b, err),
		E.MapLeft[[]byte](func(err error) error {
			return fmt.Errorf("error in converting json to byte %w", err)
		}),
		E.Map[error](func(b []byte) string { return string(b) }),
		E.Map[error](setRaw),
	)
}

func storeCachedString(s pubFetchState) IOE.IOEither[error, string] {
	rawSet := func() (bool, error) {
		return true, s.repo.Set(s.redisKey, s.rawCached)
	}
	return F.Pipe2(
		IOE.TryCatchError(rawSet),
		IOE.MapLeft[bool](func(err error) error {
			return fmt.Errorf("error in setting key in redis %w", err)
		}),
		IOE.Map[error](F.Constant1[bool](s.rawCached)),
	)
}

func fetchAndStoreArticle(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe4(
		IOE.Of[error](s),
		IOE.Bind(pubLens.Set, resolveArticleFetch),
		IOE.ChainEitherK(encodeStatePublication),
		IOE.ChainFirst(storeCachedString),
		IOE.Map[error](pubLens.Get),
	)
}

func fetchWithCache(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe2(
		IOE.Of[error](s),
		IOE.Bind(cachedLens.Set, checkCacheExists),
		IOE.Chain(branchByCache),
	)
}

// FetchCachedArticle normalizes and validates id, checks Redis when repo is
// non-nil, and fetches from NCBI/EuropePMC on a cache miss.
func (p *Publication) FetchCachedArticle(
	repo repository.Repository,
	id string,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe7(
		pubFetchState{
			identifier: id,
			repo:       repo,
			source:     p,
		},
		identifierLens.Modify(trimHTTPSDOI),
		identifierLens.Modify(trimHTTPDOI),
		identifierLens.Modify(trimDOIColon),
		IOE.Of[error, pubFetchState],
		IOE.ChainEitherK(validateIdentifierBranch),
		IOE.Let[error](redisKeyLens.Set, deriveRedisKey),
		IOE.Chain(branchByRepo),
	)
}

func onPublicationError(
	err error,
) T.Tuple2[error, *models.Publication] {
	return T.MakeTuple2[error, *models.Publication](err, nil)
}

func onPublicationSuccess(
	pub *models.Publication,
) T.Tuple2[error, *models.Publication] {
	return T.MakeTuple2[error](nil, pub)
}

// FetchPublication fetches a publication by PMID or DOI using Redis caching
// and the literature NCBI and EuropePMC clients.
func FetchPublication(
	_ context.Context,
	repo repository.Repository,
	id string,
) (*models.Publication, error) {
	newSource := NewPublication(&PublicationParams{})
	result := F.Pipe10(
		pubFetchState{
			identifier: id,
			repo:       repo,
		},
		identifierLens.Modify(trimHTTPSDOI),
		identifierLens.Modify(trimHTTPDOI),
		identifierLens.Modify(trimDOIColon),
		IOE.Of[error, pubFetchState],
		IOE.ChainEitherK(validateIdentifierBranch),
		IOE.Bind(sourceLens.Set, F.Constant1[pubFetchState](newSource)),
		IOE.Let[error](redisKeyLens.Set, deriveRedisKey),
		IOE.Chain(branchByRepo),
		ioeutils.ToEither[error, *models.Publication],
		E.Fold(onPublicationError, onPublicationSuccess),
	)
	return result.F2, result.F1
}
