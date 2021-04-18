package test

import (
	"context"
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.GetCustomerState(ctx, elarian.CustomerID(customerID))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		assert.NotNil(t, response.Data)
		assert.True(t, response.Status)
		assert.NotNil(t, response.Data)
	})

	t.Run("It should update a customer's tags", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.UpdateCustomerTag(
			ctx,
			elarian.CustomerID(customerID),
			&elarian.Tag{
				Key:   "TestTag",
				Value: "Test Tag Value",
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.AddCustomerReminder(
			ctx,
			elarian.CustomerID(customerID),
			&elarian.Reminder{
				Key:      "KEY",
				Payload:  "i am a reminder",
				RemindAt: time.Now().Add(time.Second * 3),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.AddCustomerReminderByTag(
			ctx,
			&elarian.Tag{Key: "TestTag", Value: "Test Tag value"},
			&elarian.Reminder{
				Key:      "REMINDER_KEY",
				Payload:  "i am  a reminder",
				RemindAt: time.Now().Add(time.Second * 2),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.AddCustomerReminder(
			ctx,
			elarian.CustomerID(customerID),
			&elarian.Reminder{
				Key:      key,
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
			ctx,
			elarian.CustomerID(customerID),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.AddCustomerReminderByTag(
			ctx,
			&elarian.Tag{Key: "TestTag", Value: "Test Tag value"},
			&elarian.Reminder{
				Key:      key,
				Payload:  "i am a reminder",
				RemindAt: time.Now().Add(time.Second * 2),
			},
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)

		response, err = service.CancelCustomerReminderByTag(
			ctx,
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.UpdateCustomerSecondaryID(
			ctx,
			elarian.CustomerID(customerID),
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
	})

	t.Run("It should delete a customer's Tag", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.DeleteCustomerTag(
			ctx,
			elarian.CustomerID(customerID),
			"TestTag",
		)
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
	})

	t.Run("It should delete a customer's secondary id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.DeleteCustomerSecondaryID(
			ctx,
			elarian.CustomerID(customerID),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.UpdateCustomerAppData(
			ctx,
			elarian.CustomerID(customerID),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.LeaseCustomerAppData(ctx, elarian.CustomerID(customerID))
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.DeleteCustomerAppData(ctx, elarian.CustomerID(customerID))
		if err != nil {
			t.Errorf("Error: %v \n", err)
		}
		assert.NotNil(t, response)
		assert.True(t, response.Status)
		assert.Equal(t, response.CustomerID, customerID)
	})

	t.Run("It should update a customer's metadata", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.UpdateCustomerMetaData(
			ctx,
			elarian.CustomerID(customerID),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.DeleteCustomerMetaData(
			ctx,
			elarian.CustomerID(customerID),
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.UpdateMessagingConsent(
			ctx,
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
