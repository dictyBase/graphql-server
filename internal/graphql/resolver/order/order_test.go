package order

import (
	"context"
	"testing"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/stretchr/testify/assert"
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
