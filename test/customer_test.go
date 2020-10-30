package test

import (
	"reflect"
	"testing"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Customer(t *testing.T) {
	opts := &elarian.Options{
		ApiKey: APIKey,
		AppId:  AppId,
		OrgId:  OrgId,
	}
	service, err := elarian.Initialize(opts)
	if err != nil {
		t.Fatal(err)
	}
	var customer *elarian.Customer
	customer.Id = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	customer.CustomerNumber = elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
	}

	t.Run("It Should get customer state", func(t *testing.T) {
		response, err := service.GetCustomerState(customer)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.MessagingState).IsZero() {
			t.Errorf("Expected customer messaging state by didn't get any")
		}
		t.Logf("customer state %v", response)
	})

	t.Run("It should adopt customer state", func(t *testing.T) {
		var otherCustomer elarian.Customer
		otherCustomer.CustomerNumber = elarian.CustomerNumber{
			Number:   "+254711276275",
			Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
		}

		response, err := service.AdoptCustomerState(
			customer,
			&otherCustomer,
		)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		if response.Status != true {
			t.Errorf("Expected status to be: %v but got: %v", true, false)
		}
		t.Logf("Response %v", response)
	})

	t.Run("It should add a customer reminder", func(t *testing.T) {
		reminder := &elarian.Reminder{
			Key:        "reminder_key",
			Expiration: time.Now().Add(time.Minute + 1),
			Interval:   int64(2 * time.Second),
			Payload:    "i am a reminder",
		}
		response, err := service.AddCustomerReminder(customer, reminder)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("response %v", response)
	})

	t.Run("It should add a customer reminder by tag", func(t *testing.T) {
		tag := &elarian.Tag{
			Key:   "some key",
			Value: "some value",
		}
		reminder := &elarian.Reminder{
			Key:        "reminder_key",
			Expiration: time.Now().Add(time.Minute + 1),
			Interval:   int64(2 * time.Second),
			Payload:    "i am a reminder",
		}

		response, err := service.AddCustomerReminderByTag(reminder, tag)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.WorkId).IsZero() {
			t.Errorf("Expecred a WorkId but didn't get any")
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("response %v", response)
	})

	t.Run("It should cancel a customer reminder", func(t *testing.T) {
		reminder := &elarian.Reminder{
			Key:        "reminder_key",
			Expiration: time.Now().Add(time.Minute + 1),
			Payload:    "i am a reminder",
		}
		response, err := service.AddCustomerReminder(customer, reminder)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("response %v", response)
		res, err := service.CancelCustomerReminder(customer, "reminder_key")
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should cancel a customer reminder by tag", func(t *testing.T) {

		tag := &elarian.Tag{
			Key:   "tag_key",
			Value: "i am a value",
		}
		reminder := &elarian.Reminder{
			Key:        "reminder_key",
			Expiration: time.Now().Add(time.Minute + 1),
			Payload:    "i am a reminder",
		}
		response, err := service.AddCustomerReminderByTag(reminder, tag)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.WorkId).IsZero() {
			t.Errorf("Expecred a WorkId but didn't get any")
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("response %v", response)

		res, err := service.CancelCustomerReminderByTag("reminder_key", tag)
		if err != nil {
			t.Errorf("Erorr %v", err)
		}
		if response.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("response %v", res)
	})

	t.Run("It should update a customer's tags", func(t *testing.T) {
		tags := []elarian.Tag{
			{
				Key:        "new_tag",
				Value:      "new_tag_value",
				Expiration: time.Now().Add(time.Duration(1 * time.Minute)),
			},
		}
		res, err := service.UpdateCustomerTag(customer, tags)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should delete a customer's tags", func(t *testing.T) {

		keys := []string{"new_tag"}
		res, err := service.DeleteCustomerTag(customer, keys)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Errorf("Expected a description but didn't get any")
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should update a customer's secondary ids", func(t *testing.T) {
		secondaryIds := []elarian.CustomerSecondaryId{
			{
				Key:        "my_app_customer_Id",
				Value:      "123456wq",
				Expiration: time.Now().Add(time.Duration(1 * time.Minute)),
			},
		}
		res, err := service.UpdateCustomerSecondaryId(customer, secondaryIds)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Error("Expected a description but didn't get any")
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should delete a customer's secondary ids", func(t *testing.T) {
		secondaryIds := []elarian.CustomerSecondaryId{
			{
				Key:   "my_app_customer_Id",
				Value: "123456wq",
			},
		}
		res, err := service.DeleteCustomerSecondaryId(customer, secondaryIds)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Error("Expected a description but didn't get any")
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should update a customers metadata", func(t *testing.T) {

		metadata := map[string]string{
			"some_key":       "some_value",
			"some_other_key": "some_other_value",
		}

		res, err := service.UpdateCustomerMetaData(customer, metadata)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Error("Expected a description but didn't get any")
		}
		if res.Status != true {
			t.Errorf("Expected status to be: %v but got: %v", true, false)
		}
		t.Logf("Response %v", res)
	})

	t.Run("It should delete a customer's metadata", func(t *testing.T) {

		metadata := []string{
			"some_key", "some_other_key",
		}
		res, err := service.DeleteCustomerMetaData(customer, metadata)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if res.Description == "" {
			t.Error("Expected a description but didn't get any")
		}
		if res.Status != true {
			t.Errorf("Expected status to be: %v but got: %v", true, false)
		}
		t.Logf("Response %v", res)
	})
}
