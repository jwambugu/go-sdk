package test

import (
	"sync"
	"testing"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_StreamNotifications(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
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
		Expiration: time.Now().Add(time.Duration(20 * time.Second)),
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

	streamChannel, _ := service.StreamNotifications(AppID)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		res := <-streamChannel
		reminder := res.GetReminder()
		if reminder.ProductId != ProductID {
			t.Errorf("Expected a product id %s but got %s", ProductID, reminder.ProductId)
		}
		if reminder.Reminder.Key != "reminder_key" {
			t.Errorf("Expected the reminder key to be: %s but got %s", "reminder_key", reminder.Reminder.Key)
		}
		t.Logf("response %v", res.GetReminder())
	}()
	wg.Wait()
}
