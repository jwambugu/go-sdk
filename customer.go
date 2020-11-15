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

	// SecondaryId refers to an identifier that can be used on a customer that is unique to a customer and that is provided by you and not the elarian service
	SecondaryId struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// Customer struct defines the paramters required to make any request involving a customer. Note: in every scenario either the Id or the phoneNumber is required but not  both unless otherwise specified
	Customer struct {
		Id             string          `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber `json:"phoneNumber"`
		SecondaryId    *SecondaryId    `json:"secondaryId"`
		service        Service
	}

	// CreateCustomer to create a customer
	CreateCustomer struct {
		Id             string          `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber `json:"phoneNumber,omitempty"`
	}

	// Reminder defines the composition of a reminder. The key is an identifier property. The payload is also a string.
	Reminder struct {
		Expiration time.Time `json:"expiration,omitempty"`
		Interval   int64     `json:"interval,omitempty"`
		Key        string    `json:"key,omitempty"`
		Payload    string    `json:"payload,omitempty"`
	}

	// ReminderNotification struct
	ReminderNotification struct {
		CustomerId string    `json:"customerId,omitempty"`
		WorkId     string    `json:"workId,omitempty"`
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
	CUSTOMER_NUMBER_PROVIDER_UNSPECIFIED NumberProvider = iota
	CUSTOMER_NUMBER_PROVIDER_FACEBOOK
	CUSTOMER_NUMBER_PROVIDER_TELCO
	CUSTOMER_NUMBER_PROVIDER_TELEGRAM
)

func (s *service) GetCustomerState(
	customer *Customer,
) (*hera.CustomerStateReplyData, error) {
	var request hera.GetCustomerStateRequest
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.GetCustomerStateRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	result, err := s.client.GetCustomerState(ctx, &request)
	return result.GetData(), err
}

func (s *service) AdoptCustomerState(
	customer *Customer,
	otherCustomer *Customer,
) (*hera.UpdateCustomerStateReply, error) {

	var request hera.AdoptCustomerStateRequest
	request.OrgId = s.orgId
	request.CustomerId = customer.Id

	request.OtherCustomer = &hera.AdoptCustomerStateRequest_OtherCustomerId{
		OtherCustomerId: otherCustomer.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AdoptCustomerState(ctx, &request)
}

func (s *service) AddCustomerReminder(
	customer *Customer,
	reminder *Reminder,
) (*hera.UpdateCustomerStateReply, error) {

	var request hera.AddCustomerReminderRequest
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(reminder.Expiration),
		Interval: &durationpb.Duration{
			Seconds: int64(time.Duration(reminder.Interval) * time.Second),
		},
		Key: reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: reminder.Payload,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AddCustomerReminder(ctx, &request)
}

func (s *service) AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*hera.TagCommandReply, error) {
	var request hera.AddCustomerReminderTagRequest
	request.OrgId = s.orgId

	request.Tag = &hera.IndexMapping{
		Key: tag.Key,
		Value: &wrapperspb.StringValue{
			Value: tag.Value,
		},
	}

	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(reminder.Expiration),
		Interval: &durationpb.Duration{
			Seconds: int64(time.Duration(reminder.Interval) * time.Second),
		},
		Key: reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: reminder.Payload,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AddCustomerReminderByTag(ctx, &request)
}

func (s *service) CancelCustomerReminder(customer *Customer, key string) (*hera.UpdateCustomerStateReply, error) {
	var request hera.CancelCustomerReminderRequest
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.Key = key

	if customer.Id != "" {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminder(ctx, &request)
}

func (s *service) CancelCustomerReminderByTag(tag *Tag, key string) (*hera.TagCommandReply, error) {
	var request hera.CancelCustomerReminderTagRequest
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.Key = key

	request.Tag = &hera.IndexMapping{
		Key: tag.Key,
		Value: &wrapperspb.StringValue{
			Value: tag.Value,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminderByTag(ctx, &request)
}

func (s *service) UpdateCustomerTag(customer *Customer, tags []Tag) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerTagRequest
	var heraTags []*hera.CustomerIndex
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	for _, tag := range tags {
		heraTags = append(heraTags, &hera.CustomerIndex{
			Expiration: timestamppb.New(tag.Expiration),
			Mapping: &hera.IndexMapping{
				Key: tag.Key,
				Value: &wrapperspb.StringValue{
					Value: tag.Value,
				},
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
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerTag(ctx, &request)
}

func (s *service) UpdateCustomerSecondaryId(customer *Customer,
	secondaryIds []SecondaryId,
) (*hera.UpdateCustomerStateReply, error) {
	var heraSecIds []*hera.CustomerIndex
	var request hera.UpdateCustomerSecondaryIdRequest
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	for _, secondaryId := range secondaryIds {
		heraSecIds = append(heraSecIds, &hera.CustomerIndex{
			Expiration: timestamppb.New(secondaryId.Expiration),
			Mapping: &hera.IndexMapping{
				Key: secondaryId.Key,
				Value: &wrapperspb.StringValue{
					Value: secondaryId.Value,
				},
			},
		})
	}
	request.SecondaryIds = heraSecIds

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.UpdateCustomerSecondaryId(ctx, &request)
}

func (s *service) DeleteCustomerSecondaryId(
	customer *Customer,
	secondaryIds []SecondaryId,
) (*hera.UpdateCustomerStateReply, error) {

	var request hera.DeleteCustomerSecondaryIdRequest
	var heraSecIds []*hera.IndexMapping

	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.
			DeleteCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	for _, secondaryId := range secondaryIds {
		heraSecIds = append(heraSecIds, &hera.IndexMapping{
			Key: secondaryId.Key,
			Value: &wrapperspb.StringValue{
				Value: secondaryId.Value,
			},
		})
	}
	request.OrgId = s.orgId
	request.Mappings = heraSecIds

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerSecondaryId(ctx, &request)
}

func (s *service) LeaseCustomerMetaData(
	customer *Customer,
	key string,
) (*hera.LeaseCustomerMetadataReply, error) {

	var request hera.LeaseCustomerMetadataRequest
	request.OrgId = s.orgId
	request.Key = key

	if customer.Id != "" {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.LeaseCustomerMetadata(ctx, &request)
}

func (s *service) UpdateCustomerMetaData(
	customer *Customer,
	metadata map[string]string,
) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerMetadataRequest
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
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

func (s *service) DeleteCustomerMetaData(
	customer *Customer,
	keys []string,
) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerMetadataRequest
	request.OrgId = s.orgId
	request.Keys = keys

	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.DeleteCustomerMetadata(ctx, &request)
}

func (s *service) NewCustomer(params *CreateCustomer) *Customer {
	var customer Customer
	customer.Id = params.Id
	customer.CustomerNumber = params.CustomerNumber
	customer.service = s

	return &customer
}

// GetState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
func (c *Customer) GetState() (*hera.CustomerStateReplyData, error) {
	return c.service.GetCustomerState(c)
}

// AdoptState copies the state of the second customer to this customer
func (c *Customer) AdoptState(
	otherCustomer *Customer,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AdoptCustomerState(c, otherCustomer)
}

// UpdateTag is used to add more tags to a customer
func (c *Customer) UpdateTag(
	tags []Tag,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerTag(c, tags)
}

// DeleteTag disaccosiates a tag from a customer
func (c *Customer) DeleteTag(
	keys []string,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerTag(c, keys)
}

// UpdateSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryId(
	secondaryIds []SecondaryId,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerSecondaryId(c, secondaryIds)
}

// DeleteSecondaryId deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryId(
	secondaryIds []SecondaryId,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerSecondaryId(c, secondaryIds)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(
	reminder *Reminder,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AddCustomerReminder(c, reminder)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(
	key string,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.CancelCustomerReminder(c, key)
}

// UpdateMetaData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateMetaData(
	metadata map[string]string,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerMetaData(c, metadata)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(
	keys []string,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerMetaData(c, keys)
}

// LeaseMetaData leases customer metadata
func (c *Customer) LeaseMetaData(
	key string,
) (*hera.LeaseCustomerMetadataReply, error) {
	return c.service.LeaseCustomerMetaData(c, key)
}
