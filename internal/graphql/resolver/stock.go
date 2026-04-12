package resolver

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func (mrs *MutationResolver) CreateStrain(
	ctx context.Context,
	input *models.CreateStrainInput,
) (*models.Strain, error) {
	attr := &pb.NewStrainAttributes{}
	norm := normalizeCreateStrainAttr(input)
	err := mapstructure.Decode(norm, attr)
	if err != nil {
		mrs.Logger.Error(err)
		return nil, err
	}
	n, err := mrs.GetStockClient(registry.STOCK).
		CreateStrain(ctx, &pb.NewStrain{
			Data: &pb.NewStrain_Data{
				Type:       "strain",
				Attributes: attr,
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	strainID := n.Data.Id
	// Note: InStock, Phenotypes, GeneticModification, MutagenesisMethod,
	// Characteristics, SystematicName and Genotypes will need to be implemented later.
	mrs.Logger.Debugf("successfully created new strain with ID %s", strainID)
	return stock.ConvertToStrainModel(strainID, n.Data.Attributes), nil
}

func normalizeCreateStrainAttr(
	attr *models.CreateStrainInput,
) map[string]interface{} {
	fields := structs.Fields(attr)
	newAttr := make(map[string]interface{})
	for _, k := range fields {
		if !k.IsZero() {
			newAttr[k.Name()] = k.Value()
		} else {
			switch k.Name() {
			case "Genes":
				newAttr[k.Name()] = nil
			case "Dbxrefs":
				newAttr[k.Name()] = nil
			case "Publications":
				newAttr[k.Name()] = nil
			case "Names":
				newAttr[k.Name()] = nil
			case "Phenotypes":
				newAttr[k.Name()] = nil
			case "Characteristics":
				newAttr[k.Name()] = nil
			case "Genotypes":
				newAttr[k.Name()] = nil
			default:
				newAttr[k.Name()] = ""
			}
		}
	}
	return newAttr
}

func (mrs *MutationResolver) CreatePlasmid(
	ctx context.Context,
	input *models.CreatePlasmidInput,
) (*models.Plasmid, error) {
	attr := &pb.NewPlasmidAttributes{}
	norm := normalizeCreatePlasmidAttr(input)
	err := mapstructure.Decode(norm, attr)
	if err != nil {
		mrs.Logger.Error(err)
		return nil, err
	}
	n, err := mrs.GetStockClient(registry.STOCK).
		CreatePlasmid(ctx, &pb.NewPlasmid{
			Data: &pb.NewPlasmid_Data{
				Type:       "plasmid",
				Attributes: attr,
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	plasmidID := n.Data.Id
	// Note: InStock, Keywords and GenbankAccession will need to be implemented later.
	mrs.Logger.Debugf("successfully created new plasmid with ID %s", plasmidID)
	return stock.ConvertToPlasmidModel(plasmidID, n.Data.Attributes), nil
}

func normalizeCreatePlasmidAttr(
	attr *models.CreatePlasmidInput,
) map[string]interface{} {
	fields := structs.Fields(attr)
	newAttr := make(map[string]interface{})
	for _, k := range fields {
		if !k.IsZero() {
			newAttr[k.Name()] = k.Value()
		} else {
			switch k.Name() {
			case "Genes":
				newAttr[k.Name()] = nil
			case "Dbxrefs":
				newAttr[k.Name()] = nil
			case "Publications":
				newAttr[k.Name()] = nil
			case "Keywords":
				newAttr[k.Name()] = nil
			default:
				newAttr[k.Name()] = ""
			}
		}
	}
	return newAttr
}

//nolint:dupl
func (mrs *MutationResolver) UpdateStrain(
	ctx context.Context,
	id string,
	input *models.UpdateStrainInput,
) (*models.Strain, error) {
	_, err := mrs.GetStockClient(registry.STOCK).
		GetStrain(ctx, &pb.StockId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	attr := &pb.StrainUpdateAttributes{}
	norm := normalizeUpdateStrainAttr(input)
	err = mapstructure.Decode(norm, attr)
	if err != nil {
		mrs.Logger.Error(err)
		return nil, err
	}
	n, err := mrs.GetStockClient(registry.STOCK).
		UpdateStrain(ctx, &pb.StrainUpdate{
			Data: &pb.StrainUpdate_Data{
				Type:       "strain",
				Id:         id,
				Attributes: attr,
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	strainID := n.Data.Id
	mrs.Logger.Debugf("successfully updated strain with ID %s", strainID)
	return stock.ConvertToStrainModel(strainID, n.Data.Attributes), nil
}

func normalizeUpdateStrainAttr(
	attr *models.UpdateStrainInput,
) map[string]interface{} {
	fields := structs.Fields(attr)
	newAttr := make(map[string]interface{})
	for _, k := range fields {
		if !k.IsZero() {
			newAttr[k.Name()] = k.Value()
		}
	}
	return newAttr
}

//nolint:dupl
func (mrs *MutationResolver) UpdatePlasmid(
	ctx context.Context,
	id string,
	input *models.UpdatePlasmidInput,
) (*models.Plasmid, error) {
	_, err := mrs.GetStockClient(registry.STOCK).
		GetPlasmid(ctx, &pb.StockId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	attr := &pb.PlasmidUpdateAttributes{}
	norm := normalizeUpdatePlasmidAttr(input)
	err = mapstructure.Decode(norm, attr)
	if err != nil {
		mrs.Logger.Error(err)
		return nil, err
	}
	n, err := mrs.GetStockClient(registry.STOCK).
		UpdatePlasmid(ctx, &pb.PlasmidUpdate{
			Data: &pb.PlasmidUpdate_Data{
				Type:       "plasmid",
				Id:         id,
				Attributes: attr,
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	plasmidID := n.Data.Id
	mrs.Logger.Debugf("successfully updated plasmid with ID %s", plasmidID)
	return stock.ConvertToPlasmidModel(plasmidID, n.Data.Attributes), nil
}

func normalizeUpdatePlasmidAttr(
	attr *models.UpdatePlasmidInput,
) map[string]interface{} {
	fields := structs.Fields(attr)
	newAttr := make(map[string]interface{})
	for _, k := range fields {
		if !k.IsZero() {
			newAttr[k.Name()] = k.Value()
		}
	}
	return newAttr
}

func (mrs *MutationResolver) DeleteStock(
	ctx context.Context,
	id string,
) (*models.DeleteStock, error) {
	if _, err := mrs.GetStockClient(registry.STOCK).RemoveStock(ctx, &pb.StockId{Id: id}); err != nil {
		return &models.DeleteStock{
			Success: false,
		}, err
	}
	mrs.Logger.Debugf("successfully deleted stock with ID %s", id)
	return &models.DeleteStock{
		Success: true,
	}, nil
}

func (qrs *QueryResolver) Plasmid(
	ctx context.Context,
	id string,
) (*models.Plasmid, error) {
	n, err := qrs.GetStockClient(registry.STOCK).
		GetPlasmid(ctx, &pb.StockId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	plasmidID := n.Data.Id
	qrs.Logger.Debugf("successfully found plasmid with ID %s", plasmidID)
	return stock.ConvertToPlasmidModel(plasmidID, n.Data.Attributes), nil
}

func (qrs *QueryResolver) Strain(
	ctx context.Context,
	id string,
) (*models.Strain, error) {
	n, err := qrs.GetStockClient(registry.STOCK).
		GetStrain(ctx, &pb.StockId{Id: id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	strainID := n.Data.Id
	qrs.Logger.Debugf("successfully found strain with ID %s", strainID)
	return stock.ConvertToStrainModel(strainID, n.Data.Attributes), nil
}

func (qrs *QueryResolver) ListStrains(ctx context.Context, cursor *int,
	limit *int, filter *models.StrainListFilter,
) (*models.StrainListWithCursor, error) {
	cus := resolverutils.GetCursor(cursor)
	lmt := resolverutils.GetLimit(limit)
	// no filter , get a limited set of strain
	if filter == nil {
		return qrs.listStrainsWithoutFilter(ctx, cus, lmt)
	}
	stypeQuery, err := resolverutils.StrainFilterToQuery(filter)
	if err != nil {
		return qrs.reportStrainListError(ctx, err)
	}
	strainList, err := qrs.GetStockClient(registry.STOCK).
		ListStrains(ctx, &pb.StockParameters{
			Cursor: cus,
			Limit:  lmt,
			Filter: stypeQuery,
		})
	if err != nil {
		return qrs.reportStrainListError(ctx, err)
	}

	return qrs.toStrainModelList(strainList, lmt, cus), nil
}

func (qrs *QueryResolver) ListPlasmids(
	ctx context.Context,
	cursor *int,
	limit *int,
	filter *models.PlasmidListFilter,
) (*models.PlasmidListWithCursor, error) {
	cus := resolverutils.GetCursor(cursor)
	lmt := resolverutils.GetLimit(limit)

	filterQuery, err := resolverutils.PlasmidFilterToQuery(filter)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return &models.PlasmidListWithCursor{}, err
	}

	result := foldPlasmidListResult(
		F.Pipe1(
			callListPlasmids(
				qrs.GetStockClient(registry.STOCK),
				ctx,
				buildPlasmidStockParameters(cus, lmt, filterQuery),
			),
			IOE.Map[error](func(list *pb.PlasmidCollection) *models.PlasmidListWithCursor {
				return buildPlasmidListResult(list, cus)
			}),
		),
	)

	if result.F1 != nil {
		errorutils.AddGQLError(ctx, result.F1)
		qrs.Logger.Error(result.F1)
		return result.F2, result.F1
	}

	qrs.Logger.Debugf("successfully retrieved list of %v plasmids", result.F2.TotalCount)
	return result.F2, nil
}

//nolint:dupl
func (qrs *QueryResolver) ListStrainsWithAnnotation(
	ctx context.Context,
	cursor *int,
	limit *int,
	typeArg string,
	annotation string,
) (*models.StrainListWithCursor, error) {
	strains := []*models.Strain{}
	cur := resolverutils.GetCursor(cursor)
	lmt := resolverutils.GetLimit(limit)
	onto := resolverutils.GetOntology(typeArg)
	ann, err := qrs.GetAnnotationClient(registry.ANNOTATION).
		ListAnnotations(ctx, &anno.ListParameters{
			Cursor: cur,
			Limit:  lmt,
			Filter: fmt.Sprintf("ontology==%s;tag==%s", onto, annotation),
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	for _, v := range ann.Data {
		strain, err := qrs.GetStockClient(registry.STOCK).
			GetStrain(ctx, &pb.StockId{Id: v.Attributes.EntryId})
		if err != nil {
			// errorutils.AddGQLError(ctx, err)
			qrs.Logger.Error(err)
			continue
		}
		strains = append(
			strains,
			stock.ConvertToStrainModel(strain.Data.Id, strain.Data.Attributes),
		)
	}
	/**
	  Some phenotypes list the same strain ID more than once. Consider a new approach
	  to de-duping this list while also keeping the Meta data from the annotations list.
	*/
	lm := int(ann.Meta.Limit)
	return &models.StrainListWithCursor{
		Strains:        strains,
		NextCursor:     int(ann.Meta.NextCursor),
		PreviousCursor: int(cur),
		Limit:          &lm,
		TotalCount:     len(ann.Data),
	}, nil
}

//nolint:dupl
func (qrs *QueryResolver) ListPlasmidsWithAnnotation(
	ctx context.Context,
	cursor *int,
	limit *int,
	typeArg string,
	annotation string,
) (*models.PlasmidListWithCursor, error) {
	plasmids := []*models.Plasmid{}
	cus := resolverutils.GetCursor(cursor)
	lmt := resolverutils.GetLimit(limit)
	onto := resolverutils.GetOntology(typeArg)
	ann, err := qrs.GetAnnotationClient(registry.ANNOTATION).
		ListAnnotations(ctx, &anno.ListParameters{
			Cursor: cus,
			Limit:  lmt,
			Filter: fmt.Sprintf("ontology==%s;tag==%s", onto, annotation),
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	for _, v := range ann.Data {
		plasmid, err := qrs.GetStockClient(registry.STOCK).
			GetPlasmid(ctx, &pb.StockId{Id: v.Attributes.EntryId})
		if err != nil {
			// errorutils.AddGQLError(ctx, err)
			qrs.Logger.Error(err)
			continue
		}
		plasmids = append(
			plasmids,
			stock.ConvertToPlasmidModel(
				plasmid.Data.Id,
				plasmid.Data.Attributes,
			),
		)
	}
	lm := int(ann.Meta.Limit)
	return &models.PlasmidListWithCursor{
		Plasmids:       plasmids,
		NextCursor:     int(ann.Meta.NextCursor),
		PreviousCursor: int(cus),
		Limit:          &lm,
		TotalCount:     len(ann.Data),
	}, nil
}

func (qrs *QueryResolver) AllStrains(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// ListRecentPlasmids is the resolver for the listRecentPlasmids field.
func (qrs *QueryResolver) ListRecentPlasmids(
	ctx context.Context,
	limit int,
) ([]*models.Plasmid, error) {
	return []*models.Plasmid{}, nil
}

// ListRecentStrains is the resolver for the listRecentStrains field.
func (qrs *QueryResolver) ListRecentStrains(
	ctx context.Context,
	limit int,
) ([]*models.Strain, error) {
	return []*models.Strain{}, nil
}

func (qrs *QueryResolver) listStrainsWithoutFilter(
	ctx context.Context,
	cus int64,
	lmt int64,
) (*models.StrainListWithCursor, error) {
	strainList, err := qrs.GetStockClient(registry.STOCK).
		ListStrains(ctx, &pb.StockParameters{Cursor: cus, Limit: lmt})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return &models.StrainListWithCursor{}, err
	}
	strains := make([]*models.Strain, 0)
	for _, sdata := range strainList.Data {
		strains = append(
			strains,
			stock.ConvertToStrainModel(sdata.Id, sdata.Attributes),
		)
	}
	qrs.Logger.Debugf(
		"successfully retrieved list of %d strains",
		strainList.Meta.Total,
	)
	limit := int(lmt)
	return &models.StrainListWithCursor{
		Limit:          &limit,
		NextCursor:     int(strainList.Meta.NextCursor),
		TotalCount:     int(strainList.Meta.Total),
		PreviousCursor: int(cus),
		Strains:        strains,
	}, nil
}

func (qrs *QueryResolver) toStrainModelList(
	strainList *pb.StrainCollection, limit int64, cursor int64,
) *models.StrainListWithCursor {
	smodelList := make([]*models.Strain, 0)
	for _, strain := range strainList.Data {
		smodelList = append(
			smodelList,
			stock.ConvertToStrainModel(strain.Id, strain.Attributes),
		)
	}

	lmt := int(limit)
	return &models.StrainListWithCursor{
		Strains:        smodelList,
		Limit:          &lmt,
		PreviousCursor: int(cursor),
		NextCursor:     int(strainList.Meta.NextCursor),
		TotalCount:     int(strainList.Meta.Total),
	}
}

func (qrs *QueryResolver) reportStrainListError(
	ctx context.Context,
	err error,
) (*models.StrainListWithCursor, error) {
	errorutils.AddGQLError(ctx, err)
	qrs.Logger.Error(err)
	return &models.StrainListWithCursor{}, err
}

// ListStrainsWithGene is the resolver for the listStrainsWithGene field.
func (qrs *QueryResolver) ListStrainsWithGene(
	ctx context.Context,
	gene string,
) ([]*models.Strain, error) {
	smodelList := make([]*models.Strain, 0)
	strainList, err := qrs.GetStockClient(registry.STOCK).
		ListStrains(ctx, &pb.StockParameters{
			Limit:  int64(100),
			Filter: fmt.Sprintf("gene@==%s", gene),
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return smodelList, nil
	}
	for _, strain := range strainList.Data {
		smodelList = append(
			smodelList,
			stock.ConvertToStrainModel(strain.Id, strain.Attributes),
		)
	}
	return smodelList, nil
}

func (mrs *MutationResolver) CreateGeneGeneralInfo(
	ctx context.Context,
	id string,
	input models.CreateGeneGeneralInfoInput,
) (*models.GeneGeneralInfo, error) {
	return &models.GeneGeneralInfo{}, fmt.Errorf("CreateGeneGeneralInfo is not yet implemented")
}

func (mrs *MutationResolver) UpdateGeneGeneralInfo(
	ctx context.Context,
	id string,
	input models.UpdateGeneGeneralInfoInput,
) (*models.GeneGeneralInfo, error) {
	return &models.GeneGeneralInfo{}, fmt.Errorf("UpdateGeneGeneralInfo is not yet implemented")
}

func (mrs *MutationResolver) UpdateStrainPhenotype(
	ctx context.Context,
	strainID string,
	target models.UpdateStrainPhenotypeTargetInput,
	payload models.UpdateStrainPhenotypePayloadInput,
) (*models.Strain, error) {
	return &models.Strain{}, fmt.Errorf("UpdateStrainPhenotype is not yet implemented")
}

func (mrs *MutationResolver) DeleteStrainPhenotype(
	ctx context.Context,
	strainID string,
	input models.DeleteStrainPhenotypeInput,
) (*models.DeleteStrainPhenotype, error) {
	return &models.DeleteStrainPhenotype{}, fmt.Errorf("DeleteStrainPhenotype is not yet implemented")
}

func (mrs *MutationResolver) AddStrainPhenotype(
	ctx context.Context,
	strainID string,
	input models.AddStrainPhenotypeInput,
) (*models.Strain, error) {
	return &models.Strain{}, fmt.Errorf("AddStrainPhenotype is not yet implemented")
}

func (qrs *QueryResolver) ListPhenotypeAssays(
	ctx context.Context,
	search string,
) ([]string, error) {
	return []string{}, fmt.Errorf("ListPhenotypeAssays is not yet implemented")
}

func (qrs *QueryResolver) ListPhenotypeEnvironments(
	ctx context.Context,
	search string,
) ([]string, error) {
	return []string{}, fmt.Errorf("ListPhenotypeEnvironments is not yet implemented")
}

func (qrs *QueryResolver) ListPhenotypes(
	ctx context.Context,
	search string,
) ([]string, error) {
	return []string{}, fmt.Errorf("ListPhenotypes is not yet implemented")
}
