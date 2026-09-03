package mocks

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ContentTimestamp = timestamppb.New(time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC))

var MockContentAttributes = &content.ContentAttributes{
	Name:      "test-page",
	Slug:      "news-test-page",
	CreatedBy: "art@vandelay.com",
	UpdatedBy: "h.e.@pennypacker.com",
	CreatedAt: ContentTimestamp,
	UpdatedAt: ContentTimestamp,
	Content:   "test content body",
	Namespace: "news",
}

func MockContentCollection() *content.ContentCollection {
	return &content.ContentCollection{
		Data: []*content.ContentCollection_Data{
			{Id: "100", Attributes: MockContentAttributes},
			{Id: "200", Attributes: MockContentAttributes},
			{Id: "300", Attributes: MockContentAttributes},
		},
	}
}

func MockedContentClient() *clients.ContentServiceClient {
	mockedContentClient := new(clients.ContentServiceClient)
	mockedContentClient.On("ListContents",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*content.ListParameters"),
	).Return(MockContentCollection(), nil)
	return mockedContentClient
}
