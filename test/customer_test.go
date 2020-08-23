package test

import (
	"log"
	"reflect"
	"testing"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_GetCustomerState(t *testing.T) {
	service, err := elarian.Initialize(APIKey, true)
	if err != nil {
		log.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.CustomerStateRequest
	request.AppID = AppID

	response, err := service.GetCustomerState(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if reflect.ValueOf(response.MessagingState).IsZero() {
		t.Errorf("Expected customer messaging state by didn't get any")
	}
	t.Logf("customer state %v", response.MessagingState)
}

func Test_AddCustomerReminder(t *testing.T) {
	service, err := elarian.Initialize(APIKey, true)
	if err != nil {
		log.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.CustomerReminderRequest
	request.AppID = AppID
	request.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		ProductID:  ProductID,
		Payload:    "i am a reminder",
	}

	response, err := service.AddCustomerReminder(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if response.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("response %v", response)
}
