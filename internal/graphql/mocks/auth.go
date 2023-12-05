package mocks

import (
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/stretchr/testify/mock"
)

func MockedAuthClient() authentication.LogtoClient {
	mockedAuthClient := new(clients.LogtoClient)
	mockedAuthClient.On(
		"CheckUser",
		mock.AnythingOfType("string"),
	).Return(true, "", nil)
	return mockedAuthClient
}
