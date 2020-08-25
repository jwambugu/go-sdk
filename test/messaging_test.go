package test

import (
	"log"
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Messaging(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	var customer elarian.Customer
	customer.ID = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	t.Run("It should send a text message", func(t *testing.T) {
		var request elarian.SendMessageRequest
		request.AppID = AppID
		request.ProductID = ProductID
		request.Body.Text = "Hello World"
		request.ChannelNumber = elarian.MessagingChannelNumber{
			Number:  "41012",
			Channel: elarian.MessagingChannelSMS,
		}
		response, err := service.SendMessage(&customer, &request)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.CustomerId).IsZero() {
			t.Errorf("Expected a customer id but didn't get any")
		}
		t.Logf("CustomerID %v", response.CustomerId)
	})
}
