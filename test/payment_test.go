package test

import (
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_SendPayment(t *testing.T) {
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.PaymentRequest
	request.AppID = AppID
	request.ProductID = ProductID
	request.Cash = elarian.Cash{
		CurrencyCode: "KES",
		Amount:       100.00,
	}
	request.Channel = elarian.PaymentChannelNumber{
		Number:  "+254700000001",
		Channel: elarian.Telco,
	}

	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	res, err := service.SendPayment(&customer, &request)
	if err != nil {
		t.Error(err)
	}
	if res.Description == "" {
		t.Error("SendPayment: Expected a description but got none")
	}
	t.Logf("Result %v", res)
}

func Test_CheckoutPayment(t *testing.T) {
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.PaymentRequest
	request.AppID = AppID
	request.ProductID = ProductID
	request.Cash = elarian.Cash{
		CurrencyCode: "KES",
		Amount:       100.00,
	}
	request.Channel = elarian.PaymentChannelNumber{
		Number:  "+254700000001",
		Channel: elarian.Telco,
	}

	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	res, err := service.CheckoutPayment(&customer, &request)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.Status.Number())
	if res.CustomerId.Value != customer.ID {
		t.Errorf("CheckoutPayment: Expected customer id %s but got %s", customer.ID, res.CustomerId.Value)
	}
	t.Logf("Result %v", res)
}
