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
)

type PubJsonAPI struct {
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

func FetchPublication(
	ctx context.Context,
	endpoint, id string,
) (*models.Publication, error) {
	res, err := GetResp(ctx, fmt.Sprintf("%s/%s", endpoint, id))
	if err != nil {
		return nil, fmt.Errorf("error in getting response %s", err)
	}
	pub := &PubJsonAPI{}
	if err := json.NewDecoder(res.Body).Decode(pub); err != nil {
		return nil, fmt.Errorf("error decoding json %s", err)
	}
	// convert pub_date to expected time.Time format
	pd, err := time.Parse("2006-01-02", pub.Data.Attributes.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("could not parse published date %s", err)
	}
	mpub := &models.Publication{
		ID:       id,
		Title:    pub.Data.Attributes.Title,
		Abstract: pub.Data.Attributes.Abstract,
		Journal:  pub.Data.Attributes.Journal,
		Source:   pub.Data.Attributes.Source,
		PubType:  pub.Data.Attributes.PubType,
		Doi:      &pub.Data.Attributes.Doi,
		Issue:    &pub.Data.Attributes.Issue,
		Status:   &pub.Data.Attributes.Status,
		Volume:   &pub.Data.Attributes.Volume,
		Pages:    &pub.Data.Attributes.Page,
		Issn:     &pub.Data.Attributes.Issn,
		PubDate:  &pd,
	}

	for i, a := range pub.Data.Attributes.Authors {
		mpub.Authors = append(mpub.Authors, &pb.Author{
			FirstName: a.FirstName,
			LastName:  a.LastName,
			Rank:      int64(i),
			Initials:  a.Initials,
		})
	}
	return mpub, nil
}

func FetchDOI(ctx context.Context, doi string) (*models.Publication, error) {
	res, err := getDOIResp(ctx, doi)
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
func getDOIResp(ctx context.Context, doi string) (*http.Response, error) {
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

// verifyArrayProperty checks if a property exists in the JSON and then returns
// its first child as a string
func verifyArrayProperty(jstruct *gabs.Container, val string) string {
	if jstruct.Exists(val) {
		return jstruct.Search(val).Children()[0].Data().(string)
	}
	return ""
}
