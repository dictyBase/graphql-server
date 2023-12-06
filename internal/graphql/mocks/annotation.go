package mocks

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MockTagAnno(value, tag string) *annotation.TaggedAnnotation {
	return &annotation.TaggedAnnotation{
		Data: &annotation.TaggedAnnotation_Data{
			Type: "annotation",
			Id:   "123456",
			Attributes: &annotation.TaggedAnnotationAttributes{
				Value:     value,
				EntryId:   "DBS0236922",
				CreatedBy: "dsc@dictycr.org",
				CreatedAt: timestamppb.Now(),
				Tag:       tag,
				Ontology:  registry.DictyAnnoOntology,
				Version:   1,
			},
		},
	}
}

func MockTagList(
	onto, tag, value string,
) *annotation.TaggedAnnotationCollection_Data {
	return &annotation.TaggedAnnotationCollection_Data{
		Type: "annotation",
		Id:   "888888",
		Attributes: &annotation.TaggedAnnotationAttributes{
			Value:     value,
			EntryId:   "DBS0236922",
			CreatedBy: "dsc@dictycr.org",
			CreatedAt: timestamppb.Now(),
			Tag:       tag,
			Ontology:  onto,
			Version:   1,
		},
	}
}

func MockTagGroupAnno(
	onto, tag, value string,
) *annotation.TaggedAnnotationGroup_Data {
	return &annotation.TaggedAnnotationGroup_Data{
		Type: "annotation",
		Id:   "99999999",
		Attributes: &annotation.TaggedAnnotationAttributes{
			Version:   1,
			EntryId:   "DBS0235559",
			CreatedBy: "art@vandelay.org",
			CreatedAt: timestamppb.Now(),
			Ontology:  onto,
			Tag:       tag,
			Value:     value,
		},
	}
}

func MockNamesAnno() *annotation.TaggedAnnotationCollection {
	gcdata := []*annotation.TaggedAnnotationCollection_Data{}
	gcdata = append(
		gcdata,
		MockTagList(
			registry.DictyAnnoOntology,
			registry.SynTag,
			"catenin null",
		),
	)
	gcdata = append(
		gcdata,
		MockTagList(registry.DictyAnnoOntology, registry.SynTag, "aar1-"),
	)
	return &annotation.TaggedAnnotationCollection{
		Data: gcdata,
	}
}

func MockCharacteristicsAnno() *annotation.TaggedAnnotationCollection {
	gcdata := []*annotation.TaggedAnnotationCollection_Data{}
	gcdata = append(
		gcdata,
		MockTagList(
			registry.StrainCharOnto,
			"null mutant",
			registry.EmptyValue,
		),
	)
	gcdata = append(
		gcdata,
		MockTagList(registry.StrainCharOnto, "axenic", registry.EmptyValue),
	)
	return &annotation.TaggedAnnotationCollection{
		Data: gcdata,
	}
}

func MockPhenotypeAnno() *annotation.TaggedAnnotationGroupCollection {
	gcdata := []*annotation.TaggedAnnotationGroupCollection_Data{}
	gdata := []*annotation.TaggedAnnotationGroup_Data{}
	// gdata = append(gdata, MockTagGroupAnno(registry.DictyAnnoOntology, registry.LiteratureTag, "23967067"))
	gdata = append(
		gdata,
		MockTagGroupAnno(
			registry.PhenoOntology,
			"delayed culmination",
			registry.EmptyValue,
		),
	)
	gdata = append(
		gdata,
		MockTagGroupAnno(
			registry.AssayOntology,
			"confocal microscopy",
			registry.EmptyValue,
		),
	)
	gdata = append(
		gdata,
		MockTagGroupAnno(
			registry.EnvOntology,
			"in the presence of 8-Br-cAMP",
			registry.EmptyValue,
		),
	)
	gdata = append(
		gdata,
		MockTagGroupAnno(
			registry.DictyAnnoOntology,
			registry.NoteTag,
			"this is a test note",
		),
	)
	gcdata = append(gcdata, &annotation.TaggedAnnotationGroupCollection_Data{
		Type: "annotation_group",
		Group: &annotation.TaggedAnnotationGroup{
			Data:      gdata,
			GroupId:   "4924132",
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	})
	return &annotation.TaggedAnnotationGroupCollection{
		Data: gcdata,
	}
}

// MockPhenoCollectionAnno is used to test the ListStrainsWithPhenotype method.
func MockPhenoCollectionAnno(
	tag, id string,
) *annotation.TaggedAnnotationCollection_Data {
	return &annotation.TaggedAnnotationCollection_Data{
		Type: "annotation",
		Id:   "99999999",
		Attributes: &annotation.TaggedAnnotationAttributes{
			Version:   1,
			EntryId:   id,
			CreatedBy: "art@vandelay.org",
			CreatedAt: timestamppb.Now(),
			Ontology:  registry.PhenoOntology,
			Tag:       tag,
			Value:     registry.EmptyValue,
		},
	}
}

func MockPhenotypeListAnno() *annotation.TaggedAnnotationCollection {
	cdata := []*annotation.TaggedAnnotationCollection_Data{}
	cdata = append(
		cdata,
		MockPhenoCollectionAnno("delayed culmination", "DBS123456"),
	)
	cdata = append(
		cdata,
		MockPhenoCollectionAnno("delayed culmination", "DBS987654"),
	)
	cdata = append(
		cdata,
		MockPhenoCollectionAnno("delayed culmination", "DBS000001"),
	)
	cdata = append(
		cdata,
		MockPhenoCollectionAnno("delayed culmination", "DBS000001"),
	)
	return &annotation.TaggedAnnotationCollection{
		Data: cdata,
		Meta: &annotation.Meta{
			Limit: 10,
		},
	}
}

var MockSysNameAnno = MockTagAnno("DBS0236922", registry.SysnameTag)
var MockGenModAnno = MockTagAnno("exogenous mutation", registry.MuttypeTag)
var MockMutMethodAnno = MockTagAnno("Random Insertion", registry.MutmethodTag)

var MockGenotypeAnno = MockTagAnno(
	"axeA1,axeB1,axeC1,sadA-[sadA-KO],[pSadA-GFP],bsR,neoR",
	registry.GenoTag,
)
var MockInStockAnno = MockTagAnno("DBS0236922", registry.SysnameTag)

func MockedAnnotationClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"ListAnnotations",
		mock.Anything,
		mock.AnythingOfType("*annotation.ListParameters"),
	).Return(MockPhenotypeListAnno(), nil)
	return mockedAnnoClient
}

func MockedSysNameAnnoClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"GetEntryAnnotation",
		mock.Anything,
		mock.AnythingOfType("*annotation.EntryAnnotationRequest"),
	).Return(MockSysNameAnno, nil)
	return mockedAnnoClient
}

func MockedGenModClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"GetEntryAnnotation",
		mock.Anything,
		mock.AnythingOfType("*annotation.EntryAnnotationRequest"),
	).Return(MockGenModAnno, nil)
	return mockedAnnoClient
}

func MockedMutMethodClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"GetEntryAnnotation",
		mock.Anything,
		mock.AnythingOfType("*annotation.EntryAnnotationRequest"),
	).Return(MockMutMethodAnno, nil)
	return mockedAnnoClient
}

func MockedGenotypeClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"GetEntryAnnotation",
		mock.Anything,
		mock.AnythingOfType("*annotation.EntryAnnotationRequest"),
	).Return(MockGenotypeAnno, nil)
	return mockedAnnoClient
}

func MockedInStockClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"GetEntryAnnotation",
		mock.Anything,
		mock.AnythingOfType("*annotation.EntryAnnotationRequest"),
	).Return(MockInStockAnno, nil)
	return mockedAnnoClient
}

func MockedPhenotypeClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"ListAnnotationGroups",
		mock.Anything,
		mock.AnythingOfType("*annotation.ListGroupParameters"),
	).Return(MockPhenotypeAnno(), nil)
	return mockedAnnoClient
}

func MockedNamesClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"ListAnnotations",
		mock.Anything,
		mock.AnythingOfType("*annotation.ListParameters"),
	).Return(MockNamesAnno(), nil)
	return mockedAnnoClient
}

func MockedCharacteristicsClient() *clients.TaggedAnnotationServiceClient {
	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"ListAnnotations",
		mock.Anything,
		mock.AnythingOfType("*annotation.ListParameters"),
	).Return(MockCharacteristicsAnno(), nil)
	return mockedAnnoClient
}
