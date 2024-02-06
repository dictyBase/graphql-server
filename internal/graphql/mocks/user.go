package mocks

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/stretchr/testify/mock"
)

func mockUser() *user.User {
	return &user.User{
		Data: &user.UserData{
			Id: 999,
			Attributes: &user.UserAttributes{
				Email:     "kenny@bania.com",
				FirstName: "Kenny",
				LastName:  "Bania",
				City:      "New York City",
				State:     "NY",
				Country:   "United States",
				IsActive:  true,
			},
		},
	}
}

func MockedUserClient() *clients.UserServiceClient {
	mockedUserClient := new(clients.UserServiceClient)
	mockedUserClient.On(
		"GetUserByEmail",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*jsonapi.GetEmailRequest"),
	).Return(mockUser(), nil)
	return mockedUserClient
}
