package mocks

import (
	"fmt"
	"os"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/dictyBase/graphql-server/internal/repository/redis"
	"github.com/emirpasic/gods/maps/hashmap"
	"google.golang.org/grpc"
)

type MockRegistry struct {
	ConnMap *hashmap.Map
}

func (mr *MockRegistry) AddAuthClient(
	key string, auth authentication.LogtoClient,
) {
	mr.ConnMap.Put(key, auth)
}

func (mr *MockRegistry) AddAPIEndpoint(key, endpoint string) {
	mr.ConnMap.Put(key, endpoint)
}

func (mr *MockRegistry) AddAPIConnection(key string, conn *grpc.ClientConn) {
	mr.ConnMap.Put(key, conn)
}

func (mr *MockRegistry) AddRepository(key string, st repository.Repository) {
	mr.ConnMap.Put(key, st)
}

// GetAPIClient looks up a client in the hashmap
func (mr *MockRegistry) GetAPIConnection(key string) (conn *grpc.ClientConn) {
	v, ok := mr.ConnMap.Get(key)
	if !ok {
		panic("could not get grpc client connection")
	}
	conn, _ = v.(*grpc.ClientConn)
	return conn
}

func (mr *MockRegistry) GetUserClient(key string) user.UserServiceClient {
	return MockedUserClient()
}

func (mr *MockRegistry) GetRoleClient(key string) user.RoleServiceClient {
	return MockedRoleClient()
}

func (mr *MockRegistry) GetPermissionClient(
	key string,
) user.PermissionServiceClient {
	return MockedPermissionClient()
}

func (mr *MockRegistry) GetStockClient(key string) stock.StockServiceClient {
	return MockedStockClient()
}

func (mr *MockRegistry) GetOrderClient(key string) order.OrderServiceClient {
	return MockedOrderClient()
}

func (mr *MockRegistry) GetContentClient(
	key string,
) content.ContentServiceClient {
	return MockedContentClient()
}

func (mr *MockRegistry) GetAnnotationClient(
	key string,
) annotation.TaggedAnnotationServiceClient {
	return MockedAnnotationClient()
}

func (mr *MockRegistry) GetAuthClient(key string) authentication.LogtoClient {
	return MockedAuthClient()
}

func (mr *MockRegistry) GetIdentityClient(
	key string,
) identity.IdentityServiceClient {
	return MockedIdentityClient()
}

func (mr MockRegistry) GetAPIEndpoint(key string) string {
	v, _ := mr.ConnMap.Get(key)
	endpoint, _ := v.(string)
	return endpoint
}

func (mr MockRegistry) GetRedisRepository(key string) repository.Repository {
	radd := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_SERVICE_HOST"),
		os.Getenv("REDIS_SERVICE_PORT"),
	)
	c, _ := redis.NewCache(radd)
	return c
}

func (mr *MockRegistry) ServiceMap() map[string]string {
	return map[string]string{
		"stock":   "stock",
		"user":    "user",
		"content": "content",
	}
}
