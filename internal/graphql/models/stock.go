package models

import (
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
)

// LazyStrain contains fields that are not directly mapped from model to graphql
// type. The gets resolved to their expected type on demand only
type LazyStrain struct {
	CreatedBy              string
	UpdatedBy              string
	Depositor              string
	Genes, Publications []*string
}

type Strain struct {
	LazyStrain
	ID                  string                     `json:"id"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	Summary             *string                    `json:"summary"`
	EditableSummary     *string                    `json:"editable_summary"`
	Dbxrefs             []*string                  `json:"dbxrefs"`
	SystematicName      string                     `json:"systematic_name"`
	Label               string                     `json:"label"`
	Species             string                     `json:"species"`
	Plasmid             *string                    `json:"plasmid"`
	Parent              *string                    `json:"parent"`
	Names               []*string                  `json:"names"`
	InStock             bool                       `json:"in_stock"`
	GeneticModification *string                    `json:"genetic_modification"`
	MutagenesisMethod   *string                    `json:"mutagenesis_method"`
	Characteristics     []*string                  `json:"characteristics"`
	Genotypes           []*string                  `json:"genotypes"`
}

type Plasmid struct {
	ID               string                     `json:"id"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	CreatedBy        *user.User                 `json:"created_by"`
	UpdatedBy        *user.User                 `json:"updated_by"`
	Summary          *string                    `json:"summary"`
	EditableSummary  *string                    `json:"editable_summary"`
	Depositor        *string                    `json:"depositor"`
	Genes            []*Gene                    `json:"genes"`
	Dbxrefs          []*string                  `json:"dbxrefs"`
	Publications     []*publication.Publication `json:"publications"`
	ImageMap         *string                    `json:"image_map"`
	Sequence         *string                    `json:"sequence"`
	Name             string                     `json:"name"`
	InStock          bool                       `json:"in_stock"`
	Keywords         []*string                  `json:"keywords"`
	GenbankAccession *string                    `json:"genbank_accession"`
}
