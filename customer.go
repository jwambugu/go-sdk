package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// NumberProvider is an enum that defines a type of customer number provider. it could be a telco, facebook, telegram or unspecified
	NumberProvider int32

	// CustomerNumber struct
	CustomerNumber struct {
		Number    string         `json:"number,omitempty"`
		Provider  NumberProvider `json:"provider,omitempty"`
		Partition string         `json:"partition,omitempty"`
	}

	// SecondaryID refers to an identifier that can be used on a customer that is unique to a customer and that is provided by you and not the elarian service
	SecondaryID struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// Customer struct defines the parameters required to make any request involving a customer. Note: in every scenario either the Id or the phoneNumber is required but not  both unless otherwise specified
	Customer struct {
		ID             string          `json:"id,omitempty"`
		CustomerNumber *CustomerNumber `json:"customerNumber,omitempty"`
		SecondaryID    *SecondaryID    `json:"secondaryId,omitempty"`
		service        Service
	}

	// CreateCustomer to create a customer
	CreateCustomer struct {
		ID             string          `json:"id,omitempty"`
		CustomerNumber *CustomerNumber `json:"customerNumber,omitempty"`
	}

	// Reminder defines the composition of a reminder. The key is an identifier property. The payload is also a string.
	Reminder struct {
		Interval   time.Duration `json:"interval,omitempty"`
		Key        string        `json:"key,omitempty"`
		Payload    string        `json:"payload,omitempty"`
		Expiration time.Time     `json:"expiration,omitempty"`
	}

	// ReminderNotification struct
	ReminderNotification struct {
		CustomerID string    `json:"customerId,omitempty"`
		WorkID     string    `json:"workId,omitempty"`
		Reminder   *Reminder `json:"reminder,omitempty"`
		Tag        *Tag      `json:"tag,omitempty"`
	}

	// Tag defines a customer tag
	Tag struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}
)

const (
	// CustomerNumberProviderUnspecified type of NumberProvider
	CustomerNumberProviderUnspecified NumberProvider = iota

	// CustomerNumberProviderFacebook type of NumberProvider
	CustomerNumberProviderFacebook

	// CustomerNumberProviderTelco type of NumberProvider represents a telecommunication company
	CustomerNumberProviderTelco

	// CustomerNumberProviderTelegram type of NumberProvider
	CustomerNumberProviderTelegram
)

func (s *service) GetCustomerState(customer *Customer) (*hera.CustomerStateReplyData, error) {
	var request hera.GetCustomerStateRequest
	request.OrgId = s.orgID
	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.GetCustomerStateRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	result, err := s.client.GetCustomerState(ctx, &request)
	return result.GetData(), err
}

func (s *service) AdoptCustomerState(customerID string, otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error) {
	var request hera.AdoptCustomerStateRequest
	request.OrgId = s.orgID
	request.CustomerId = customerID

	if !reflect.ValueOf(otherCustomer.SecondaryID).IsZero() {
		request.OtherCustomer = &hera.
			AdoptCustomerStateRequest_OtherSecondaryId{
			OtherSecondaryId: s.secondaryID(otherCustomer),
		}
	}
	if !reflect.ValueOf(otherCustomer.CustomerNumber).IsZero() {
		request.OtherCustomer = &hera.AdoptCustomerStateRequest_OtherCustomerNumber{
			OtherCustomerNumber: s.customerNumber(otherCustomer),
		}
	}
	if otherCustomer.ID != "" {
		request.OtherCustomer = &hera.
			AdoptCustomerStateRequest_OtherCustomerId{
			OtherCustomerId: otherCustomer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AdoptCustomerState(ctx, &request)
}

func (s *service) AddCustomerReminder(customer *Customer, reminder *Reminder) (*hera.UpdateCustomerStateReply, error) {
	var request hera.AddCustomerReminderRequest
	request.OrgId = s.orgID
	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	request.Reminder = &hera.CustomerReminder{
		Key:        reminder.Key,
		Interval:   durationpb.New(reminder.Interval),
		Payload:    wrapperspb.String(reminder.Payload),
		Expiration: timestamppb.New(reminder.Expiration),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AddCustomerReminder(ctx, &request)
}

func (s *service) AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*hera.TagCommandReply, error) {
	var request hera.AddCustomerReminderTagRequest
	request.OrgId = s.orgID
	request.Tag = &hera.IndexMapping{
		Key:   tag.Key,
		Value: wrapperspb.String(tag.Value),
	}
	request.Reminder = &hera.CustomerReminder{
		Key:        reminder.Key,
		Interval:   durationpb.New(reminder.Interval),
		Payload:    wrapperspb.String(reminder.Payload),
		Expiration: timestamppb.New(reminder.Expiration),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AddCustomerReminderByTag(ctx, &request)
}

func (s *service) CancelCustomerReminder(customer *Customer, key string) (*hera.UpdateCustomerStateReply, error) {
	var request hera.CancelCustomerReminderRequest
	request.AppId = s.appID
	request.OrgId = s.orgID
	request.Key = key
	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminder(ctx, &request)
}

func (s *service) CancelCustomerReminderByTag(tag *Tag, key string) (*hera.TagCommandReply, error) {
	var request hera.CancelCustomerReminderTagRequest
	request.AppId = s.appID
	request.OrgId = s.orgID
	request.Key = key

	request.Tag = &hera.IndexMapping{
		Key:   tag.Key,
		Value: wrapperspb.String(tag.Value),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminderByTag(ctx, &request)
}

func (s *service) UpdateCustomerTag(customer *Customer, tags []Tag) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerTagRequest
	var heraTags []*hera.CustomerIndex
	request.OrgId = s.orgID

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	for _, tag := range tags {
		heraTags = append(heraTags, &hera.CustomerIndex{
			Expiration: timestamppb.New(tag.Expiration),
			Mapping: &hera.IndexMapping{
				Key:   tag.Key,
				Value: wrapperspb.String(tag.Value),
			},
		})
	}
	request.Tags = heraTags
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.UpdateCustomerTag(ctx, &request)
}

func (s *service) DeleteCustomerTag(customer *Customer, keys []string) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerTagRequest
	request.Keys = keys
	request.OrgId = s.orgID

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerTag(ctx, &request)
}

func (s *service) UpdateCustomerSecondaryID(customer *Customer, secondaryIDs []SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	var heraSecIDs []*hera.CustomerIndex
	var request hera.UpdateCustomerSecondaryIdRequest
	request.OrgId = s.orgID

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	for _, secondaryID := range secondaryIDs {
		heraSecIDs = append(heraSecIDs, &hera.CustomerIndex{
			Expiration: timestamppb.New(secondaryID.Expiration),
			Mapping: &hera.IndexMapping{
				Key:   secondaryID.Key,
				Value: wrapperspb.String(secondaryID.Value),
			},
		})
	}
	request.SecondaryIds = heraSecIDs
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.UpdateCustomerSecondaryId(ctx, &request)
}

func (s *service) DeleteCustomerSecondaryID(customer *Customer, secondaryIDs []SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerSecondaryIdRequest
	var heraSecIDs []*hera.IndexMapping

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.
			DeleteCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	for _, secondaryID := range secondaryIDs {
		heraSecIDs = append(heraSecIDs, &hera.IndexMapping{
			Key: secondaryID.Key,
			Value: &wrapperspb.StringValue{
				Value: secondaryID.Value,
			},
		})
	}
	request.OrgId = s.orgID
	request.Mappings = heraSecIDs

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerSecondaryId(ctx, &request)
}

func (s *service) LeaseCustomerMetaData(customer *Customer, key string) (*hera.LeaseCustomerMetadataReply, error) {
	var request hera.LeaseCustomerMetadataRequest
	request.OrgId = s.orgID
	request.Key = key

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.LeaseCustomerMetadata(ctx, &request)
}

func (s *service) UpdateCustomerMetaData(customer *Customer, metadata map[string]string) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerMetadataRequest
	request.OrgId = s.orgID
	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	request.Metadata = map[string]*hera.DataMapValue{}
	for key, value := range metadata {
		request.Metadata[key] = &hera.DataMapValue{
			Value: &hera.DataMapValue_StringVal{
				StringVal: value,
			},
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.UpdateCustomerMetadata(ctx, &request)
}

func (s *service) DeleteCustomerMetaData(customer *Customer, keys []string) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerMetadataRequest
	request.OrgId = s.orgID
	request.Keys = keys

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerMetadata(ctx, &request)
}

func (s *service) NewCustomer(params *CreateCustomer) *Customer {
	var customer Customer
	customer.ID = params.ID
	customer.CustomerNumber = params.CustomerNumber
	customer.service = s
	return &customer
}

// GetState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
func (c *Customer) GetState() (*hera.CustomerStateReplyData, error) {
	return c.service.GetCustomerState(c)
}

// AdoptState copies the state of the second customer to this customer
func (c *Customer) AdoptState(otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AdoptCustomerState(c.ID, otherCustomer)
}

// UpdateTag is used to add more tags to a customer
func (c *Customer) UpdateTag(tags []Tag) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerTag(c, tags)
}

// DeleteTag disaccosiates a tag from a customer
func (c *Customer) DeleteTag(keys []string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerTag(c, keys)
}

// UpdateSecondaryID adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryID(secondaryIds []SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerSecondaryID(c, secondaryIds)
}

// DeleteSecondaryID deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryID(secondaryIds []SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerSecondaryID(c, secondaryIds)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(reminder *Reminder) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AddCustomerReminder(c, reminder)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(key string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.CancelCustomerReminder(c, key)
}

// UpdateMetaData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateMetaData(metadata map[string]string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerMetaData(c, metadata)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(keys []string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerMetaData(c, keys)
}

// LeaseMetaData leases customer metadata
func (c *Customer) LeaseMetaData(key string) (*hera.LeaseCustomerMetadataReply, error) {
	return c.service.LeaseCustomerMetaData(c, key)
}
