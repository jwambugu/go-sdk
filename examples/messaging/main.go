package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func sendMessage(service elarian.Service) {
	var (
		custNumber  *elarian.CustomerNumber
		channel     *elarian.MessagingChannelNumber
		messageBody *elarian.MessageBody
	)

	custNumber = &elarian.CustomerNumber{
		Number:   "+254708752502",
		Provider: elarian.CustomerNumberProviderCellular,
	}
	customer := service.NewCustomer(&elarian.CreateCustomer{
		CustomerNumber: custNumber,
	})

	channel = &elarian.MessagingChannelNumber{
		Number:  "",
		Channel: elarian.MessagingChannelSms,
	}
	messageBody = &elarian.MessageBody{
		Text: "Hello world from the go sdk",
	}

	response, err := service.SendMessage(customer, channel, messageBody)
	if err != nil {
		log.Fatalf("Message not send %v \n", err.Error())
	}
	log.Printf("Status %d Description %s \n customerID %s \n", response.Status, response.Description, response.CustomerId)
}

func main() {
	const (
		AppID  string = "zordTest"
		OrgID  string = "og-hv3yFs"
		APIKey string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
	)
	service, err := elarian.Connect(&elarian.Options{
		OrgID:  OrgID,
		AppID:  AppID,
		APIKey: APIKey,
	}, nil)
	if err != nil {
		log.Fatalf("Error Initializing Elarian: %v \n", err)
	}
	sendMessage(service)
}
