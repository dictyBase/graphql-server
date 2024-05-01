package stock

import (
	"context"
	"fmt"
	"regexp"

	"github.com/dictyBase/aphgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
)

var doiRgxp = regexp.MustCompile(`^(doi:)?10.\d{4,9}/[-._;()/:A-Z0-9]+$`)

type StrainResolver struct {
	Client           pb.StockServiceClient
	UserClient       authentication.LogtoClient
	AnnotationClient annotation.TaggedAnnotationServiceClient
	Registry         registry.Registry
	Logger           *logrus.Entry
}

func (srs *StrainResolver) CreatedBy(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	u, err := userByEmail(ctx, obj.CreatedBy, srs.UserClient, srs.Logger)
	if err != nil {
		srs.Logger.Error(err)
		return nil, err
	}
	return u, nil
}

func (srs *StrainResolver) UpdatedBy(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	u, err := userByEmail(ctx, obj.UpdatedBy, srs.UserClient, srs.Logger)
	if err != nil {
		srs.Logger.Error(err)
		return nil, err
	}
	return u, nil
}

func (srs *StrainResolver) Depositor(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	d, err := userByEmail(ctx, obj.Depositor, srs.UserClient, srs.Logger)
	if err != nil {
		srs.Logger.Error(err)
		return nil, err
	}
	return d, nil
}

func (srs *StrainResolver) Genes(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Gene, error) {
	gntype := []*models.Gene{}
	redis := srs.Registry.GetRedisRepository(cache.RedisKey)
	for _, gne := range obj.Genes {
		gene, err := cache.GetGeneFromCache(ctx, redis, gne)
		if err != nil {
			srs.Logger.Error(err)
			continue
		}
		gntype = append(gntype, gene)
	}
	return gntype, nil
}

func (srs *StrainResolver) Publications(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Publication, error) {
	redis := srs.Registry.GetRedisRepository(cache.RedisKey)
	pubs := make([]*models.Publication, 0)
	for _, id := range obj.Publications {
		// GWDI IDs come back as 10.1101/582072 or doi:10.1101/582072
		if doiRgxp.MatchString(id) {
			p, err := fetch.FetchDOI(ctx, fmt.Sprintf("https://doi.org/%s", id))
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				srs.Logger.Error(err)
				return pubs, err
			}
			pubs = append(pubs, p)
		} else {
			p, err := fetch.FetchPublicationFromEuroPMC(
				ctx,
				redis,
				srs.Registry.GetAPIEndpoint(registry.PUBLICATION),
				id,
			)
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				srs.Logger.Error(err)
				return pubs, err
			}
			pubs = append(pubs, p)
		}
	}
	return pubs, nil
}

func (srs *StrainResolver) Parent(
	ctx context.Context,
	obj *models.Strain,
) (*models.Strain, error) {
	parent := obj.Parent
	if parent == nil {
		return &models.Strain{}, nil
	}
	n, err := srs.Client.GetStrain(ctx, &pb.StockId{Id: *parent})
	if err != nil {
		srs.Logger.Debugf("could not find parent strain with ID %s", *parent)
		return nil, nil
	}
	srs.Logger.Debugf("successfully found parent strain with ID %s", *parent)
	return ConvertToStrainModel(*parent, n.Data.Attributes), nil
}

func (srs *StrainResolver) Names(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	names := make([]string, 0)
	names = append(names, obj.Names...)
	n, err := srs.AnnotationClient.ListAnnotations(
		ctx,
		&annotation.ListParameters{
			Limit: 20,
			Filter: fmt.Sprintf(
				"entry_id===%s;tag===%s;ontology===%s",
				obj.ID, registry.SynTag, registry.DictyAnnoOntology,
			)})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return names, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return names, err
	}
	for _, syn := range n.Data {
		names = append(names, syn.Attributes.Value)
	}
	return names, nil
}

func (srs *StrainResolver) Phenotypes(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Phenotype, error) {
	p := []*models.Phenotype{}
	strainID := obj.ID
	gc, err := srs.AnnotationClient.ListAnnotationGroups(
		ctx,
		&annotation.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;ontology==%s",
				strainID,
				registry.PhenoOntology,
			),
			Limit: 30,
		})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return p, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return p, err
	}
	p = getPhenotypes(ctx, srs, gc.Data)
	return p, nil
}

func (srs *StrainResolver) GeneticModification(
	ctx context.Context,
	obj *models.Strain,
) (*string, error) {
	var gm string
	gc, err := srs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.MuttypeTag,
			Ontology: registry.DictyAnnoOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &gm, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return &gm, err
	}
	gm = gc.Data.Attributes.Value
	return &gm, nil
}

func (srs *StrainResolver) MutagenesisMethod(
	ctx context.Context,
	obj *models.Strain,
) (*string, error) {
	var method string
	geneticChange, err := srs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.MutmethodTag,
			Ontology: registry.MutagenesisOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &method, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return &method, err
	}
	method = geneticChange.Data.Attributes.Value
	return &method, nil
}

func (srs *StrainResolver) SystematicName(
	ctx context.Context,
	obj *models.Strain,
) (string, error) {
	sn, err := srs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.SysnameTag,
			Ontology: registry.DictyAnnoOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return "", err
	}
	return sn.Data.Attributes.Value, nil
}

func (srs *StrainResolver) Characteristics(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	pslice := make([]string, 0)
	cg, err := srs.AnnotationClient.ListAnnotations(
		ctx, &annotation.ListParameters{Filter: fmt.Sprintf(
			"entry_id===%s;ontology===%s",
			obj.ID, registry.StrainCharOnto,
		)},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return pslice, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return pslice, err
	}
	for _, item := range cg.Data {
		pslice = append(pslice, item.Attributes.Tag)
	}
	return pslice, nil
}

func (srs *StrainResolver) Genotypes(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	gntype := make([]string, 0)
	gl, err := srs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			EntryId:  obj.ID,
			Ontology: registry.DictyAnnoOntology,
			Tag:      registry.GenoTag,
		})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return gntype, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Error(err)
		return gntype, err
	}
	gntype = append(gntype, gl.Data.Attributes.Value)
	return gntype, nil
}

func (srs *StrainResolver) InStock(
	ctx context.Context,
	obj *models.Strain,
) (bool, error) {
	id := obj.ID
	_, err := srs.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.StrainInvTag,
			Ontology: registry.StrainInvOnto,
			EntryId:  id,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		errorutils.AddGQLError(ctx, err)
		srs.Logger.Errorf("error getting %s from inventory: %v", id, err)
		return false, err
	}
	return true, nil
}

func getPhenotypes(
	ctx context.Context,
	r *StrainResolver,
	data []*annotation.TaggedAnnotationGroupCollection_Data,
) []*models.Phenotype {
	redis := r.Registry.GetRedisRepository(cache.RedisKey)
	p := []*models.Phenotype{}
	for _, item := range data {
		m := &models.Phenotype{}
		for _, gntype := range item.Group.Data {
			switch gntype.Attributes.Ontology {
			case registry.PhenoOntology:
				m.Phenotype = gntype.Attributes.Tag
			case registry.EnvOntology:
				m.Environment = &gntype.Attributes.Tag
			case registry.AssayOntology:
				m.Assay = &gntype.Attributes.Tag
			case registry.DictyAnnoOntology:
				if gntype.Attributes.Tag == registry.LiteratureTag {
					pub, err := fetch.FetchPublicationFromEuroPMC(
						ctx,
						redis,
						r.Registry.GetAPIEndpoint(registry.PUBLICATION),
						gntype.Attributes.Value,
					)
					if err != nil {
						r.Logger.Error(err)
						errorutils.AddGQLError(ctx, err)
					}
					m.Publication = pub
				}
				if gntype.Attributes.Tag == registry.NoteTag {
					m.Note = &gntype.Attributes.Value
				}
			}
		}
		p = append(p, m)
	}
	return p
}

func ConvertToStrainModel(id string, attr *pb.StrainAttributes) *models.Strain {
	return &models.Strain{
		ID:              id,
		CreatedAt:       aphgrpc.ProtoTimeStamp(attr.CreatedAt),
		UpdatedAt:       aphgrpc.ProtoTimeStamp(attr.UpdatedAt),
		Label:           attr.Label,
		Species:         attr.Species,
		Summary:         &attr.Summary,
		EditableSummary: &attr.EditableSummary,
		Plasmid:         &attr.Plasmid,
		Parent:          &attr.Parent,
		Dbxrefs:         attr.Dbxrefs,
		Names:           attr.Names,
		LazyStock: models.LazyStock{
			CreatedBy:    attr.CreatedBy,
			UpdatedBy:    attr.UpdatedBy,
			Depositor:    attr.Depositor,
			Genes:        attr.Genes,
			Publications: attr.Publications,
		},
	}
}

func sliceConverter[T any](aslice []T) []*T {
	// make a copy so that the passed slice gets
	// garbage collected
	// https://stackoverflow.com/a/48459980
	nslice := make([]T, 0)
	nslice = append(nslice, aslice...)
	pslice := make([]*T, 0)
	for i := range nslice {
		pslice = append(pslice, &nslice[i])
	}

	return pslice
}
