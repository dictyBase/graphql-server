package mocks

import (
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
)

// func mockContent() *content.Content {
// 	return &content.Content{}
// }

func MockedContentClient() *clients.ContentServiceClient {
	mockedContentClient := new(clients.ContentServiceClient)
	return mockedContentClient
}
