package stock

import (
	"context"
	"fmt"
	"regexp"

	"github.com/dictyBase/aphgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
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
	UserClient       user.UserServiceClient
	AnnotationClient annotation.TaggedAnnotationServiceClient
	Registry         registry.Registry
	Logger           *logrus.Entry
}

func (r *StrainResolver) CreatedBy(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	u, err := getUserByEmail(ctx, r.UserClient, obj.CreatedBy)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), err
	}
	return u, nil
}

func (r *StrainResolver) UpdatedBy(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	u, err := getUserByEmail(ctx, r.UserClient, obj.UpdatedBy)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), err
	}
	return u, nil
}

func (r *StrainResolver) Depositor(
	ctx context.Context,
	obj *models.Strain,
) (*user.User, error) {
	d, err := getUserByEmail(ctx, r.UserClient, obj.Depositor)
	if err != nil {
		r.Logger.Error(err)
		return newUser(), nil
	}
	return d, nil
}

func (r *StrainResolver) Genes(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Gene, error) {
	gntype := []*models.Gene{}
	redis := r.Registry.GetRedisRepository(cache.RedisKey)
	for _, gne := range obj.Genes {
		gene, err := cache.GetGeneFromCache(ctx, redis, gne)
		if err != nil {
			r.Logger.Error(err)
			continue
		}
		gntype = append(gntype, gene)
	}
	return gntype, nil
}

func (r *StrainResolver) Publications(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Publication, error) {
	pubs := make([]*models.Publication, 0)
	for _, id := range obj.Publications {
		// GWDI IDs come back as 10.1101/582072 or doi:10.1101/582072
		if doiRgxp.MatchString(id) {
			p, err := fetch.FetchDOI(ctx, fmt.Sprintf("https://doi.org/%s", id))
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				r.Logger.Error(err)
				return pubs, err
			}
			pubs = append(pubs, p)
		} else {
			endpoint := r.Registry.GetAPIEndpoint(registry.PUBLICATION)
			p, err := fetch.FetchPublication(ctx, endpoint, id)
			if err != nil {
				errorutils.AddGQLError(ctx, err)
				r.Logger.Error(err)
				return pubs, err
			}
			pubs = append(pubs, p)
		}
	}
	return pubs, nil
}

func (r *StrainResolver) Parent(
	ctx context.Context,
	obj *models.Strain,
) (*models.Strain, error) {
	parent := obj.Parent
	if parent == nil {
		return &models.Strain{}, nil
	}
	n, err := r.Client.GetStrain(ctx, &pb.StockId{Id: *parent})
	if err != nil {
		r.Logger.Debugf("could not find parent strain with ID %s", *parent)
		return nil, nil
	}
	r.Logger.Debugf("successfully found parent strain with ID %s", *parent)
	return ConvertToStrainModel(*parent, n.Data.Attributes), nil
}

func (r *StrainResolver) Names(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	names := make([]string, 0)
	for _, v := range obj.Names {
		names = append(names, v)
	}
	n, err := r.AnnotationClient.ListAnnotations(
		ctx,
		&annotation.ListParameters{
			Limit: 20,
			Filter: fmt.Sprintf(
				"entry_id===%s;tag===%s;ontology===%s",
				obj.ID, registry.SynTag, registry.DictyAnnoOntology,
			)})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return names, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return names, err
	}
	for _, syn := range n.Data {
		names = append(names, syn.Attributes.Value)
	}
	return names, nil
}

func (r *StrainResolver) Phenotypes(
	ctx context.Context,
	obj *models.Strain,
) ([]*models.Phenotype, error) {
	p := []*models.Phenotype{}
	strainId := obj.ID
	gc, err := r.AnnotationClient.ListAnnotationGroups(
		ctx,
		&annotation.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;ontology==%s",
				strainId,
				registry.PhenoOntology,
			),
			Limit: 30,
		})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return p, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return p, err
	}
	p = getPhenotypes(ctx, r, gc.Data)
	return p, nil
}

func (r *StrainResolver) GeneticModification(
	ctx context.Context,
	obj *models.Strain,
) (*string, error) {
	var gm string
	gc, err := r.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.MuttypeTag,
			Ontology: registry.DictyAnnoOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return &gm, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return &gm, err
	}
	gm = gc.Data.Attributes.Value
	return &gm, nil
}

func (r *StrainResolver) MutagenesisMethod(
	ctx context.Context,
	obj *models.Strain,
) (*string, error) {
	var m string
	gc, err := r.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.MutmethodTag,
			Ontology: registry.MutagenesisOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return &m, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return &m, err
	}
	m = gc.Data.Attributes.Value
	return &m, nil
}

func (r *StrainResolver) SystematicName(
	ctx context.Context,
	obj *models.Strain,
) (string, error) {
	sn, err := r.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.SysnameTag,
			Ontology: registry.DictyAnnoOntology,
			EntryId:  obj.ID,
		},
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return "", nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return "", err
	}
	return sn.Data.Attributes.Value, nil
}

func (r *StrainResolver) Characteristics(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	pslice := make([]string, 0)
	cg, err := r.AnnotationClient.ListAnnotations(
		ctx, &annotation.ListParameters{Filter: fmt.Sprintf(
			"entry_id===%s;ontology===%s",
			obj.ID, registry.StrainCharOnto,
		)},
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return pslice, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return pslice, err
	}
	for _, item := range cg.Data {
		pslice = append(pslice, item.Attributes.Tag)
	}
	return pslice, nil
}

func (r *StrainResolver) Genotypes(
	ctx context.Context,
	obj *models.Strain,
) ([]string, error) {
	gntype := make([]string, 0)
	gl, err := r.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			EntryId:  obj.ID,
			Ontology: registry.DictyAnnoOntology,
			Tag:      registry.GenoTag,
		})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return gntype, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Error(err)
		return gntype, err
	}
	gntype = append(gntype, gl.Data.Attributes.Value)
	return gntype, nil
}

func (r *StrainResolver) InStock(
	ctx context.Context,
	obj *models.Strain,
) (bool, error) {
	id := obj.ID
	_, err := r.AnnotationClient.GetEntryAnnotation(
		ctx,
		&annotation.EntryAnnotationRequest{
			Tag:      registry.StrainInvTag,
			Ontology: registry.StrainInvOnto,
			EntryId:  id,
		},
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return false, nil
		}
		errorutils.AddGQLError(ctx, err)
		r.Logger.Errorf("error getting %s from inventory: %v", id, err)
		return false, err
	}
	return true, nil
}

func getPhenotypes(
	ctx context.Context,
	r *StrainResolver,
	data []*annotation.TaggedAnnotationGroupCollection_Data,
) []*models.Phenotype {
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
					endpoint := r.Registry.GetAPIEndpoint(registry.PUBLICATION)
					pub, err := fetch.FetchPublication(
						ctx,
						endpoint,
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

func getUserByEmail(
	ctx context.Context,
	uc user.UserServiceClient,
	email string,
) (*user.User, error) {
	u := &user.User{}
	if email == "" {
		return u, fmt.Errorf("got an empty email address %s", email)
	}
	gntype, err := uc.GetUserByEmail(
		ctx,
		&jsonapi.GetEmailRequest{Email: email},
	)
	if err != nil {
		errorutils.AddGQLError(
			ctx,
			fmt.Errorf("could not find user with email %s", email),
		)
		return u, err
	}
	return gntype, nil
}

func newUser() *user.User {
	return &user.User{
		Data: &user.UserData{
			Attributes: &user.UserAttributes{
				FirstName:     "",
				LastName:      "",
				Email:         "",
				Organization:  "",
				GroupName:     "",
				FirstAddress:  "",
				SecondAddress: "",
				City:          "",
				State:         "",
				Zipcode:       "",
				Country:       "",
				Phone:         "",
				IsActive:      false,
			},
		},
	}
}
