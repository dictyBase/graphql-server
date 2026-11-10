package order

import (
	"context"
	"testing"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserInfoResolvers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ord := &OrderResolver{Logger: mocks.TestLogger()}
	consumer := &pb.UserInfo{FirstName: "Art", City: "New York"}
	payer := &pb.UserInfo{Phone: "555-0100"}
	obj := &pb.Order{Data: &pb.Order_Data{Attributes: &pb.OrderAttributes{
		ConsumerInfo: consumer,
		PayerInfo:    payer,
	}}}
	got, err := ord.ConsumerInfo(context.Background(), obj)
	assert.NoError(err, "expect no error resolving consumer info")
	assert.Same(consumer, got, "should return embedded consumer profile")
	got, err = ord.PayerInfo(context.Background(), obj)
	assert.NoError(err, "expect no error resolving payer info")
	assert.Same(payer, got, "should return embedded payer profile")

	legacy := &pb.Order{Data: &pb.Order_Data{Attributes: &pb.OrderAttributes{}}}
	got, err = ord.ConsumerInfo(context.Background(), legacy)
	assert.NoError(err, "expect no error on legacy order")
	assert.Nil(got, "legacy order should have nil consumer profile")
	got, err = ord.PayerInfo(context.Background(), legacy)
	assert.NoError(err, "expect no error on legacy order")
	assert.Nil(got, "legacy order should have nil payer profile")
}

func TestOrderUserResolversAuthDisabled(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	ord := &OrderResolver{
		UserClient: authentication.NoopClient{},
		Logger:     mocks.TestLogger(),
	}
	obj := &pb.Order{Data: &pb.Order_Data{Attributes: &pb.OrderAttributes{
		Consumer:  "consumer@example.org",
		Payer:     "payer@example.org",
		Purchaser: "purchaser@example.org",
	}}}

	consumer, err := ord.Consumer(context.Background(), obj)
	req.NoError(err)
	req.Equal("Dicty", consumer.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", consumer.GetData().GetAttributes().GetLastName())

	payer, err := ord.Payer(context.Background(), obj)
	req.NoError(err)
	req.Equal("Dicty", payer.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", payer.GetData().GetAttributes().GetLastName())

	purchaser, err := ord.Purchaser(context.Background(), obj)
	req.NoError(err)
	req.Equal("Dicty", purchaser.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", purchaser.GetData().GetAttributes().GetLastName())
}
