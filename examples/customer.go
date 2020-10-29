package main

import (
	"log"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func getCustomerState(client elarian.Service) {
	var customer elarian.Customer
	customer.Id = "customer_id"

	response, err := client.GetCustomerState(&customer)

	calvin, err := client.NewCustomer(&elarian.CreateCustomerParams{Id: ""})
	res, err := calvin.GetState()
	log.Println(res)

	if err != nil {
		log.Fatalf("could not get customer state %v", err)
	}
	log.Printf("customer state %v", response)
}

func addCustomerReminder(client elarian.Service) {
	var customer elarian.Customer
	customer.Id = "customer_id"

	var request elarian.CustomerReminderRequest
	request.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		Payload:    "I am a reminder",
	}

	response, err := client.AddCustomerReminder(&customer, &request)
	if err != nil {
		log.Fatalf("could not set a reminder %v", err)
	}
	log.Printf("response %v", response)
}

func adoptCustomerState(client elarian.Service) {
	var customer elarian.Customer
	customer.Id = "customer_id"

	var otherCustomer elarian.Customer
	otherCustomer.Id = "otherCustomer_id"

	response, err := client.AdoptCustomerState(
		&customer,
		&otherCustomer,
	)

	if err != nil {
		log.Fatalf("could not adopt customeer state %v", err)
	}
	log.Printf("response %v", response)
}
