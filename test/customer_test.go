package test

import (
	"reflect"
	"testing"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_GetCustomerState(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
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
	t.Logf("customer state %v", response)
}

func Test_AdoptCustomerState(t *testing.T) {
	var customer elarian.Customer
	var otherCustomer elarian.Customer
	var request elarian.AdoptCustomerStateRequest

	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	customer.ID = "el_cst_c617c20cec9b52bf7698ea58695fb8bc"
	otherCustomer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254711276275",
		Provider: elarian.CustomerNumberProviderTelco,
	}
	request.AppID = AppID

	response, err := service.AdoptCustomerState(&customer, &otherCustomer, &request)
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
}

func Test_AddCustomerReminder(t *testing.T) {
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
func Test_AddCustomerReminderByTag(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}

	var request elarian.CustomerReminderByTagRequest
	request.AppID = AppID
	request.Tag = elarian.Tag{
		Key:   "some key",
		Value: "some value",
	}
	request.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		ProductID:  ProductID,
		Payload:    "i am a reminder",
	}

	response, err := service.AddCustomerReminderByTag(&request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if reflect.ValueOf(response.WorkId).IsZero() {
		t.Errorf("Expecred a WorkID but didn't get any")
	}
	if response.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("response %v", response)
}

func Test_CancelCustomerReminder(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var reminderRequest elarian.CustomerReminderRequest
	reminderRequest.AppID = AppID
	reminderRequest.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		ProductID:  ProductID,
		Payload:    "i am a reminder",
	}
	response, err := service.AddCustomerReminder(&customer, &reminderRequest)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if response.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("response %v", response)

	var cancelReminderRequest elarian.CancelCustomerReminderRequest
	cancelReminderRequest.AppID = AppID
	cancelReminderRequest.Key = "reminder_key"
	cancelReminderRequest.ProductID = ProductID

	res, err := service.CancelCustomerReminder(&customer, &cancelReminderRequest)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if res.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("Response %v", res)
}

func Test_CancelCustomerReminderByTag(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}

	var reminderRequest elarian.CustomerReminderByTagRequest
	reminderRequest.AppID = AppID
	reminderRequest.Tag = elarian.Tag{
		Key:   "tag_key",
		Value: "i am a value",
	}
	reminderRequest.Reminder = elarian.Reminder{
		Key:        "reminder_key",
		Expiration: time.Now().Add(time.Minute + 1),
		ProductID:  ProductID,
		Payload:    "i am a reminder",
	}

	response, err := service.AddCustomerReminderByTag(&reminderRequest)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if reflect.ValueOf(response.WorkId).IsZero() {
		t.Errorf("Expecred a WorkID but didn't get any")
	}
	if response.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("response %v", response)

	var cancelReminderRequest elarian.CancelCustomerReminderByTagRequest
	cancelReminderRequest.AppID = AppID
	cancelReminderRequest.Key = "reminder_key"
	cancelReminderRequest.ProductID = ProductID
	cancelReminderRequest.Tag = elarian.Tag{
		Key:   "tag_key",
		Value: "i am a value",
	}
	res, err := service.CancelCustomerReminderByTag(&cancelReminderRequest)
	if err != nil {
		t.Errorf("Erorr %v", err)
	}
	if response.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("response %v", res)
}
func Test_UpdateCustomerTag(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}
	var request elarian.UpdateCustomerTagRequest
	request.AppID = AppID
	request.Tags = []elarian.Tag{
		{
			Key:        "new_tag",
			Value:      "new_tag_value",
			Expiration: time.Now().Add(time.Duration(1 * time.Minute)),
		},
	}
	res, err := service.UpdateCustomerTag(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if res.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("Response %v", res)
}
func Test_DeleteCustomerTag(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}
	var request elarian.DeleteCustomerTagRequest
	request.AppID = AppID
	request.Tags = []string{"new_tag"}
	res, err := service.DeleteCustomerTag(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if res.Description == "" {
		t.Errorf("Expected a description but didn't get any")
	}
	t.Logf("Response %v", res)
}
func Test_UpdateCustomerSecondaryID(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}
	var request elarian.UpdateCustomerSecondaryIDRequest
	request.AppID = AppID
	request.SecondaryIDs = []elarian.CustomerSecondaryID{
		{
			Key:        "my_app_customer_Id",
			Value:      "123456wq",
			Expiration: time.Now().Add(time.Duration(1 * time.Minute)),
		},
	}
	res, err := service.UpdateCustomerSecondaryID(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if res.Description == "" {
		t.Error("Expected a description but didn't get any")
	}
	t.Logf("Response %v", res)
}
func Test_DeleteCustomerSecondaryID(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}
	var request elarian.DeleteCustomerSecondaryIDRequest
	request.AppID = AppID
	request.SecondaryIDs = []elarian.CustomerSecondaryID{
		{
			Key:   "my_app_customer_Id",
			Value: "123456wq",
		},
	}
	res, err := service.DeleteCustomerSecondaryID(&customer, &request)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if res.Description == "" {
		t.Error("Expected a description but didn't get any")
	}
	t.Logf("Response %v", res)
}
func Test_UpdateCustomerMetaData(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.UpdateCustomerMetadataRequest
	request.AppID = AppID
	request.Metadata = map[string]string{
		"some_key":       "some_value",
		"some_other_key": "some_other_value",
	}

	res, err := service.UpdateCustomerMetaData(&customer, &request)
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
}
func Test_DeleteCustomerMetaData(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.PhoneNumber = elarian.PhoneNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	var request elarian.DeleteCustomerMetadataRequest
	request.AppID = AppID
	request.Metadata = []string{
		"some_key", "some_other_key",
	}
	res, err := service.DeleteCustomerMetaData(&customer, &request)
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
}