package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// CreatedBy is the resolver for the created_by field.
func (r *plasmidResolver) CreatedBy(ctx context.Context, obj *models.Plasmid) (*user.User, error) {
	panic("not implemented")
}

// UpdatedBy is the resolver for the updated_by field.
func (r *plasmidResolver) UpdatedBy(ctx context.Context, obj *models.Plasmid) (*user.User, error) {
	panic("not implemented")
}

// Depositor is the resolver for the depositor field.
func (r *plasmidResolver) Depositor(ctx context.Context, obj *models.Plasmid) (*user.User, error) {
	panic("not implemented")
}

// Genes is the resolver for the genes field.
func (r *plasmidResolver) Genes(ctx context.Context, obj *models.Plasmid) ([]*models.Gene, error) {
	panic("not implemented")
}

// Publications is the resolver for the publications field.
func (r *plasmidResolver) Publications(ctx context.Context, obj *models.Plasmid) ([]*models.Publication, error) {
	panic("not implemented")
}

// InStock is the resolver for the in_stock field.
func (r *plasmidResolver) InStock(ctx context.Context, obj *models.Plasmid) (bool, error) {
	panic("not implemented")
}

// Keywords is the resolver for the keywords field.
func (r *plasmidResolver) Keywords(ctx context.Context, obj *models.Plasmid) ([]string, error) {
	panic("not implemented")
}

// GenbankAccession is the resolver for the genbank_accession field.
func (r *plasmidResolver) GenbankAccession(ctx context.Context, obj *models.Plasmid) (*string, error) {
	panic("not implemented")
}

// CreatedBy is the resolver for the created_by field.
func (r *strainResolver) CreatedBy(ctx context.Context, obj *models.Strain) (*user.User, error) {
	panic("not implemented")
}

// UpdatedBy is the resolver for the updated_by field.
func (r *strainResolver) UpdatedBy(ctx context.Context, obj *models.Strain) (*user.User, error) {
	panic("not implemented")
}

// Depositor is the resolver for the depositor field.
func (r *strainResolver) Depositor(ctx context.Context, obj *models.Strain) (*user.User, error) {
	panic("not implemented")
}

// Genes is the resolver for the genes field.
func (r *strainResolver) Genes(ctx context.Context, obj *models.Strain) ([]*models.Gene, error) {
	panic("not implemented")
}

// Publications is the resolver for the publications field.
func (r *strainResolver) Publications(ctx context.Context, obj *models.Strain) ([]*models.Publication, error) {
	panic("not implemented")
}

// SystematicName is the resolver for the systematic_name field.
func (r *strainResolver) SystematicName(ctx context.Context, obj *models.Strain) (string, error) {
	panic("not implemented")
}

// Parent is the resolver for the parent field.
func (r *strainResolver) Parent(ctx context.Context, obj *models.Strain) (*models.Strain, error) {
	panic("not implemented")
}

// Names is the resolver for the names field.
func (r *strainResolver) Names(ctx context.Context, obj *models.Strain) ([]string, error) {
	panic("not implemented")
}

// InStock is the resolver for the in_stock field.
func (r *strainResolver) InStock(ctx context.Context, obj *models.Strain) (bool, error) {
	panic("not implemented")
}

// Phenotypes is the resolver for the phenotypes field.
func (r *strainResolver) Phenotypes(ctx context.Context, obj *models.Strain) ([]*models.Phenotype, error) {
	panic("not implemented")
}

// GeneticModification is the resolver for the genetic_modification field.
func (r *strainResolver) GeneticModification(ctx context.Context, obj *models.Strain) (*string, error) {
	panic("not implemented")
}

// MutagenesisMethod is the resolver for the mutagenesis_method field.
func (r *strainResolver) MutagenesisMethod(ctx context.Context, obj *models.Strain) (*string, error) {
	panic("not implemented")
}

// Characteristics is the resolver for the characteristics field.
func (r *strainResolver) Characteristics(ctx context.Context, obj *models.Strain) ([]string, error) {
	panic("not implemented")
}

// Genotypes is the resolver for the genotypes field.
func (r *strainResolver) Genotypes(ctx context.Context, obj *models.Strain) ([]string, error) {
	panic("not implemented")
}

// Plasmid returns generated.PlasmidResolver implementation.
func (r *Resolver) Plasmid() generated.PlasmidResolver { return &plasmidResolver{r} }

// Strain returns generated.StrainResolver implementation.
func (r *Resolver) Strain() generated.StrainResolver { return &strainResolver{r} }

type plasmidResolver struct{ *Resolver }
type strainResolver struct{ *Resolver }
