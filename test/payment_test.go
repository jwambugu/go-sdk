package test

import (
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Payments(t *testing.T) {
	var opts *elarian.Options
	opts.ApiKey = APIKey
	opts.AppId = AppId
	opts.OrgId = OrgId

	service, err := elarian.Initialize(opts)
	if err != nil {
		t.Fatal(err)
	}

	var customer elarian.Customer
	customer.Id = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	customer.CustomerNumber = &elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
	}

	t.Run("It should send a payment", func(t *testing.T) {
		var request *elarian.Paymentrequest
		request.Cash = elarian.Cash{
			CurrencyCode: "KES",
			Amount:       100.00,
		}
		request.Channel = elarian.PaymentChannelNumber{
			Number:  "+254700000001",
			Channel: elarian.PAYMENT_CHANNEL_TELCO,
		}

		res, err := service.InitiatePayment(&customer, request)
		if err != nil {
			t.Error(err)
		}
		if res.Description == "" {
			t.Error("SendPayment: Expected a description but got none")
		}
		t.Logf("Result %v", res)
	})
}
