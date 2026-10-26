package fetch

import (
	"errors"
	"fmt"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	T "github.com/IBM/fp-go/v2/tuple"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictybase/literature"
	"github.com/stretchr/testify/require"
)

func runPublicationIOE(
	ioe IOE.IOEither[error, *models.Publication],
) (*models.Publication, error) {
	result := F.Pipe2(
		ioe,
		ioeutils.ToEither[error, *models.Publication],
		E.Fold(onPublicationError, onPublicationSuccess),
	)
	return result.F2, result.F1
}

const (
	testPMID  = "29997103"
	testDOI   = "10.1101/582072"
	huberLast = "Huber"
)

type fakePubMed struct {
	articles map[string]*literature.Article
	errs     map[string]error
	calls    []string
}

func (f *fakePubMed) GetArticle(pmid string) (*literature.Article, error) {
	f.calls = append(f.calls, pmid)
	if err := f.errs[pmid]; err != nil {
		return nil, err
	}
	return f.articles[pmid], nil
}

func newFakePubMed(article *literature.Article) *fakePubMed {
	return &fakePubMed{
		articles: map[string]*literature.Article{testPMID: article},
		errs:     map[string]error{},
	}
}

func newFakePubMedWithError(pmid string, err error) *fakePubMed {
	return &fakePubMed{
		articles: map[string]*literature.Article{},
		errs:     map[string]error{pmid: err},
	}
}

type fakeEuropePMC struct {
	byPMID    map[string]*literature.EuropePMCArticle
	byDOI     map[string]*literature.EuropePMCArticle
	errs      map[string]error
	pmidCalls []string
	doiCalls  []string
}

func (f *fakeEuropePMC) GetArticle(
	pmid string,
) (*literature.EuropePMCArticle, error) {
	f.pmidCalls = append(f.pmidCalls, pmid)
	if err := f.errs[pmid]; err != nil {
		return nil, err
	}
	return f.byPMID[pmid], nil
}

func (f *fakeEuropePMC) GetArticleByDOI(
	doi string,
) (*literature.EuropePMCArticle, error) {
	f.doiCalls = append(f.doiCalls, doi)
	if err := f.errs[doi]; err != nil {
		return nil, err
	}
	return f.byDOI[doi], nil
}

func newFakeEuropePMC(article *literature.EuropePMCArticle) *fakeEuropePMC {
	return &fakeEuropePMC{
		byPMID: map[string]*literature.EuropePMCArticle{testPMID: article},
		byDOI:  map[string]*literature.EuropePMCArticle{testDOI: article},
		errs:   map[string]error{},
	}
}

func newFakeEuropePMCWithError(id string, err error) *fakeEuropePMC {
	return &fakeEuropePMC{
		byPMID: map[string]*literature.EuropePMCArticle{},
		byDOI:  map[string]*literature.EuropePMCArticle{},
		errs:   map[string]error{id: err},
	}
}

type fakeRepo struct {
	store map[string]string
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{store: map[string]string{}}
}

func (r *fakeRepo) Get(key string) (string, error) {
	return r.store[key], nil
}

func (r *fakeRepo) Set(key, val string) error {
	r.store[key] = val
	return nil
}

func (r *fakeRepo) SetWithTTL(
	key, val string,
	_ time.Duration,
) error {
	return r.Set(key, val)
}

func (r *fakeRepo) Exists(key string) (bool, error) {
	_, ok := r.store[key]
	return ok, nil
}

func (r *fakeRepo) HSet(hash, key, val string) error { return nil }
func (r *fakeRepo) HGet(hash, key string) (string, error) {
	return "", nil
}
func (r *fakeRepo) HExists(hash, key string) (bool, error) {
	return false, nil
}

func samplePubMedArticle() *literature.Article {
	return &literature.Article{
		PMID:        testPMID,
		Title:       "Cln5 in Dictyostelium",
		Abstract:    "Study of Cln5.",
		Journal:     "BBA Molecular Cell Research",
		PublishDate: time.Date(2018, time.July, 23, 0, 0, 0, 0, time.UTC),
		DOI:         "10.1016/j.bbamcr.2018.07.017",
		Volume:      "1865",
		Issue:       "10",
		Pages:       "1432-1444",
		Authors: []literature.Author{
			{FirstName: "Robert J", LastName: huberLast, FullName: "Huber RJ"},
			{FirstName: "Shreya", LastName: "Mathavarajah", FullName: "Mathavarajah S"},
		},
	}
}

func sampleEuropeArticle() *literature.EuropePMCArticle {
	pubDate := time.Date(2018, time.July, 23, 0, 0, 0, 0, time.UTC)
	return &literature.EuropePMCArticle{
		ID:          testPMID,
		PMID:        testPMID,
		DOI:         "10.1016/j.bbamcr.2018.07.017",
		Title:       "Cln5 in Dictyostelium (EuropePMC)",
		Abstract:    "Study of Cln5.",
		Source:      "MED",
		PageInfo:    "1432-1444",
		PubTypes:    []string{"Journal Article"},
		PublishDate: &pubDate,
		Journal: literature.EuropePMCJournal{
			Title:  "BBA Molecular Cell Research",
			ISSN:   "0167-4889",
			Volume: "1865",
			Issue:  "10",
		},
		Authors: []literature.EuropePMCAuthor{
			{
				FirstName: "Robert J",
				LastName:  huberLast,
				Initials:  "RJ",
				FullName:  "Huber RJ",
			},
		},
	}
}

func TestFetchArticleByPMIDUsesNCBIFirst(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ncbi := newFakePubMed(samplePubMedArticle())
	epmc := newFakeEuropePMC(sampleEuropeArticle())
	src := &Publication{pubMed: ncbi, europePMC: epmc}

	pub, err := runPublicationIOE(src.FetchCachedArticle(nil, testPMID))
	req.NoError(err)
	req.Equal(testPMID, pub.ID)
	req.Equal("Cln5 in Dictyostelium", pub.Title)
	req.Equal("Study of Cln5.", pub.Abstract)
	req.Equal("BBA Molecular Cell Research", pub.Journal)
	req.Equal("MED", pub.Source)
	req.Equal("published", *pub.Status)
	req.Equal("10.1016/j.bbamcr.2018.07.017", *pub.Doi)
	req.Equal("1865", *pub.Volume)
	req.Equal("10", *pub.Issue)
	req.Equal("1432-1444", *pub.Pages)
	req.NotNil(pub.PubDate)
	req.Equal(2018, pub.PubDate.Year())
	req.Len(pub.Authors, 2)
	req.Equal(huberLast, pub.Authors[0].LastName)
	req.Equal(int64(0), pub.Authors[0].Rank)
	req.Equal(int64(1), pub.Authors[1].Rank)
	req.Equal([]string{testPMID}, ncbi.calls)
	req.Empty(epmc.pmidCalls, "EuropePMC must not be called when NCBI succeeds")
}

func TestFetchArticleByPMIDFallsBackToEuropePMC(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ncbi := newFakePubMedWithError(testPMID, errors.New("ncbi unavailable"))
	epmc := newFakeEuropePMC(sampleEuropeArticle())
	src := &Publication{pubMed: ncbi, europePMC: epmc}

	pub, err := runPublicationIOE(src.FetchCachedArticle(nil, testPMID))
	req.NoError(err)
	req.Equal(testPMID, pub.ID)
	req.Equal("Cln5 in Dictyostelium (EuropePMC)", pub.Title)
	req.Equal("Journal Article", pub.PubType)
	req.Equal("0167-4889", *pub.Issn)
	req.Equal("RJ", pub.Authors[0].Initials)
	req.Equal([]string{testPMID}, ncbi.calls)
	req.Equal([]string{testPMID}, epmc.pmidCalls)
}

func TestFetchArticleByPMIDBothFail(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ncbi := newFakePubMedWithError(testPMID, errors.New("ncbi down"))
	epmc := newFakeEuropePMCWithError(testPMID, errors.New("europepmc down"))
	src := &Publication{pubMed: ncbi, europePMC: epmc}

	_, err := runPublicationIOE(src.FetchCachedArticle(nil, testPMID))
	req.Error(err)
	req.ErrorContains(err, testPMID)
	req.ErrorContains(err, "europepmc down")
	req.Equal([]string{testPMID}, ncbi.calls)
	req.Equal([]string{testPMID}, epmc.pmidCalls)
}

func TestFetchArticleByDOIPrefixes(t *testing.T) {
	t.Parallel()
	for _, rawInput := range []string{
		testDOI,
		"doi:" + testDOI,
		"https://doi.org/" + testDOI,
		"http://doi.org/" + testDOI,
	} {
		t.Run(rawInput, func(t *testing.T) {
			t.Parallel()
			req := require.New(t)
			ncbi := newFakePubMed(samplePubMedArticle())
			art := sampleEuropeArticle()
			art.PMID = ""
			art.ID = "PPR12345"
			art.DOI = testDOI
			epmc := newFakeEuropePMC(art)
			src := &Publication{pubMed: ncbi, europePMC: epmc}

			pub, err := runPublicationIOE(src.FetchCachedArticle(nil, rawInput))
			req.NoError(err)
			req.Equal("PPR12345", pub.ID)
			req.Equal(testDOI, *pub.Doi)
			req.Equal([]string{testDOI}, epmc.doiCalls)
			req.Empty(ncbi.calls, "DOI lookup must go directly to EuropePMC")
		})
	}
}

func TestFetchArticleZeroPublishDate(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	art := samplePubMedArticle()
	art.PublishDate = time.Time{}
	ncbi := newFakePubMed(art)
	epmc := newFakeEuropePMC(sampleEuropeArticle())
	src := &Publication{pubMed: ncbi, europePMC: epmc}

	pub, err := runPublicationIOE(src.FetchCachedArticle(nil, testPMID))
	req.NoError(err)
	req.Nil(pub.PubDate)
}

func TestFetchArticleInvalidPMID(t *testing.T) {
	t.Parallel()
	for _, pmid := range []string{"", "   ", "abc", "123x", "0"} {
		t.Run(fmt.Sprintf("pmid %q", pmid), func(t *testing.T) {
			t.Parallel()
			req := require.New(t)
			ncbi := newFakePubMed(samplePubMedArticle())
			ncbi.articles[pmid] = ncbi.articles[testPMID]
			epmc := newFakeEuropePMC(sampleEuropeArticle())
			epmc.byPMID[pmid] = epmc.byPMID[testPMID]
			src := &Publication{pubMed: ncbi, europePMC: epmc}

			_, err := runPublicationIOE(src.FetchCachedArticle(nil, pmid))
			req.Error(err)
			req.ErrorContains(err, "pmid must contain decimal digits")

			var litErr *literature.Error
			req.ErrorAs(err, &litErr)
			req.Equal(literature.ErrorTypeInvalidInput, litErr.Type)
			req.Empty(ncbi.calls)
			req.Empty(epmc.pmidCalls)
			req.Empty(epmc.doiCalls)
		})
	}
}

func TestFetchCachedArticleMissAndHit(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ncbi := newFakePubMed(samplePubMedArticle())
	epmc := newFakeEuropePMC(sampleEuropeArticle())
	repo := newFakeRepo()
	src := &Publication{pubMed: ncbi, europePMC: epmc}

	first, err := runPublicationIOE(src.FetchCachedArticle(repo, testPMID))
	req.NoError(err)
	req.Equal(testPMID, first.ID)
	req.Len(ncbi.calls, 1)

	second, err := runPublicationIOE(src.FetchCachedArticle(repo, testPMID))
	req.NoError(err)
	req.Equal(testPMID, second.ID)
	req.Len(ncbi.calls, 1, "cache hit must not call client again")
}

func TestPublicationOptions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		args       *PublicationParams
		wantPubMed int
		wantEurope int
	}{
		{name: "empty params", args: &PublicationParams{}, wantPubMed: 0, wantEurope: 0},
		{
			name:       "api key only",
			args:       &PublicationParams{APIKey: "key"},
			wantPubMed: 1,
			wantEurope: 0,
		},
		{
			name:       "tool only",
			args:       &PublicationParams{Tool: "graphql-server"},
			wantPubMed: 1,
			wantEurope: 1,
		},
		{
			name:       "email only",
			args:       &PublicationParams{Email: "dev@dictybase.org"},
			wantPubMed: 1,
			wantEurope: 1,
		},
		{
			name: "full identity",
			args: &PublicationParams{
				APIKey: "key",
				Tool:   "graphql-server",
				Email:  "dev@dictybase.org",
			},
			wantPubMed: 3,
			wantEurope: 2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Len(t, pubMedOptions(tc.args), tc.wantPubMed)
			require.Len(t, europePMCOptions(tc.args), tc.wantEurope)
		})
	}
}

func TestNewPublicationEmptyParams(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	result := F.Pipe2(
		NewPublication(&PublicationParams{}),
		ioeutils.ToEither[error, *Publication],
		E.Fold(
			func(err error) T.Tuple2[error, *Publication] {
				return T.MakeTuple2[error, *Publication](err, nil)
			},
			func(pub *Publication) T.Tuple2[error, *Publication] {
				return T.MakeTuple2[error](nil, pub)
			},
		),
	)
	req.NoError(result.F1)
	req.NotNil(result.F2)
}
