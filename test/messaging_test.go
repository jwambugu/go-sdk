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
	customer := &elarian.Customer{
		CustomerNumber: &elarian.CustomerNumber{
			Number:   "+254712876967",
			Provider: elarian.CustomerNumberProviderCellular,
		},
	}

	t.Run("It should send a text message", func(t *testing.T) {
		channelNumber := &elarian.MessagingChannelNumber{
			Number:  "21356",
			Channel: elarian.MessagingChannelSms,
		}
		response, err := service.SendMessage(
			customer.CustomerNumber,
			channelNumber,
			elarian.TextMessage("Hello World"),
		)

		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.NotEqual(t, response.CustomerID, "")
	})
}
