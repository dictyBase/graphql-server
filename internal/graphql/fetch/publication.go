package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Jeffail/gabs/v2"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/fatih/structs"
)

const (
	RedisKey = "PUBLICATION_KEY"
)

var Status string = "published"

type EuroPMC struct {
	HitCount       int64  `json:"hitCount"`
	NextCursorMark string `json:"nextCursorMark"`
	Request        struct {
		CursorMark string `json:"cursorMark"`
		PageSize   int64  `json:"pageSize"`
		Query      string `json:"query"`
		ResultType string `json:"resultType"`
		Sort       string `json:"sort"`
		Synonym    bool   `json:"synonym"`
	} `json:"request"`
	ResultList struct {
		Result []struct {
			AbstractText string `json:"abstractText"`
			Affiliation  string `json:"affiliation"`
			AuthMan      string `json:"authMan"`
			AuthorList   struct {
				Author []struct {
					Affiliation string `json:"affiliation"`
					FirstName   string `json:"firstName"`
					FullName    string `json:"fullName"`
					Initials    string `json:"initials"`
					LastName    string `json:"lastName"`
				} `json:"author"`
			} `json:"authorList"`
			AuthorString              string `json:"authorString"`
			CitedByCount              int64  `json:"citedByCount"`
			DateOfCreation            string `json:"dateOfCreation"`
			DateOfRevision            string `json:"dateOfRevision"`
			Doi                       string `json:"doi"`
			ElectronicPublicationDate string `json:"electronicPublicationDate"`
			EpmcAuthMan               string `json:"epmcAuthMan"`
			FirstPublicationDate      string `json:"firstPublicationDate"`
			FullTextURLList           struct {
				FullTextURL []struct {
					Availability     string `json:"availability"`
					AvailabilityCode string `json:"availabilityCode"`
					DocumentStyle    string `json:"documentStyle"`
					Site             string `json:"site"`
					URL              string `json:"url"`
				} `json:"fullTextUrl"`
			} `json:"fullTextUrlList"`
			HasBook               string `json:"hasBook"`
			HasDBCrossReferences  string `json:"hasDbCrossReferences"`
			HasLabsLinks          string `json:"hasLabsLinks"`
			HasPDF                string `json:"hasPDF"`
			HasReferences         string `json:"hasReferences"`
			HasTMAccessionNumbers string `json:"hasTMAccessionNumbers"`
			HasTextMinedTerms     string `json:"hasTextMinedTerms"`
			ID                    string `json:"id"`
			InEPMC                string `json:"inEPMC"`
			InPMC                 string `json:"inPMC"`
			IsOpenAccess          string `json:"isOpenAccess"`
			JournalInfo           struct {
				DateOfPublication string `json:"dateOfPublication"`
				Journal           struct {
					Essn                string `json:"essn"`
					Isoabbreviation     string `json:"isoabbreviation"`
					Issn                string `json:"issn"`
					MedlineAbbreviation string `json:"medlineAbbreviation"`
					Nlmid               string `json:"nlmid"`
					Title               string `json:"title"`
				} `json:"journal"`
				JournalIssueID       int64  `json:"journalIssueId"`
				MonthOfPublication   int64  `json:"monthOfPublication"`
				PrintPublicationDate string `json:"printPublicationDate"`
				YearOfPublication    int64  `json:"yearOfPublication"`
				Issue                string `json:"issue"`
				Volume               string `json:"volume"`
			} `json:"journalInfo"`
			KeywordList struct {
				Keyword []string `json:"keyword"`
			} `json:"keywordList"`
			Language    string `json:"language"`
			NihAuthMan  string `json:"nihAuthMan"`
			PageInfo    string `json:"pageInfo"`
			Pmid        string `json:"pmid"`
			PubModel    string `json:"pubModel"`
			PubTypeList struct {
				PubType []string `json:"pubType"`
			} `json:"pubTypeList"`
			PubYear string `json:"pubYear"`
			Source  string `json:"source"`
			Title   string `json:"title"`
		} `json:"result"`
	} `json:"resultList"`
	Version string `json:"version"`
}

type PubJSONAPI struct {
	Data  *PubData `json:"data"`
	Links *Links   `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}

type PubData struct {
	Type       string       `json:"type"`
	ID         string       `json:"id"`
	Attributes *Publication `json:"attributes"`
}

type Publication struct {
	Abstract      string    `json:"abstract"`
	Doi           string    `json:"doi,omitempty"`
	FullTextURL   string    `json:"full_text_url,omitempty"`
	PubmedURL     string    `json:"pubmed_url"`
	Journal       string    `json:"journal"`
	Issn          string    `json:"issn,omitempty"`
	Page          string    `json:"page,omitempty"`
	Pubmed        string    `json:"pubmed"`
	Title         string    `json:"title"`
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	PubType       string    `json:"pub_type"`
	Issue         string    `json:"issue"`
	Volume        string    `json:"volume"`
	PublishedDate string    `json:"publication_date"`
	Authors       []*Author `json:"authors"`
}

type Author struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Initials  string `json:"initials"`
}

func EuroPMC2Pub(pmc *EuroPMC, id string) (*models.Publication, error) {
	if len(pmc.ResultList.Result) == 0 {
		return nil, fmt.Errorf("no publication found with id %s", id)
	}
	result := pmc.ResultList.Result[0]
	pub := &models.Publication{
		ID:       result.Pmid,
		Title:    result.Title,
		Abstract: result.AbstractText,
		Journal:  result.JournalInfo.Journal.Title,
		Source:   result.Source,
		Doi:      &result.Doi,
		Issue:    &result.JournalInfo.Issue,
		Status:   &Status,
		Volume:   &result.JournalInfo.Volume,
		Pages:    &result.PageInfo,
		Issn:     &result.JournalInfo.Journal.Issn,
	}
	pdate, err := time.Parse(time.DateOnly, result.FirstPublicationDate)
	if err != nil {
		return pub, fmt.Errorf("error in parsing date %s", err)
	}
	pub.PubDate = &pdate
	rstruct := structs.New(result)
	if !rstruct.Field("PubTypeList").IsZero() {
		pub.PubType = result.PubTypeList.PubType[0]
	}
	for i, a := range result.AuthorList.Author {
		pub.Authors = append(pub.Authors, &pb.Author{
			FirstName: a.FirstName,
			LastName:  a.LastName,
			Rank:      int64(i),
			Initials:  a.Initials,
		})
	}
	return pub, nil
}

func FetchPublicationFromEuroPMC(
	ctx context.Context,
	endpoint, id string,
) (*models.Publication, error) {
	pmodel := new(models.Publication)
	res, err := http.Get(fmt.Sprintf(
		"%s?format=json&resultType=core&query=ext_id:%s",
		endpoint, id,
	))
	if err != nil {
		return pmodel, fmt.Errorf("error in fetching data %s", err)
	}
	defer res.Body.Close()
	epmc := &EuroPMC{}
	err = json.NewDecoder(res.Body).Decode(epmc)
	if err != nil {
		return pmodel, fmt.Errorf("error in decoding data %s", err)
	}
	return EuroPMC2Pub(epmc, id)
}

func FetchPublication(
	ctx context.Context,
	repo repository.Repository,
	endpoint, id string,
) (*models.Publication, error) {
	ok, pubeuro, err := FetchPublicationFromCache(repo, id)
	if err != nil {
		return nil, fmt.Errorf(
			"error in fetching publication from cache %s",
			err,
		)
	}
	if ok {
		return pubeuro, nil
	}
	pubeuro, err = FetchPublicationFromEuroPMC(ctx, endpoint, id)
	if err != nil {
		return nil, fmt.Errorf("error in fetching publication %s", err)
	}
	if err := StorePublicationInCache(id, repo, pubeuro); err != nil {
		return nil, fmt.Errorf("error in storing publication in cache %s", err)
	}
	return pubeuro, nil
}

func FetchDOI(ctx context.Context, doi string) (*models.Publication, error) {
	res, err := getDOIResp(doi)
	if err != nil {
		return nil, fmt.Errorf("error in getting response for doi %s", err)
	}
	defer res.Body.Close()
	jstruct, err := gabs.ParseJSONDecoder(json.NewDecoder(res.Body))
	if err != nil {
		return nil, fmt.Errorf("error in decoding json %s", err)
	}
	authors := getDOIAuthors(jstruct.Search("author").Children())
	dateStr, ok := jstruct.Search("created", "date-time").Data().(string)
	if !ok {
		return nil, fmt.Errorf(
			"error in converting date time to string %s",
			err,
		)
	}
	pd, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse published date %s", err)
	}
	mpub := &models.Publication{
		ID:       "",
		Doi:      &doi,
		Title:    verifyStringProperty(jstruct, "title"),
		Abstract: verifyStringProperty(jstruct, "abstract"),
		Journal: verifyStringProperty(
			jstruct,
			"container-title-short",
		),
		Source:  verifyStringProperty(jstruct, "source"),
		PubType: verifyStringProperty(jstruct, "type"),
		Issue:   verifyStringPointerProperty(jstruct, "issue"),
		Status:  verifyStringPointerProperty(jstruct, "subtype"),
		Volume:  verifyStringPointerProperty(jstruct, "volume"),
		Pages:   verifyStringPointerProperty(jstruct, "page"),
		Issn:    verifyStringPointerProperty(jstruct, "ISSN"),
		Authors: authors,
		PubDate: &pd,
	}
	return mpub, nil
}

// getDOIResp makes HTTP request with necessary
// headers for DOI and returns the response
func getDOIResp(doi string) (*http.Response, error) {
	url, err := url.Parse(doi)
	if err != nil {
		return nil, fmt.Errorf("error in parsing doi url %s", err)
	}
	res, err := http.DefaultClient.Do(
		&http.Request{
			Method: "GET",
			URL:    url,
			Header: map[string][]string{
				"Accept": {"application/vnd.citationstyles.csl+json"},
			},
		})
	if err != nil {
		return res, fmt.Errorf("error in doi http request %s", err)
	}
	return res, nil
}

// getDOIAuthors converts DOI authors data into expected Author format
func getDOIAuthors(authors []*gabs.Container) []*pb.Author {
	a := []*pb.Author{}
	for _, v := range authors {
		n := &pb.Author{}
		for key, val := range v.ChildrenMap() {
			if key == "given" {
				n.FirstName = val.Data().(string)
			}
			if key == "family" {
				n.LastName = val.Data().(string)
			}
		}
		a = append(a, n)
	}
	return a
}

func verifyStringPointerProperty(jstruct *gabs.Container, val string) *string {
	emptyStr := ""
	if !jstruct.Exists(val) {
		return &emptyStr
	}
	str, ok := jstruct.Search(val).Data().(string)
	if !ok {
		return &emptyStr
	}
	return &str
}

// verifyStringProperty checks if a property exists in the JSON and returns the
// value if true
func verifyStringProperty(jstruct *gabs.Container, val string) string {
	if !jstruct.Exists(val) {
		return ""
	}
	str, ok := jstruct.Search(val).Data().(string)
	if !ok {
		return ""
	}
	return str
}

func FetchPublicationFromCache(
	repo repository.Repository,
	id string,
) (bool, *models.Publication, error) {
	pmodel := &models.Publication{}
	rkey := makeRedisKey(id)
	ok, err := repo.Exists(rkey)
	if err != nil {
		return ok, pmodel, fmt.Errorf(
			"error in checking for key %s",
			err,
		)
	}
	if !ok {
		return ok, pmodel, nil
	}
	pubr, err := repo.Get(rkey)
	if err != nil {
		return ok, pmodel, fmt.Errorf(
			"error in getting existing key %s",
			err,
		)
	}
	if err := json.Unmarshal([]byte(pubr), pmodel); err != nil {
		return ok, pmodel, fmt.Errorf("error in decoding json %s", err)
	}
	return ok, pmodel, nil
}

func StorePublicationInCache(
	id string,
	repo repository.Repository,
	pub *models.Publication,
) error {
	cnt, err := json.Marshal(pub)
	if err != nil {
		return fmt.Errorf("error in converting json to byte %s", err)
	}
	if err := repo.Set(makeRedisKey(id), string(cnt)); err != nil {
		return fmt.Errorf("error in setting key in redis %s", err)
	}
	return nil
}

func makeRedisKey(id string) string {
	return fmt.Sprintf("%s/%s", RedisKey, id)
}

func GetResp(ctx context.Context, url string) (*http.Response, error) {
	res, err := http.Get(url) //nolint:gosec
	if err != nil {
		return res, fmt.Errorf("error in http get request with %s", err)
	}
	if res.StatusCode != 200 {
		return res, fmt.Errorf(
			"error fetching data with status code %d",
			res.StatusCode,
		)
	}
	return res, nil
}
