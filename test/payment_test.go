package test

import (
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Payments(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}

	var customer elarian.Customer
	customer.ID = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	t.Run("It should send a payment", func(t *testing.T) {
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

		res, err := service.SendPayment(&customer, &request)
		if err != nil {
			t.Error(err)
		}
		if res.Description == "" {
			t.Error("SendPayment: Expected a description but got none")
		}
		t.Logf("Result %v", res)
	})

	t.Run("It should checkout a payment", func(t *testing.T) {
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
		res, err := service.CheckoutPayment(&customer, &request)
		if err != nil {
			t.Error(err)
		}
		t.Log(res.Status.Number())
		if res.CustomerId.Value != customer.ID {
			t.Errorf("CheckoutPayment: Expected customer id %s but got %s", customer.ID, res.CustomerId.Value)
		}
		t.Logf("Result %v", res)
	})
}
