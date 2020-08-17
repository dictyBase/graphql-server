// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
)

type Stock interface {
	IsStock()
}

type Citation struct {
	Authors  string `json:"authors"`
	Title    string `json:"title"`
	Journal  string `json:"journal"`
	PubmedID string `json:"pubmed_id"`
}

type CreateContentInput struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	Content   string `json:"content"`
	Namespace string `json:"namespace"`
}

type CreateOrderInput struct {
	Courier          string     `json:"courier"`
	CourierAccount   string     `json:"courier_account"`
	Comments         *string    `json:"comments"`
	Payment          string     `json:"payment"`
	PurchaseOrderNum *string    `json:"purchase_order_num"`
	Status           StatusEnum `json:"status"`
	Consumer         string     `json:"consumer"`
	Payer            string     `json:"payer"`
	Purchaser        string     `json:"purchaser"`
	Items            []*string  `json:"items"`
}

type CreatePermissionInput struct {
	Permission  string `json:"permission"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
}

type CreatePlasmidInput struct {
	CreatedBy        string    `json:"created_by"`
	UpdatedBy        string    `json:"updated_by"`
	Summary          *string   `json:"summary"`
	EditableSummary  *string   `json:"editable_summary"`
	Depositor        *string   `json:"depositor"`
	Genes            []*string `json:"genes"`
	Dbxrefs          []*string `json:"dbxrefs"`
	Publications     []*string `json:"publications"`
	Name             string    `json:"name"`
	ImageMap         *string   `json:"image_map"`
	Sequence         *string   `json:"sequence"`
	InStock          bool      `json:"in_stock"`
	Keywords         []*string `json:"keywords"`
	GenbankAccession *string   `json:"genbank_accession"`
}

type CreateRoleInput struct {
	Role        string `json:"role"`
	Description string `json:"description"`
}

type CreateStrainInput struct {
	CreatedBy           string    `json:"created_by"`
	UpdatedBy           string    `json:"updated_by"`
	Summary             *string   `json:"summary"`
	EditableSummary     *string   `json:"editable_summary"`
	Depositor           *string   `json:"depositor"`
	Genes               []*string `json:"genes"`
	Dbxrefs             []*string `json:"dbxrefs"`
	Publications        []*string `json:"publications"`
	SystematicName      string    `json:"systematic_name"`
	Label               string    `json:"label"`
	Species             string    `json:"species"`
	Plasmid             *string   `json:"plasmid"`
	Parent              *string   `json:"parent"`
	Names               []*string `json:"names"`
	InStock             bool      `json:"in_stock"`
	Phenotypes          []*string `json:"phenotypes"`
	GeneticModification *string   `json:"genetic_modification"`
	MutagenesisMethod   *string   `json:"mutagenesis_method"`
	Characteristics     []*string `json:"characteristics"`
	Genotypes           []*string `json:"genotypes"`
}

type CreateUserInput struct {
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	Organization  *string `json:"organization"`
	GroupName     *string `json:"group_name"`
	FirstAddress  *string `json:"first_address"`
	SecondAddress *string `json:"second_address"`
	City          *string `json:"city"`
	State         *string `json:"state"`
	Zipcode       *string `json:"zipcode"`
	Country       *string `json:"country"`
	Phone         *string `json:"phone"`
	IsActive      bool    `json:"is_active"`
}

type DeleteContent struct {
	Success bool `json:"success"`
}

type DeletePermission struct {
	Success bool `json:"success"`
}

type DeleteRole struct {
	Success bool `json:"success"`
}

type DeleteStock struct {
	Success bool `json:"success"`
}

type DeleteUser struct {
	Success bool `json:"success"`
}

type Download struct {
	Title string          `json:"title"`
	Items []*DownloadItem `json:"items"`
}

type DownloadItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Extension struct {
	ID       string `json:"id"`
	Db       string `json:"db"`
	Relation string `json:"relation"`
	Name     string `json:"name"`
}

type GOAnnotation struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Date         string       `json:"date"`
	EvidenceCode string       `json:"evidence_code"`
	GoTerm       string       `json:"go_term"`
	Qualifier    string       `json:"qualifier"`
	Publication  string       `json:"publication"`
	With         []*With      `json:"with"`
	Extensions   []*Extension `json:"extensions"`
	AssignedBy   string       `json:"assigned_by"`
}

type Gene struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Goas []*GOAnnotation `json:"goas"`
}

type Identity struct {
	ID         string    `json:"id"`
	Identifier string    `json:"identifier"`
	Provider   string    `json:"provider"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListOrderInput struct {
	Cursor *int    `json:"cursor"`
	Limit  *int    `json:"limit"`
	Filter *string `json:"filter"`
}

type ListStockInput struct {
	Cursor *int    `json:"cursor"`
	Limit  *int    `json:"limit"`
	Filter *string `json:"filter"`
}

type ListStrainsWithPhenotypeInput struct {
	Cursor    *int   `json:"cursor"`
	Limit     *int   `json:"limit"`
	Phenotype string `json:"phenotype"`
}

type LoginInput struct {
	ClientID    string `json:"client_id"`
	State       string `json:"state"`
	Code        string `json:"code"`
	Scopes      string `json:"scopes"`
	Provider    string `json:"provider"`
	RedirectURL string `json:"redirect_url"`
}

type Logout struct {
	Success bool `json:"success"`
}

type OrderListWithCursor struct {
	Orders         []*order.Order `json:"orders"`
	NextCursor     int            `json:"nextCursor"`
	PreviousCursor int            `json:"previousCursor"`
	Limit          *int           `json:"limit"`
	TotalCount     int            `json:"totalCount"`
}

type Organism struct {
	TaxonID        string      `json:"taxon_id"`
	ScientificName string      `json:"scientific_name"`
	Citations      []*Citation `json:"citations"`
	Downloads      []*Download `json:"downloads"`
}

type Phenotype struct {
	Phenotype   string                   `json:"phenotype"`
	Note        *string                  `json:"note"`
	Assay       *string                  `json:"assay"`
	Environment *string                  `json:"environment"`
	Publication *publication.Publication `json:"publication"`
}

type PlasmidListWithCursor struct {
	Plasmids       []*Plasmid `json:"plasmids"`
	NextCursor     int        `json:"nextCursor"`
	PreviousCursor int        `json:"previousCursor"`
	Limit          *int       `json:"limit"`
	TotalCount     int        `json:"totalCount"`
}

type StrainListWithCursor struct {
	Strains        []*Strain `json:"strains"`
	NextCursor     int       `json:"nextCursor"`
	PreviousCursor int       `json:"previousCursor"`
	Limit          *int      `json:"limit"`
	TotalCount     int       `json:"totalCount"`
}

type UpdateContentInput struct {
	ID        string `json:"id"`
	UpdatedBy string `json:"updated_by"`
	Content   string `json:"content"`
}

type UpdateOrderInput struct {
	Courier          *string     `json:"courier"`
	CourierAccount   *string     `json:"courier_account"`
	Comments         *string     `json:"comments"`
	Payment          *string     `json:"payment"`
	PurchaseOrderNum *string     `json:"purchase_order_num"`
	Status           *StatusEnum `json:"status"`
	Items            []*string   `json:"items"`
}

type UpdatePermissionInput struct {
	Permission  string `json:"permission"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
}

type UpdatePlasmidInput struct {
	UpdatedBy        string    `json:"updated_by"`
	Summary          *string   `json:"summary"`
	EditableSummary  *string   `json:"editable_summary"`
	Depositor        *string   `json:"depositor"`
	Genes            []*string `json:"genes"`
	Dbxrefs          []*string `json:"dbxrefs"`
	Publications     []*string `json:"publications"`
	Name             *string   `json:"name"`
	ImageMap         *string   `json:"image_map"`
	Sequence         *string   `json:"sequence"`
	InStock          *bool     `json:"in_stock"`
	Keywords         []*string `json:"keywords"`
	GenbankAccession *string   `json:"genbank_accession"`
}

type UpdateRoleInput struct {
	Role        string `json:"role"`
	Description string `json:"description"`
}

type UpdateStrainInput struct {
	UpdatedBy           string    `json:"updated_by"`
	Summary             *string   `json:"summary"`
	EditableSummary     *string   `json:"editable_summary"`
	Depositor           *string   `json:"depositor"`
	Genes               []*string `json:"genes"`
	Dbxrefs             []*string `json:"dbxrefs"`
	Publications        []*string `json:"publications"`
	SystematicName      *string   `json:"systematic_name"`
	Label               *string   `json:"label"`
	Species             *string   `json:"species"`
	Plasmid             *string   `json:"plasmid"`
	Parent              *string   `json:"parent"`
	Names               []*string `json:"names"`
	InStock             *bool     `json:"in_stock"`
	Phenotypes          []*string `json:"phenotypes"`
	GeneticModification *string   `json:"genetic_modification"`
	MutagenesisMethod   *string   `json:"mutagenesis_method"`
	Characteristics     []*string `json:"characteristics"`
	Genotypes           []*string `json:"genotypes"`
}

type UpdateUserInput struct {
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	Organization  *string `json:"organization"`
	GroupName     *string `json:"group_name"`
	FirstAddress  *string `json:"first_address"`
	SecondAddress *string `json:"second_address"`
	City          *string `json:"city"`
	State         *string `json:"state"`
	Zipcode       *string `json:"zipcode"`
	Country       *string `json:"country"`
	Phone         *string `json:"phone"`
	IsActive      *bool   `json:"is_active"`
}

type UserList struct {
	Users      []*user.User `json:"users"`
	PageNum    *string      `json:"pageNum"`
	PageSize   *string      `json:"pageSize"`
	TotalCount int          `json:"totalCount"`
}

type With struct {
	ID   string `json:"id"`
	Db   string `json:"db"`
	Name string `json:"name"`
}

type StatusEnum string

const (
	StatusEnumInPreparation StatusEnum = "IN_PREPARATION"
	StatusEnumGrowing       StatusEnum = "GROWING"
	StatusEnumCancelled     StatusEnum = "CANCELLED"
	StatusEnumShipped       StatusEnum = "SHIPPED"
)

var AllStatusEnum = []StatusEnum{
	StatusEnumInPreparation,
	StatusEnumGrowing,
	StatusEnumCancelled,
	StatusEnumShipped,
}

func (e StatusEnum) IsValid() bool {
	switch e {
	case StatusEnumInPreparation, StatusEnumGrowing, StatusEnumCancelled, StatusEnumShipped:
		return true
	}
	return false
}

func (e StatusEnum) String() string {
	return string(e)
}

func (e *StatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StatusEnum", str)
	}
	return nil
}

func (e StatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
