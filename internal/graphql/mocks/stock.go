package mocks

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	t              = time.Date(2020, time.January, 0o1, 0o1, 0, 0, 0, time.UTC)
	StockTimestamp = timestamppb.New(t)
)

var MockPlasmidAttributes = &stock.PlasmidAttributes{
	CreatedAt:       StockTimestamp,
	UpdatedAt:       StockTimestamp,
	CreatedBy:       "art@vandelay.com",
	UpdatedBy:       "art@vandelay.com",
	Summary:         "test summary",
	EditableSummary: "editable test summary",
	Genes:           []string{"DDB_G0285425"},
	Dbxrefs:         []string{"test1"},
	Publications:    []string{"99999"},
	ImageMap:        "https://eric.dictybase.dev/test.jpg",
	Sequence:        "ABCDEF",
	Name:            "pTest",
}

var MockStrainAttributes = &stock.StrainAttributes{
	CreatedAt:       StockTimestamp,
	UpdatedAt:       StockTimestamp,
	CreatedBy:       "art@vandelay.com",
	UpdatedBy:       "art@vandelay.com",
	Summary:         "test summary",
	EditableSummary: "editable test summary",
	Genes:           []string{"DDB_G0285425"},
	Dbxrefs:         []string{"test1"},
	Publications:    []string{"99999"},
	Label:           "test99",
	Species:         "human",
	Plasmid:         "pTest",
	Names:           []string{"fusilli"},
}

func MockStrain() *stock.Strain {
	return &stock.Strain{
		Data: &stock.Strain_Data{
			Type:       "strain",
			Id:         "DBS123456",
			Attributes: MockStrainAttributes,
		},
	}
}

func MockStrainInputWithParams(
	depositor, parent string,
) *stock.StrainAttributes {
	s := MockStrainAttributes
	s.Depositor = depositor
	s.Parent = parent
	return s
}

func MockPlasmidInputWithParams(depositor string) *stock.PlasmidAttributes {
	p := MockPlasmidAttributes
	p.Depositor = depositor
	return p
}

var MockUpdateStrainAttributes = &stock.StrainAttributes{
	CreatedAt:       StockTimestamp,
	UpdatedAt:       StockTimestamp,
	CreatedBy:       "art@vandelay.com",
	UpdatedBy:       "h.e.@pennypacker.com",
	Summary:         "updated summary",
	EditableSummary: "editable updated summary",
	Depositor:       "puddy@nyrangers.com",
	Genes:           []string{"sadA"},
	Dbxrefs:         []string{"test1"},
	Publications:    []string{"99999"},
	Label:           "test99",
	Species:         "human",
	Plasmid:         "pTest",
	Names:           []string{"fusilli"},
}

var MockUpdatePlasmidAttributes = &stock.PlasmidAttributes{
	CreatedAt:       StockTimestamp,
	UpdatedAt:       StockTimestamp,
	CreatedBy:       "art@vandelay.com",
	UpdatedBy:       "h.e.@pennypacker.com",
	Summary:         "updated summary",
	EditableSummary: "editable updated summary",
	Depositor:       "puddy@nyrangers.com",
	Genes:           []string{"sadA"},
	Dbxrefs:         []string{"test1"},
	Publications:    []string{"99999"},
	ImageMap:        "https://eric.dictybase.dev/test.jpg",
	Sequence:        "ABCDEF",
	Name:            "pTest",
}

var mockStrainData = &stock.StrainCollection_Data{
	Type:       "strain",
	Id:         "DBS123456",
	Attributes: MockStrainAttributes,
}

var mockPlasmidList = &stock.PlasmidCollection_Data{
	Type:       "plasmid",
	Id:         "DBP123456",
	Attributes: MockPlasmidAttributes,
}

func MockStrainCollection() *stock.StrainCollection {
	var strains []*stock.StrainCollection_Data
	for i := 0; i < 3; i++ {
		strains = append(strains, mockStrainData)
	}
	return &stock.StrainCollection{
		Data: strains,
		Meta: &stock.Meta{
			Limit: 10,
			Total: int64(len(strains)),
		},
	}
}

func MockStrainList() *stock.StrainList {
	return &stock.StrainList{
		Data: []*stock.StrainList_Data{
			{Type: "strain", Id: "DBS000001", Attributes: MockStrainAttributes},
			{Type: "strain", Id: "DBS000002", Attributes: MockStrainAttributes},
		},
	}
}

func MockPlasmidCollection() *stock.PlasmidCollection {
	var plasmids []*stock.PlasmidCollection_Data
	plasmids = append(
		plasmids,
		mockPlasmidList,
		mockPlasmidList,
		mockPlasmidList,
	)
	return &stock.PlasmidCollection{
		Data: plasmids,
		Meta: &stock.Meta{
			Limit: 10,
			Total: int64(len(plasmids)),
		},
	}
}

func mockPlasmidWithParams(depositor string) *stock.Plasmid {
	attr := MockPlasmidAttributes
	attr.Depositor = depositor
	return &stock.Plasmid{
		Data: &stock.Plasmid_Data{
			Type:       "plasmid",
			Id:         "DBP123456",
			Attributes: attr,
		},
	}
}

func mockStrainWithParams(depositor, parent string) *stock.Strain {
	attr := MockStrainAttributes
	attr.Depositor = depositor
	attr.Parent = parent
	return &stock.Strain{
		Data: &stock.Strain_Data{
			Type:       "strain",
			Id:         "DBS123456",
			Attributes: attr,
		},
	}
}

func mockUpdatePlasmid() *stock.Plasmid {
	return &stock.Plasmid{
		Data: &stock.Plasmid_Data{
			Type:       "plasmid",
			Id:         "DBP123456",
			Attributes: MockUpdatePlasmidAttributes,
		},
	}
}

func mockUpdateStrain() *stock.Strain {
	return &stock.Strain{
		Data: &stock.Strain_Data{
			Type:       "strain",
			Id:         "DBS123456",
			Attributes: MockUpdateStrainAttributes,
		},
	}
}

func MockedStockClient() *clients.StockServiceClient {
	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On("GetPlasmid",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockId"),
	).Return(mockPlasmidWithParams("kenny@bania.com"), nil).On("GetStrain",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&stock.StockId{Id: "DBS987654"},
	).Return(mockStrainWithParams("kenny@bania.com", "DBS987654"), nil).On("GetStrain",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&stock.StockId{Id: "DBS123456"},
	).Return(MockStrain(), nil).
		On("GetStrain",
			mock.AnythingOfType("*context.emptyCtx"),
			&stock.StockId{Id: "DBS000001"},
		).
		Return(MockStrain(), nil).
		On("GetStrain",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			&stock.StockId{Id: "DBS000002"},
		).
		Return(MockStrain(), nil).
		On("ListStrains",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.StockParameters"),
		).Return(MockStrainCollection(), nil).
		On("ListPlasmids",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.StockParameters"),
		).
		Return(MockPlasmidCollection(), nil).
		On("CreateStrain",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.NewStrain"),
		).
		Return(mockStrainWithParams("kenny@bania.com", "DBS987654"), nil).
		On("CreatePlasmid",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.NewPlasmid"),
		).
		Return(mockPlasmidWithParams("kenny@bania.com"), nil).
		On("UpdatePlasmid",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.PlasmidUpdate"),
		).
		Return(mockUpdatePlasmid(), nil).
		On("UpdateStrain",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.StrainUpdate"),
		).
		Return(mockUpdateStrain(), nil).
		On("ListStrainsByIds",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(MockStrainList(), nil)
	return mockedStockClient
}
