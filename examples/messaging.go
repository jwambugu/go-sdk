package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func sendMessage(client elarian.Elarian) {
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "customer_phone_number",
		Provider: elarian.CustomerNumberProviderFacebook,
	}

	var message elarian.SendMessageRequest
	message.AppID = "app_id"
	message.ProductID = "product_id"
	message.ChannelNumber = elarian.MessagingChannelNumber{
		Number:  "channel_number",
		Channel: elarian.MessagingChannelGoogleRCS,
	}
	message.Body = elarian.MessageBody{
		Text: "Hello world",
	}

	response, err := client.SendMessage(&customer, &message)
	if err != nil {
		log.Fatalf("Could not send sms %v", err)
	}
	log.Printf("Customer id %v", response.GetCustomerId())
}
