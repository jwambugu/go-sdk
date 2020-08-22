// Package main implements a the simple examples that demonstrate hot to use the elarian SDK.
package main

import (
	"log"
	"time"

	elarian "github.com/elarianltd/go-sdk"
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func sendSms(client elarian.Elarian) {
	customer := &elarian.Customer{
		PhoneNumber: elarian.PhoneNumber{
			Number:   "",
			Provider: hera.CustomerNumberProvider_CUSTOMER_NUMBER_PROVIDER_FACEBOOK,
		},
	}

	message := &elarian.SendMessageRequest{
		AppID:     "",
		ProductID: "",
		ChannelNumber: elarian.MessagingChannelNumber{
			Number:  "",
			Channel: hera.MessagingChannel_MESSAGING_CHANNEL_FB_MESSENGER,
		},
		Body: elarian.MessageBody{
			Text: "",
		},
	}

	res, err := client.SendMessage(customer, message)
	if err != nil {
		log.Fatalf("Could not send sms %v", err)
	}
	log.Printf("Customer id %v", res.GetCustomerId())
}

func getCustomerState(client elarian.Elarian) {
	var customer elarian.Customer
	customer.ID = "customer_id"
	var request elarian.CustomerStateRequest
	request.AppID = "app_id"

	res, err := client.GetCustomerState(&customer, &request)
	if err != nil {
		log.Fatalf("could not get customer state %v", err)
	}
	log.Printf("customer state %v", res)
}

func addCustomerReminder(client elarian.Elarian) {
	var customer elarian.Customer
	customer.ID = "customer_id"

	var request elarian.CustomerReminderRequest
	request.AppID = "app_id"
	request.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		Payload:    "I am a reminder",
		ProductID:  "product_id",
	}

	res, err := client.AddCustomerReminder(&customer, &request)
	if err != nil {
		log.Fatalf("could not set a reminder %v", err)
	}
	log.Printf("response %v", res)
}

func streamCustomerNotifications(service elarian.Elarian) {
	streamChan, errorChan := service.StreamNotifications("app_id")
	res := <-streamChan
	err := <-errorChan

	if err != nil {
		log.Fatalf("notification stream error %v", err)
	}
	log.Printf("response %v", res.GetReminder())
}

func main() {
	service, err := elarian.Initialize("api_key", true)
	if err != nil {
		log.Fatal(err)
	}
	sendSms(service)
	getCustomerState(service)
	addCustomerReminder(service)
}
