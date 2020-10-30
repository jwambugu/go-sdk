package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func sendMessage(service elarian.Service) {
	customer := &elarian.Customer{
		CustomerNumber: elarian.CustomerNumber{
			Number:   "customer_phone_number",
			Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
		},
	}
	channelNumber := &elarian.MessagingChannelNumber{
		Number:  "channel_number",
		Channel: elarian.MESSAGING_CHANNEL_GOOGLE_RCS,
	}
	body := &elarian.MessageBody{
		Text: "Hello world",
	}
	response, err := service.SendMessage(customer, body, channelNumber)
	if err != nil {
		log.Fatalf("Could not send sms %v", err)
	}
	log.Printf("Customer id %v", response.GetCustomerId())
}
