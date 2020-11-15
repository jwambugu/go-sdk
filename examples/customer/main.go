package main

import (
	"fmt"
	"log"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func getCustomerState(service elarian.Service) {
	var cust elarian.Customer
	cust.Id = ""
	res, err := service.GetCustomerState(&cust)
	if err != nil {
		log.Fatalf("Error getting customer state %v \n", err)
	}
	log.Printf("State %v \n", res)
}

func getState(service elarian.Service) {
	cust := service.NewCustomer(&elarian.CreateCustomer{Id: ""})
	res, err := cust.GetState()
	if err != nil {
		log.Fatalf("Error getting customer state %v \n", err)
	}
	log.Printf("State %v \n", res)
}

func addReminder(service elarian.Service) {
	err := service.AddNotificationSubscriber(
		elarian.ELARIAN_REMINDER_NOTIFICATION,
		func(svc elarian.Service, cust *elarian.Customer, data interface{}) {
			notf, ok := data.(elarian.ReminderNotification)
			if !ok {
				log.Fatalf("Corrupted notification data")
			}
			fmt.Printf("Reminder %v \n", notf)
		},
	)
	if err != nil {
		log.Fatalf("Error adding a subscriber %v \n", err)
	}

	var reminder elarian.Reminder
	reminder.Key = "reminder_key"
	reminder.Payload = "i am a reminder"
	reminder.Expiration = time.Now().Add(time.Minute + 1)

	cust := service.NewCustomer(&elarian.CreateCustomer{Id: ""})
	res, err := cust.AddReminder(&reminder)
	if err != nil {
		log.Fatalf("could not set a reminder %v", err)
	}
	log.Printf("response %v", res)
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
	getCustomerState(service)
	getState(service)
	addReminder(service)
}
