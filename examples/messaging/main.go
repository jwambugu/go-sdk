package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func sendMessage(service elarian.Service) {
	var custNumber elarian.CustomerNumber
	custNumber.Number = "+254708752702"
	custNumber.Provider = elarian.CUSTOMER_NUMBER_PROVIDER_TELCO

	var cust elarian.Customer
	cust.CustomerNumber = &custNumber

	var channel elarian.MessagingChannelNumber
	channel.Channel = elarian.MESSAGING_CHANNEL_SMS
	channel.Number = "Elarian"

	var messageBody elarian.MessageBody
	messageBody.Text = "Hello world from the go sdk"

	response, err := service.SendMessage(&cust, &channel, &messageBody)
	if err != nil {
		log.Fatalf("Message not send %v \n", err)
	}
	log.Printf("Status %d Description %s \n customerID %s \n", response.Status, response.Description, response.CustomerId)
}

func main() {
	service, err := elarian.Initialize(&elarian.Options{
		ApiKey: "test_api_key",
		OrgId:  "test_org",
		AppId:  "test_app",
	})
	if err != nil {
		log.Fatalf("Error Initializing Elarian: %v \n", err)
	}
	sendMessage(service)
}
