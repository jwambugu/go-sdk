package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func sendMessage(service elarian.Service) {
	var customer elarian.Customer
	customer.CustomerNumber = elarian.CustomerNumber{
		Number:   "customer_phone_number",
		Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
	}

	var message elarian.SendMessageRequest
	message.AppId = "app_id"
	message.ChannelNumber = elarian.MessagingChannelNumber{
		Number:  "channel_number",
		Channel: elarian.MESSAGING_CHANNEL_GOOGLE_RCS,
	}
	message.Body = elarian.MessageBody{
		Text: "Hello world",
	}

	response, err := service.SendMessage(&customer, &message)
	if err != nil {
		log.Fatalf("Could not send sms %v", err)
	}
	log.Printf("Customer id %v", response.GetCustomerId())
}
