package test

import (
	"log"
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_SendMessage(t *testing.T) {
	service, err := elarian.Initialize(APIKey, true)
	if err != nil {
		log.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

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
}
