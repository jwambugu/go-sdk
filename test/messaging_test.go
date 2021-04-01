package test

import (
	"log"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
	"github.com/stretchr/testify/assert"
)

func Test_Messaging(t *testing.T) {

	service, err := elarian.Connect(GetOpts())
	if err != nil {
		log.Fatal(err)
	}
	defer service.Disconnect()
	customer := &elarian.Customer{}
	customer.CustomerNumber = &elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderCellular,
	}

	t.Run("It should send a text message", func(t *testing.T) {
		messageBody := &elarian.MessageBody{}
		messageBody.Text = "hello world"

		channelNumber := &elarian.MessagingChannelNumber{}
		channelNumber.Number = "21356"
		channelNumber.Channel = elarian.MessagingChannelSms

		response, err := service.SendMessage(
			customer,
			channelNumber,
			messageBody,
		)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.NotEqual(t, response.CustomerId, "")
		// customerID messageId description status
		t.Logf("CustomerId %v", response.CustomerId.Value)
	})
}
