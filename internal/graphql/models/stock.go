package models

import (
	"time"
)

type StockCommon struct{}

func (stck StockCommon) IsStock() {}

// LazyStock contains fields that are not directly mapped from model to graphql
// type. These fields gets resolved on demand only
type LazyStock struct {
	CreatedBy           string
	UpdatedBy           string
	Depositor           string
	Genes, Publications []string
}

type Strain struct {
	LazyStock
	StockCommon
	ID                  string    `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Summary             *string   `json:"summary"`
	EditableSummary     *string   `json:"editable_summary"`
	Dbxrefs             []string `json:"dbxrefs"`
	SystematicName      string    `json:"systematic_name"`
	Label               string    `json:"label"`
	Species             string    `json:"species"`
	Plasmid             *string   `json:"plasmid"`
	Parent              *string   `json:"parent"`
	Names               []string `json:"names"`
	InStock             bool      `json:"in_stock"`
	GeneticModification *string   `json:"genetic_modification"`
	MutagenesisMethod   *string   `json:"mutagenesis_method"`
	Characteristics     []string `json:"characteristics"`
	Genotypes           []string `json:"genotypes"`
}

type Plasmid struct {
	StockCommon
	LazyStock
	ID               string    `json:"id"`
	Summary          *string   `json:"summary"`
	EditableSummary  *string   `json:"editable_summary"`
	ImageMap         *string   `json:"image_map"`
	Sequence         *string   `json:"sequence"`
	Name             string    `json:"name"`
	GenbankAccession *string   `json:"genbank_accession"`
	InStock          bool      `json:"in_stock"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Keywords         []*string `json:"keywords"`
	Dbxrefs          []*string `json:"dbxrefs"`
}
