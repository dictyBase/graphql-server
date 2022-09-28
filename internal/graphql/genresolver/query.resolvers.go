package genresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// GetRefreshToken is the resolver for the getRefreshToken field.
func (r *queryResolver) GetRefreshToken(ctx context.Context, token string) (*auth.Auth, error) {
	panic("not implemented")
}

// Content is the resolver for the content field.
func (r *queryResolver) Content(ctx context.Context, id string) (*content.Content, error) {
	panic("not implemented")
}

// ContentBySlug is the resolver for the contentBySlug field.
func (r *queryResolver) ContentBySlug(ctx context.Context, slug string) (*content.Content, error) {
	panic("not implemented")
}

// Organism is the resolver for the organism field.
func (r *queryResolver) Organism(ctx context.Context, taxonID string) (*models.Organism, error) {
	panic("not implemented")
}

// ListOrganisms is the resolver for the listOrganisms field.
func (r *queryResolver) ListOrganisms(ctx context.Context) ([]*models.Organism, error) {
	panic("not implemented")
}

// Gene is the resolver for the gene field.
func (r *queryResolver) Gene(ctx context.Context, gene string) (*models.Gene, error) {
	panic("not implemented")
}

// AllStrains is the resolver for the allStrains field.
func (r *queryResolver) AllStrains(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: AllStrains - allStrains"))
}

// AllPublications is the resolver for the allPublications field.
func (r *queryResolver) AllPublications(ctx context.Context, gene string, limit *int, sortBy *string) (*models.NumberOfPublicationsWithGene, error) {
	panic(fmt.Errorf("not implemented: AllPublications - allPublications"))
}

// AllOrthologs is the resolver for the allOrthologs field.
func (r *queryResolver) AllOrthologs(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: AllOrthologs - allOrthologs"))
}

// ListRecentGenes is the resolver for the listRecentGenes field.
func (r *queryResolver) ListRecentGenes(ctx context.Context, limit int) ([]*models.Gene, error) {
	panic(fmt.Errorf("not implemented: ListRecentGenes - listRecentGenes"))
}

// GeneralInformation is the resolver for the generalInformation field.
func (r *queryResolver) GeneralInformation(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: GeneralInformation - generalInformation"))
}

// GetAssociatedSequnces is the resolver for the getAssociatedSequnces field.
func (r *queryResolver) GetAssociatedSequnces(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: GetAssociatedSequnces - getAssociatedSequnces"))
}

// GetLinks is the resolver for the getLinks field.
func (r *queryResolver) GetLinks(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: GetLinks - getLinks"))
}

// GetProteinInformation is the resolver for the getProteinInformation field.
func (r *queryResolver) GetProteinInformation(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: GetProteinInformation - getProteinInformation"))
}

// ListGeneProductInfo is the resolver for the listGeneProductInfo field.
func (r *queryResolver) ListGeneProductInfo(ctx context.Context, gene string) (*models.Gene, error) {
	panic(fmt.Errorf("not implemented: ListGeneProductInfo - listGeneProductInfo"))
}

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*order.Order, error) {
	panic("not implemented")
}

// ListOrders is the resolver for the listOrders field.
func (r *queryResolver) ListOrders(ctx context.Context, cursor *int, limit *int, filter *string) (*models.OrderListWithCursor, error) {
	panic("not implemented")
}

// Publication is the resolver for the publication field.
func (r *queryResolver) Publication(ctx context.Context, id string) (*models.Publication, error) {
	panic("not implemented")
}

// ListRecentPublications is the resolver for the listRecentPublications field.
func (r *queryResolver) ListRecentPublications(ctx context.Context, limit int) ([]*models.Publication, error) {
	panic(fmt.Errorf("not implemented: ListRecentPublications - listRecentPublications"))
}

// Plasmid is the resolver for the plasmid field.
func (r *queryResolver) Plasmid(ctx context.Context, id string) (*models.Plasmid, error) {
	panic("not implemented")
}

// Strain is the resolver for the strain field.
func (r *queryResolver) Strain(ctx context.Context, id string) (*models.Strain, error) {
	panic("not implemented")
}

// ListStrains is the resolver for the listStrains field.
func (r *queryResolver) ListStrains(ctx context.Context, cursor *int, limit *int, filter *models.StrainListFilter) (*models.StrainListWithCursor, error) {
	panic("not implemented")
}

// ListPlasmids is the resolver for the listPlasmids field.
func (r *queryResolver) ListPlasmids(ctx context.Context, cursor *int, limit *int, filter *string) (*models.PlasmidListWithCursor, error) {
	panic("not implemented")
}

// ListStrainsWithAnnotation is the resolver for the listStrainsWithAnnotation field.
func (r *queryResolver) ListStrainsWithAnnotation(ctx context.Context, cursor *int, limit *int, typeArg string, annotation string) (*models.StrainListWithCursor, error) {
	panic("not implemented")
}

// ListPlasmidsWithAnnotation is the resolver for the listPlasmidsWithAnnotation field.
func (r *queryResolver) ListPlasmidsWithAnnotation(ctx context.Context, cursor *int, limit *int, typeArg string, annotation string) (*models.PlasmidListWithCursor, error) {
	panic("not implemented")
}

// ListRecentPlasmids is the resolver for the listRecentPlasmids field.
func (r *queryResolver) ListRecentPlasmids(ctx context.Context, limit int) ([]*models.Plasmid, error) {
	panic(fmt.Errorf("not implemented: ListRecentPlasmids - listRecentPlasmids"))
}

// ListRecentStrains is the resolver for the listRecentStrains field.
func (r *queryResolver) ListRecentStrains(ctx context.Context, limit int) ([]*models.Strain, error) {
	panic(fmt.Errorf("not implemented: ListRecentStrains - listRecentStrains"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*user.User, error) {
	panic("not implemented")
}

// UserByEmail is the resolver for the userByEmail field.
func (r *queryResolver) UserByEmail(ctx context.Context, email string) (*user.User, error) {
	panic("not implemented")
}

// ListUsers is the resolver for the listUsers field.
func (r *queryResolver) ListUsers(ctx context.Context, pagenum string, pagesize string, filter string) (*models.UserList, error) {
	panic("not implemented")
}

// Role is the resolver for the role field.
func (r *queryResolver) Role(ctx context.Context, id string) (*user.Role, error) {
	panic("not implemented")
}

// ListRoles is the resolver for the listRoles field.
func (r *queryResolver) ListRoles(ctx context.Context) ([]*user.Role, error) {
	panic("not implemented")
}

// Permission is the resolver for the permission field.
func (r *queryResolver) Permission(ctx context.Context, id string) (*user.Permission, error) {
	panic("not implemented")
}

// ListPermissions is the resolver for the listPermissions field.
func (r *queryResolver) ListPermissions(ctx context.Context) ([]*user.Permission, error) {
	panic("not implemented")
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
