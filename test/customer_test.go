package test

import (
	"log"
	"testing"
	"time"

	elarian "github.com/elarianltd/go-sdk"
	"github.com/stretchr/testify/assert"
)

func Test_Customers(t *testing.T) {
	service, err := elarian.Connect(GetOpts())
	if err != nil {
		log.Fatalln(err)
	}
	defer service.Disconnect()

	t.Run("It Should get customer state", func(t *testing.T) {
		response, err := service.GetCustomerState(&elarian.Customer{ID: customerID})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		assert.NotNil(t, response.Data)
		assert.True(t, response.Status)
		assert.NotNil(t, response.Data)
	})

	t.Run("It should update a customer's tags", func(t *testing.T) {
		response, err := service.UpdateCustomerTag(
			&elarian.Customer{ID: customerID},
			&elarian.Tag{
				Key:        "TestTag",
				Value:      "Test Tag Value",
				Expiration: time.Now().Add(time.Minute * 2),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	t.Run("It should add a customer reminder", func(t *testing.T) {
		response, err := service.AddCustomerReminder(
			&elarian.Customer{ID: customerID},
			&elarian.Reminder{Key: "KEY",
				Payload:  "i am a reminder",
				RemindAt: time.Now().Add(time.Second * 3),
				Interval: time.Duration(time.Second * 60),
			},
		)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.NotEmpty(t, response.CustomerID)
		assert.Equal(t, response.CustomerID, customerID)
		time.Sleep(time.Duration(time.Second * 5))
	})

	t.Run("It should add a customer reminder by tag", func(t *testing.T) {
		response, err := service.AddCustomerReminderByTag(
			&elarian.Tag{Key: "TestTag", Value: "Test Tag value"},
			&elarian.Reminder{Key: "REMINDER_KEY",
				Payload:  "i am  a reminder",
				RemindAt: time.Now().Add(time.Second * 2),
				Interval: time.Duration(time.Hour * 2),
			},
		)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should cancel a customer reminder", func(t *testing.T) {
		key := "REMINDER_KEY"
		response, err := service.AddCustomerReminder(
			&elarian.Customer{ID: customerID},
			&elarian.Reminder{Key: key,
				Payload:  "i am a reminder",
				RemindAt: time.Now().Add(time.Second * 2),
				Interval: time.Duration(time.Hour * 2),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)

		response, err = service.CancelCustomerReminder(
			&elarian.Customer{ID: customerID},
			key,
		)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should cancel a customer reminder by tag", func(t *testing.T) {
		key := "REMINDER_KEY"
		response, err := service.AddCustomerReminderByTag(
			&elarian.Tag{Key: "TestTag", Value: "Test Tag value"},
			&elarian.Reminder{Key: key,
				Payload:  "i am a reminder",
				RemindAt: time.Now().Add(time.Second * 2),
				Interval: time.Duration(time.Hour * 2),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)

		response, err = service.CancelCustomerReminderByTag(
			&elarian.Tag{Key: "TestTag", Value: "Test Tag value"},
			key,
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should update a customer's secondary id", func(t *testing.T) {
		response, err := service.UpdateCustomerSecondaryID(
			&elarian.Customer{ID: customerID},
			&elarian.SecondaryID{
				Key:        "email",
				Value:      "fakeemail@test.com",
				Expiration: time.Now().Add(time.Minute * 2),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should delete a customer's Tag", func(t *testing.T) {
		response, err := service.DeleteCustomerTag(
			&elarian.Customer{ID: customerID},
			"TestTag",
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should delete a customer's secondary id", func(t *testing.T) {
		response, err := service.DeleteCustomerSecondaryID(
			&elarian.Customer{ID: customerID},
			&elarian.SecondaryID{
				Key:   "email",
				Value: "fakeemail@test.com",
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	// update customer appData
	t.Run("It should update a customer's app data", func(t *testing.T) {
		response, err := service.UpdateCustomerAppData(
			&elarian.Customer{ID: customerID},
			&elarian.Appdata{
				Value: elarian.StringDataValue(`{"sessionId": "fake-session-id", "properties": { "ok": 1, "val": "false" } }`),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	// lease customer appData
	t.Run("It should lease a customer's app data", func(t *testing.T) {
		response, err := service.LeaseCustomerAppData(&elarian.Customer{ID: customerID})
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
		assert.Contains(t, response.Appdata.Value, "properties")
	})

	// delete customer appData
	t.Run("It should delete a customer's app data", func(t *testing.T) {
		response, err := service.DeleteCustomerAppData(&elarian.Customer{ID: customerID})
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	t.Run("It should update a customer's metadata", func(t *testing.T) {
		response, err := service.UpdateCustomerMetaData(
			&elarian.Customer{ID: customerID},
			&elarian.Metadata{
				Key:   "DOB",
				Value: elarian.StringDataValue(`{"year": 2020, "day": 13, "month": 10 }`),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	t.Run("It should delete a customer's metadata", func(t *testing.T) {
		response, err := service.DeleteCustomerMetaData(
			&elarian.Customer{ID: customerID},
			"DOB",
		)
		if err != nil {
			t.Errorf("Errror: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	t.Run("It should update message consent", func(t *testing.T) {
		response, err := service.UpdateMessagingConsent(
			&elarian.CustomerNumber{
				Number:   "+254712876967",
				Provider: elarian.CustomerNumberProviderCellular,
			},
			&elarian.MessagingChannelNumber{
				Number:  "21356",
				Channel: elarian.MessagingChannelSms,
			},
			elarian.MessagingConsentUpdateAllow,
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.Equal(t, response.CustomerID, customerID)
	})

	// get customer activity
	// t.Run("It should get a customer's activity", func(t *testing.T) {
	// 	response, err := service.GetCustomerActivity(
	// 		&elarian.CustomerNumber{},
	// 		&elarian.ActivityChannelNumber{},
	// 		"",
	// 	)
	// 	if err != nil {
	// 		t.Errorf("Error: %v \n", err)
	// 	}
	// 	assert.NotNil(t, response)
	// 	assert.True(t, response.Status)
	// 	assert.Equal(t, response.CustomerId.Value, customerID)
	// })

	// adopt customer state

	// t.Run("It should adopt customer state", func(t *testing.T) {
	// 	response, err := service.AdoptCustomerState(
	// 		customerID,
	// 		&elarian.Customer{},
	// 	)
	// 	if err != nil {
	// 		t.Errorf("Error: %v \n", err)
	// 	}
	// 	assert.NotNil(t, response)
	// 	assert.True(t, response.Status)
	// 	assert.Equal(t, response.CustomerId.Value, customerID)
	// })

}
