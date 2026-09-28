package resolver

import (
	"context"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/emirpasic/gods/maps/hashmap"

	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrder(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ord := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	id := "999"
	o, err := ord.Order(context.Background(), id)
	assert.NoError(err, "expect no error from getting order information")
	assert.Exactly(o.Data.Id, id, "should match id")
	assert.Exactly(o.Data.Attributes.Courier, mocks.MockOrderAttributes.Courier, "should match courier")
	assert.Exactly(o.Data.Attributes.CourierAccount, mocks.MockOrderAttributes.CourierAccount, "should match courier account")
	assert.Exactly(o.Data.Attributes.Comments, mocks.MockOrderAttributes.Comments, "should match comments")
	assert.Exactly(o.Data.Attributes.Payment, mocks.MockOrderAttributes.Payment, "should match payment")
	assert.Exactly(o.Data.Attributes.PurchaseOrderNum, mocks.MockOrderAttributes.PurchaseOrderNum, "should match purchase order number")
	assert.Exactly(o.Data.Attributes.Status, order.OrderStatus_IN_PREPARATION, "should match status")
	assert.Exactly(o.Data.Attributes.Consumer, mocks.MockOrderAttributes.Consumer, "should match consumer")
	assert.Exactly(o.Data.Attributes.Payer, mocks.MockOrderAttributes.Payer, "should match payer")
	assert.Exactly(o.Data.Attributes.Purchaser, mocks.MockOrderAttributes.Purchaser, "should match purchaser")
	assert.ElementsMatch(o.Data.Attributes.Items, mocks.MockOrderAttributes.Items, "should match items")
}

func TestListOrders(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ord := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	filter := "type===strain"
	o, err := ord.ListOrders(context.Background(), &cursor, &limit, &filter)
	assert.NoError(err, "expect no error from getting list of orders")
	assert.Exactly(o.Limit, &limit, "should match limit")
	assert.Exactly(o.PreviousCursor, 0, "should match previous cursor")
	assert.Exactly(o.NextCursor, 10000, "should match next cursor")
	assert.Exactly(o.TotalCount, 3, "should match total count (length) of items")
	assert.Len(o.Orders, 3, "should have three orders")
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ord := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	comments := "first order"
	pon := "987654"
	id := "DBS123456"
	input := &models.CreateOrderInput{
		Courier:          "USPS",
		CourierAccount:   "123456",
		Comments:         &comments,
		Payment:          "credit",
		PurchaseOrderNum: &pon,
		Status:           models.StatusEnumInPreparation,
		Consumer:         "art@vandelayindustries.com",
		Payer:            "george@costanza.com",
		Purchaser:        "thatsgold@jerry.org",
		Items:            []string{id},
	}
	o, err := ord.CreateOrder(context.Background(), input)
	assert.NoError(err, "expect no error from creating an order")
	assert.Exactly(o.Data.Attributes.Courier, input.Courier, "should match courier")
	assert.Exactly(o.Data.Attributes.CourierAccount, input.CourierAccount, "should match courier account")
	assert.Exactly(&o.Data.Attributes.Comments, input.Comments, "should match comments")
	assert.Exactly(o.Data.Attributes.Payment, input.Payment, "should match payment")
	assert.Exactly(&o.Data.Attributes.PurchaseOrderNum, input.PurchaseOrderNum, "should match purchase order number")
	assert.Exactly(o.Data.Attributes.Status, order.OrderStatus_IN_PREPARATION, "should match status")
	assert.Exactly(o.Data.Attributes.Consumer, input.Consumer, "should match consumer")
	assert.Exactly(o.Data.Attributes.Payer, input.Payer, "should match payer")
	assert.Exactly(o.Data.Attributes.Purchaser, input.Purchaser, "should match purchaser")
	assert.ElementsMatch(o.Data.Attributes.Items, []string{"DBS123456"}, "should match items")
}

func TestUpdateOrder(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ord := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	courier := "FedEx"
	courierAccount := "444444"
	comments := "Please send ASAP"
	status := models.StatusEnumGrowing
	o, err := ord.UpdateOrder(
		context.Background(),
		"999",
		&models.UpdateOrderInput{
			Courier:        &courier,
			CourierAccount: &courierAccount,
			Comments:       &comments,
			Status:         &status,
		},
	)
	assert.NoError(err, "expect no error from updating an order")
	assert.Exactly(o.Data.Attributes.Courier, courier, "should match updated courier")
	assert.Exactly(o.Data.Attributes.CourierAccount, courierAccount, "should match updated courier account")
	assert.Exactly(o.Data.Attributes.Comments, comments, "should match updated comments")
	assert.Exactly(o.Data.Attributes.PurchaseOrderNum, "987654", "should match purchase order number")
	assert.Exactly(o.Data.Attributes.Status, order.OrderStatus_GROWING, "should match updated status")
	assert.Exactly(o.Data.Attributes.Payment, mocks.MockOrderAttributes.Payment, "should match existing payment")
	assert.Exactly(o.Data.Attributes.Consumer, mocks.MockOrderAttributes.Consumer, "should match existing consumer")
	assert.Exactly(o.Data.Attributes.Payer, mocks.MockOrderAttributes.Payer, "should match existing payer")
	assert.Exactly(o.Data.Attributes.Purchaser, mocks.MockOrderAttributes.Purchaser, "should match existing purchaser")
	assert.ElementsMatch(o.Data.Attributes.Items, mocks.MockOrderAttributes.Items, "should match existing items")
}

func TestCreateOrderWithUserInfo(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := new(clients.OrderServiceClient)
	var captured *order.NewOrder
	client.On(
		"CreateOrder",
		mock.Anything,
		mock.AnythingOfType("*order.NewOrder"),
	).Run(func(args mock.Arguments) {
		captured = args.Get(1).(*order.NewOrder)
	}).Return(&order.Order{
		Data: &order.Order_Data{
			Type:       "order",
			Id:         "999",
			Attributes: mocks.MockOrderAttributes,
		},
	}, nil)
	reg := &mocks.MockRegistry{ConnMap: hashmap.New()}
	reg.ConnMap.Put(registry.ORDER, client)
	ord := &MutationResolver{Registry: reg, Logger: mocks.TestLogger()}
	first := "Art"
	city := "New York"
	phone := "555-0100"
	input := &models.CreateOrderInput{
		Courier:        "USPS",
		CourierAccount: "123456",
		Payment:        "credit",
		Status:         models.StatusEnumInPreparation,
		Consumer:       "art@vandelayindustries.com",
		Payer:          "george@costanza.com",
		Purchaser:      "thatsgold@jerry.org",
		Items:          []string{"DBS123456"},
		ConsumerInfo:   &models.UserInfoInput{FirstName: &first, City: &city},
		PayerInfo:      &models.UserInfoInput{Phone: &phone},
	}
	_, err := ord.CreateOrder(context.Background(), input)
	assert.NoError(err, "expect no error from creating an order")
	attr := captured.Data.Attributes
	assert.Empty(attr.Comments, "absent comments should be empty")
	assert.Empty(attr.PurchaseOrderNum, "absent purchase order number should be empty")
	assert.Exactly(first, attr.ConsumerInfo.FirstName, "should map consumer first name")
	assert.Exactly(city, attr.ConsumerInfo.City, "should map consumer city")
	assert.Empty(attr.ConsumerInfo.LastName, "unset consumer field should be empty")
	assert.Exactly(phone, attr.PayerInfo.Phone, "should map payer phone")
	assert.Empty(attr.PayerInfo.FirstName, "unset payer field should be empty")
}

func TestUserInfoFromInput(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.Nil(userInfoFromInput(nil), "nil input should yield nil profile")
	empty := userInfoFromInput(&models.UserInfoInput{})
	assert.NotNil(empty, "empty input should yield empty profile")
	assert.Empty(empty.FirstName, "empty input should have empty fields")
	name := "Kramer"
	org := "Kramerica"
	addr1 := "129 W 81st"
	addr2 := "Apt 5B"
	state := "NY"
	zip := "10024"
	country := "USA"
	full := userInfoFromInput(&models.UserInfoInput{
		FirstName:     &name,
		LastName:      &name,
		Organization:  &org,
		FirstAddress:  &addr1,
		SecondAddress: &addr2,
		City:          &addr1,
		State:         &state,
		Zipcode:       &zip,
		Country:       &country,
		Phone:         &zip,
	})
	assert.Exactly(&order.UserInfo{
		FirstName:     name,
		LastName:      name,
		Organization:  org,
		FirstAddress:  addr1,
		SecondAddress: addr2,
		City:          addr1,
		State:         state,
		Zipcode:       zip,
		Country:       country,
		Phone:         zip,
	}, full, "should map every field")
}
