package fetch

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	B "github.com/IBM/fp-go/v2/bytes"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	J "github.com/IBM/fp-go/v2/json"
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
	buildParamsLens = L.MakeLens(
		func(s clientBuildState) *PublicationParams { return s.params },
		func(s clientBuildState, v *PublicationParams) clientBuildState {
			s.params = v
			return s
		},
	)
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
	pubMedClientLens = L.MakeLensRef(
		func(p *Publication) pubMedClient { return p.pubMed },
		func(p *Publication, v pubMedClient) *Publication {
			p.pubMed = v
			return p
		},
	)
	europePMCClientLens = L.MakeLensRef(
		func(p *Publication) europePMCClient { return p.europePMC },
		func(p *Publication, v europePMCClient) *Publication {
			p.europePMC = v
			return p
		},
	)
	authorFirstNameLens = L.MakeLensRef(
		func(a *pb.Author) string { return a.FirstName },
		func(a *pb.Author, v string) *pb.Author {
			a.FirstName = v
			return a
		},
	)
	authorLastNameLens = L.MakeLensRef(
		func(a *pb.Author) string { return a.LastName },
		func(a *pb.Author, v string) *pb.Author {
			a.LastName = v
			return a
		},
	)
	authorInitialsLens = L.MakeLensRef(
		func(a *pb.Author) string { return a.Initials },
		func(a *pb.Author, v string) *pb.Author {
			a.Initials = v
			return a
		},
	)
	authorRankLens = L.MakeLensRef(
		func(a *pb.Author) int64 { return a.Rank },
		func(a *pb.Author, v int64) *pb.Author {
			a.Rank = v
			return a
		},
	)
	pubModelIDLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.ID },
		func(p *models.Publication, v string) *models.Publication {
			p.ID = v
			return p
		},
	)
	pubModelTitleLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.Title },
		func(p *models.Publication, v string) *models.Publication {
			p.Title = v
			return p
		},
	)
	pubModelAbstractLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.Abstract },
		func(p *models.Publication, v string) *models.Publication {
			p.Abstract = v
			return p
		},
	)
	pubModelJournalLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.Journal },
		func(p *models.Publication, v string) *models.Publication {
			p.Journal = v
			return p
		},
	)
	pubModelPubDateLens = L.MakeLensRef(
		func(p *models.Publication) *time.Time { return p.PubDate },
		func(p *models.Publication, v *time.Time) *models.Publication {
			p.PubDate = v
			return p
		},
	)
	pubModelDoiLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Doi },
		func(p *models.Publication, v *string) *models.Publication {
			p.Doi = v
			return p
		},
	)
	pubModelVolumeLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Volume },
		func(p *models.Publication, v *string) *models.Publication {
			p.Volume = v
			return p
		},
	)
	pubModelIssueLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Issue },
		func(p *models.Publication, v *string) *models.Publication {
			p.Issue = v
			return p
		},
	)
	pubModelPagesLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Pages },
		func(p *models.Publication, v *string) *models.Publication {
			p.Pages = v
			return p
		},
	)
	pubModelIssnLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Issn },
		func(p *models.Publication, v *string) *models.Publication {
			p.Issn = v
			return p
		},
	)
	pubModelPubTypeLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.PubType },
		func(p *models.Publication, v string) *models.Publication {
			p.PubType = v
			return p
		},
	)
	pubModelSourceLens = L.MakeLensRef(
		func(p *models.Publication) string { return p.Source },
		func(p *models.Publication, v string) *models.Publication {
			p.Source = v
			return p
		},
	)
	pubModelStatusLens = L.MakeLensRef(
		func(p *models.Publication) *string { return p.Status },
		func(p *models.Publication, v *string) *models.Publication {
			p.Status = v
			return p
		},
	)
	pubModelAuthorsLens = L.MakeLensRef(
		func(p *models.Publication) []*pb.Author { return p.Authors },
		func(p *models.Publication, v []*pb.Author) *models.Publication {
			p.Authors = v
			return p
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
	normalizeIdentifier = F.Flow3(
		trimHTTPSDOI,
		trimHTTPDOI,
		trimDOIColon,
	)
	isPMIDOrDOI = F.Pipe1(
		isDOIStr,
		P.Or(validPMID),
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
	decodeCachedPub = F.Flow3(
		S.ToBytes,
		J.Unmarshal[*models.Publication],
		E.MapLeft[*models.Publication](wrapDecodeJSONError),
	)
	marshalCachedPub = F.Flow4(
		pubLens.Get,
		J.Marshal[*models.Publication],
		E.MapLeft[[]byte](wrapEncodeJSONError),
		E.Map[error](B.ToString),
	)
	encodeStatePublication = F.Flow2(
		E.Of[error, pubFetchState],
		E.Bind(rawCachedLens.Set, marshalCachedPub),
	)
)

// IsResolvableIdentifier reports whether the publication pipeline can
// resolve the stored identifier: a decimal PMID or a DOI-shaped string
// after DOI prefix normalization.
func IsResolvableIdentifier(id string) bool {
	return F.Pipe2(
		id,
		normalizeIdentifier,
		isPMIDOrDOI,
	)
}

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
	opts := F.Pipe2(s, buildParamsLens.Get, pubMedOptions)

	return F.Pipe1(
		IOE.TryCatchError(func() (*literature.Client, error) {
			return literature.New(opts...)
		}),
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
	opts := F.Pipe2(s, buildParamsLens.Get, europePMCOptions)

	return F.Pipe1(
		IOE.TryCatchError(func() (*literature.EuropePMCClient, error) {
			return literature.NewEuropePMCClient(opts...)
		}),
		IOE.MapLeft[*literature.EuropePMCClient](func(err error) error {
			return fmt.Errorf(
				"error creating literature europepmc client: %w",
				err,
			)
		}),
	)
}

func assemblePublication(s clientBuildState) *Publication {
	pubMed := buildPubMedLens.Get(s)
	europe := buildEuropeLens.Get(s)
	return F.Pipe2(
		&Publication{},
		pubMedClientLens.Set(pubMed),
		europePMCClientLens.Set(europe),
	)
}

// NewPublication builds the publication source with both NCBI PubMed and
// EuropePMC clients.
func NewPublication(
	p *PublicationParams,
) IOE.IOEither[error, *Publication] {
	return F.Pipe5(
		clientBuildState{},
		buildParamsLens.Set(p),
		IOE.Of[error, clientBuildState],
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
	rank := int64(idx)
	return F.Pipe3(
		&pb.Author{},
		authorFirstNameLens.Set(a.FirstName),
		authorLastNameLens.Set(a.LastName),
		authorRankLens.Set(rank),
	)
}

func toPBEuropeAuthor(idx int, a literature.EuropePMCAuthor) *pb.Author {
	rank := int64(idx)
	return F.Pipe4(
		&pb.Author{},
		authorFirstNameLens.Set(a.FirstName),
		authorLastNameLens.Set(a.LastName),
		authorInitialsLens.Set(a.Initials),
		authorRankLens.Set(rank),
	)
}

func pubMedArticleToPublication(a *literature.Article) *models.Publication {
	emptyIssn := ""
	authors := F.Pipe1(a.Authors, A.MapWithIndex(toPBAuthor))
	pubDate := F.Pipe3(
		a.PublishDate,
		O.FromPredicate(isNonZeroTime),
		O.Map(func(t time.Time) *time.Time { return new(t) }),
		O.GetOrElse(F.Constant(nilTime)),
	)

	return F.Pipe14(
		&models.Publication{},
		pubModelIDLens.Set(a.PMID),
		pubModelTitleLens.Set(a.Title),
		pubModelAbstractLens.Set(a.Abstract),
		pubModelJournalLens.Set(a.Journal),
		pubModelPubDateLens.Set(pubDate),
		pubModelDoiLens.Set(&a.DOI),
		pubModelVolumeLens.Set(&a.Volume),
		pubModelIssueLens.Set(&a.Issue),
		pubModelPagesLens.Set(&a.Pages),
		pubModelIssnLens.Set(&emptyIssn),
		pubModelPubTypeLens.Set(""),
		pubModelSourceLens.Set(pubMedSource),
		pubModelStatusLens.Set(&Status),
		pubModelAuthorsLens.Set(authors),
	)
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

	return F.Pipe14(
		&models.Publication{},
		pubModelIDLens.Set(id),
		pubModelTitleLens.Set(a.Title),
		pubModelAbstractLens.Set(a.Abstract),
		pubModelJournalLens.Set(a.Journal.Title),
		pubModelPubDateLens.Set(a.PublishDate),
		pubModelDoiLens.Set(&a.DOI),
		pubModelVolumeLens.Set(&a.Journal.Volume),
		pubModelIssueLens.Set(&a.Journal.Issue),
		pubModelPagesLens.Set(&a.PageInfo),
		pubModelIssnLens.Set(&a.Journal.ISSN),
		pubModelPubTypeLens.Set(pubType),
		pubModelSourceLens.Set(a.Source),
		pubModelStatusLens.Set(&Status),
		pubModelAuthorsLens.Set(authors),
	)
}

func validateStatePMID(s pubFetchState) E.Either[error, pubFetchState] {
	return F.Pipe2(
		identifierLens.Get(s),
		validatePMID,
		E.Map[error](F.Constant1[string](s)),
	)
}

func callPubMedByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*literature.Article, error) {
			return s.source.pubMed.GetArticle(s.identifier)
		}),
		IOE.MapLeft[*literature.Article](func(err error) error {
			return fmt.Errorf(
				"error fetching publication %s: %w",
				identifierLens.Get(s),
				err,
			)
		}),
		IOE.Map[error](pubMedArticleToPublication),
	)
}

func callEuropeByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*literature.EuropePMCArticle, error) {
			return s.source.europePMC.GetArticle(s.identifier)
		}),
		IOE.MapLeft[*literature.EuropePMCArticle](func(err error) error {
			return fmt.Errorf(
				"error fetching publication %s: %w",
				identifierLens.Get(s),
				err,
			)
		}),
		IOE.Map[error](europeArticleToPublication),
	)
}

func fetchByPMID(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	europeAttempt := callEuropeByPMID(s)
	return F.Pipe2(
		s,
		callPubMedByPMID,
		IOE.OrElse(F.Constant1[error](europeAttempt)),
	)
}

func fetchEuropeByDOI(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*literature.EuropePMCArticle, error) {
			return s.source.europePMC.GetArticleByDOI(s.identifier)
		}),
		IOE.MapLeft[*literature.EuropePMCArticle](func(err error) error {
			return fmt.Errorf(
				"error fetching doi %s: %w",
				identifierLens.Get(s),
				err,
			)
		}),
		IOE.Map[error](europeArticleToPublication),
	)
}

func deriveRedisKey(s pubFetchState) string {
	prefix := strutils.JoinStrings(RedisKey, "/")
	return strutils.JoinStrings(prefix, identifierLens.Get(s))
}

func checkCacheExists(s pubFetchState) IOE.IOEither[error, bool] {
	return F.Pipe1(
		IOE.TryCatchError(func() (bool, error) {
			return s.repo.Exists(s.redisKey)
		}),
		IOE.MapLeft[bool](func(err error) error {
			return fmt.Errorf("error in checking for key %w", err)
		}),
	)
}

func wrapDecodeJSONError(err error) error {
	return fmt.Errorf("error in decoding json %w", err)
}

func wrapEncodeJSONError(err error) error {
	return fmt.Errorf("error in converting json to byte %w", err)
}

func loadFromCache(
	s pubFetchState,
) IOE.IOEither[error, *models.Publication] {
	return F.Pipe2(
		IOE.TryCatchError(func() (string, error) {
			return s.repo.Get(s.redisKey)
		}),
		IOE.MapLeft[string](func(err error) error {
			return fmt.Errorf("error in getting existing key %w", err)
		}),
		IOE.ChainEitherK(decodeCachedPub),
	)
}

func storeCachedString(s pubFetchState) IOE.IOEither[error, string] {
	return F.Pipe2(
		IOE.TryCatchError(func() (bool, error) {
			return true, s.repo.Set(s.redisKey, s.rawCached)
		}),
		IOE.MapLeft[bool](func(err error) error {
			return fmt.Errorf("error in setting key in redis %w", err)
		}),
		IOE.Map[error](F.Constant1[bool](rawCachedLens.Get(s))),
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
	return F.Pipe10(
		pubFetchState{},
		identifierLens.Set(id),
		repoLens.Set(repo),
		sourceLens.Set(p),
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

// FetchPublicationIOE fetches a publication by PMID or DOI as a lazy IOEither
// using Redis caching and the literature NCBI and EuropePMC clients.
func FetchPublicationIOE(
	repo repository.Repository,
	id string,
) IOE.IOEither[error, *models.Publication] {
	newSource := NewPublication(&PublicationParams{})
	return F.Pipe10(
		pubFetchState{},
		identifierLens.Set(id),
		repoLens.Set(repo),
		identifierLens.Modify(trimHTTPSDOI),
		identifierLens.Modify(trimHTTPDOI),
		identifierLens.Modify(trimDOIColon),
		IOE.Of[error, pubFetchState],
		IOE.ChainEitherK(validateIdentifierBranch),
		IOE.Bind(sourceLens.Set, F.Constant1[pubFetchState](newSource)),
		IOE.Let[error](redisKeyLens.Set, deriveRedisKey),
		IOE.Chain(branchByRepo),
	)
}

// FetchPublication fetches a publication by PMID or DOI using Redis caching
// and the literature NCBI and EuropePMC clients.
func FetchPublication(
	_ context.Context,
	repo repository.Repository,
	id string,
) (*models.Publication, error) {
	result := F.Pipe2(
		FetchPublicationIOE(repo, id),
		ioeutils.ToEither[error, *models.Publication],
		E.Fold(onPublicationError, onPublicationSuccess),
	)
	return result.F2, result.F1
}
