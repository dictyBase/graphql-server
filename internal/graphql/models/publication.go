package models

import (
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/publication"
)

type Publication struct {
	ID       string                `json:"id"`
	Doi      *string               `json:"doi"`
	Title    string                `json:"title"`
	Abstract string                `json:"abstract"`
	Journal  string                `json:"journal"`
	PubDate  *time.Time            `json:"pub_date"`
	Volume   *string               `json:"volume"`
	Pages    *string               `json:"pages"`
	Issn     *string               `json:"issn"`
	PubType  string                `json:"pub_type"`
	Source   string                `json:"source"`
	Issue    *string               `json:"issue"`
	Status   *string               `json:"status"`
	Authors  []*publication.Author `json:"authors"`
}

func (Publication) IsBasePublication() {}
