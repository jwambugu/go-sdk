package test

import (
	"log"
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Messaging(t *testing.T) {
	var opts *elarian.Options
	opts.ApiKey = APIKey
	opts.AppId = AppId
	opts.OrgId = OrgId

	service, err := elarian.Initialize(opts)
	if err != nil {
		log.Fatal(err)
	}
	var customer *elarian.Customer
	customer.Id = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	customer.CustomerNumber = &elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
	}

	t.Run("It should send a text message", func(t *testing.T) {
		var messageBody *elarian.MessageBody
		messageBody.Text = "hello world"

		var channelNumber *elarian.MessagingChannelNumber
		channelNumber.Number = "41012"
		channelNumber.Channel = elarian.MESSAGING_CHANNEL_SMS

		response, err := service.SendMessage(
			customer,
			channelNumber,
			messageBody,
		)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.CustomerId).IsZero() {
			t.Errorf("Expected a customer id but didn't get any")
		}
		t.Logf("CustomerId %v", response.CustomerId)
	})
}
