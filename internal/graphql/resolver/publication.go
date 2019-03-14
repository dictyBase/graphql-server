package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dictyBase/apihelpers/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/graphql-server/internal/registry"
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
	Issue         int64     `json:"issue,omitempty"`
	PublishedDate string    `json:"publication_date"`
	Authors       []*Author `json:"authors"`
}

type Author struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Initials  string `json:"initials"`
}

// Publication is the resolver for getting an individual publication by ID.
func (q *QueryResolver) Publication(ctx context.Context, id string) (*pb.Publication, error) {
	// get publication endpoint from registry
	endpoint := q.GetAPIEndpoint(registry.PUBLICATION)
	url := endpoint + "/" + id
	// get data from endpoint
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in http get request %s", err)
	}
	// close body when done reading from it
	defer res.Body.Close()
	// if response is not 200 (OK) then there's an error
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching publication with ID %s", id)
	}
	// create new json decoder
	decoder := json.NewDecoder(res.Body)
	var pub PubJsonAPI
	// decode the json into established struct
	err = decoder.Decode(&pub)
	if err != nil {
		return nil, fmt.Errorf("error decoding json %s", err)
	}
	attr := pub.Data.Attributes
	// convert pub_date to expected time.Time format
	pd, err := time.Parse("2006-01-02", attr.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("could not parse published date %s", err)
	}
	// convert issue to expected string format
	issue := strconv.Itoa(int(attr.Issue))
	// convert to expected pb model
	p := &pb.Publication{
		Data: &pb.Publication_Data{
			Type: "publication",
			Id:   id,
			Attributes: &pb.PublicationAttributes{
				Doi:      attr.Doi,
				Title:    attr.Title,
				Abstract: attr.Abstract,
				Journal:  attr.Journal,
				PubDate:  aphgrpc.TimestampProto(pd),
				Pages:    attr.Page,
				Issn:     attr.Issn,
				PubType:  attr.PubType,
				Source:   attr.Source,
				Issue:    issue,
				Status:   attr.Status,
				Volume:   "", // field does not exist yet
			},
		},
	}
	var authors []*pb.Author
	// get the list of authors and add it to model
	for _, a := range attr.Authors {
		authors = append(authors, &pb.Author{
			FirstName: a.FirstName,
			LastName:  a.LastName,
			// Rank:      0, // field does not exist yet
			Initials: a.Initials,
		})
	}
	p.Data.Attributes.Authors = authors
	q.Logger.Debugf("successfully found publication with ID %s", pub.Data.ID)
	return p, nil
}
