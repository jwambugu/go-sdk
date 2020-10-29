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

	// CustomerSecondaryId refers to an identifier that can be used on a customer that is unique to a customer and that is provided by you and not the elarian service
	CustomerSecondaryId struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// Customer struct defines the paramters required to make any request involving a customer. Note: in every scenario either the Id or the phoneNumber is required but not  both unless otherwise specified
	Customer struct {
		Id             string              `json:"customerId,omitempty"`
		CustomerNumber CustomerNumber      `json:"phoneNumber"`
		SecondaryId    CustomerSecondaryId `json:"secondaryId"`
		service        Service
	}

	// CreateCustomerParams to create a customer
	CreateCustomerParams struct {
		Id             string         `json:"customerId,omitempty"`
		CustomerNumber CustomerNumber `json:"phoneNumber,omitempty"`
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
		CustomerId string   `json:"customerId,omitempty"`
		WorkId     string   `json:"workId,omitempty"`
		Reminder   Reminder `json:"reminder,omitempty"`
		Tag        Tag      `json:"tag,omitempty"`
	}

	// Tag defines a customer tag
	Tag struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// CustomerReminderRequest defines the arguments required to create a reminder based on a customer.
	CustomerReminderRequest struct {
		Reminder Reminder `json:"reminder,omitempty"`
	}

	// CancelCustomerReminderRequest defines the arguments required to cancel a reminder set on a customer
	CancelCustomerReminderRequest struct {
		AppId    string `json:"appId,omitempty"`
		Key      string `json:"key,omitempty"`
		Reminder string `json:"reminder,omitempty"`
	}

	// CustomerReminderByTagRequest defines the arguments required to Create a reminder based on a tag. With this the reminder can applly to a group of customers.
	CustomerReminderByTagRequest struct {
		Reminder Reminder `json:"reminder,omitempty"`
		Tag      Tag      `json:"tag,omitempty"`
	}

	// CancelCustomerReminderByTagRequest defines the arguments required to cancel a reminder set on a tag
	CancelCustomerReminderByTagRequest struct {
		AppId string `json:"appId,omitempty"`
		Key   string `json:"key,omitempty"`
		Tag   Tag    `json:"tag,omitempty"`
	}

	// DeleteCustomerTagRequest defines the arguments required to delete a customer's tags
	DeleteCustomerTagRequest struct {
		Keys []string `json:"keys,omitempty"`
	}

	// UpdateCustomerTagRequest defines the arguments required to update a customer's tags, you can add one or more tags
	UpdateCustomerTagRequest struct {
		OrgId string `json:"orgId,omitempty"`
		Tags  []Tag  `json:"tags,omitempty"`
	}

	// UpdateCustomerSecondaryIdRequest defines the arguments required to update a customer's secondary Ids you can add one or more secondary Ids
	UpdateCustomerSecondaryIdRequest struct {
		SecondaryIds []CustomerSecondaryId `json:"secondaryIds,omitempty"`
	}

	// DeleteCustomerSecondaryIdRequest defines the arguments required to delete a customer's secondary Identifiers. You can provide one or more secondary Ids you want to delete
	DeleteCustomerSecondaryIdRequest struct {
		SecondaryIds []CustomerSecondaryId `json:"secondaryIds,omitempty"`
	}

	// UpdateCustomerMetadataRequest defines the arguments required to update a customer's metadata
	UpdateCustomerMetadataRequest struct {
		Metadata map[string]string `json:"metadata,omitempty"`
	}

	// DeleteCustomerMetadataRequest defines the arguments required to delete a customer's metadata
	DeleteCustomerMetadataRequest struct {
		Metadata []string `json:"metadata,omitempty"`
	}

	// LeaseCustomerMetadataRequest defines the arguments required to lease metadata
	LeaseCustomerMetadataRequest struct {
		Key string `json:"key,omitempty"`
	}
)

const (
	CUSTOMER_NUMBER_PROVIDER_UNSPECIFIED NumberProvider = iota
	CUSTOMER_NUMBER_PROVIDER_FACEBOOK
	CUSTOMER_NUMBER_PROVIDER_TELCO
	CUSTOMER_NUMBER_PROVIDER_TELEGRAM
)

func (s *service) setCustomerNumber(customer *Customer) *hera.CustomerNumber {
	return &hera.CustomerNumber{
		Number:   customer.CustomerNumber.Number,
		Provider: hera.CustomerNumberProvider(customer.CustomerNumber.Provider),
		Partition: &wrapperspb.StringValue{
			Value: customer.CustomerNumber.Partition,
		},
	}
}

func (s *service) setCustomerNumbers(customerNumbers []*CustomerNumber) []*hera.CustomerNumber {
	var numbers []*hera.CustomerNumber

	for _, number := range customerNumbers {
		numbers = append(numbers, &hera.CustomerNumber{
			Number:    number.Number,
			Provider:  hera.CustomerNumberProvider(number.Provider),
			Partition: wrapperspb.String(number.Partition),
		})
	}
	return numbers
}

func (s *service) setCustomerSecondaryId(
	customer *Customer,
) *hera.IndexMapping {
	return &hera.IndexMapping{
		Key: customer.SecondaryId.Key,
		Value: &wrapperspb.StringValue{
			Value: customer.SecondaryId.Value,
		},
	}
}

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
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AdoptCustomerState(ctx, &request)
}

func (s *service) AddCustomerReminder(
	customer *Customer,
	params *CustomerReminderRequest,
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
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(params.Reminder.Expiration),
		Interval: &durationpb.Duration{
			Seconds: int64(time.Duration(params.Reminder.Interval) * time.Second),
		},
		Key: params.Reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: params.Reminder.Payload,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AddCustomerReminder(ctx, &request)
}

func (s *service) AddCustomerReminderByTag(
	params *CustomerReminderByTagRequest,
) (*hera.TagCommandReply, error) {
	var request hera.AddCustomerReminderTagRequest
	request.OrgId = s.orgId

	request.Tag = &hera.IndexMapping{
		Key: params.Tag.Key,
		Value: &wrapperspb.StringValue{
			Value: params.Tag.Value,
		},
	}

	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(params.Reminder.Expiration),
		Interval: &durationpb.Duration{
			Seconds: int64(time.Duration(params.Reminder.Interval) * time.Second),
		},
		Key: params.Reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: params.Reminder.Payload,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AddCustomerReminderByTag(ctx, &request)
}

func (s *service) CancelCustomerReminder(
	customer *Customer,
	params *CancelCustomerReminderRequest,
) (*hera.UpdateCustomerStateReply, error) {
	var request hera.CancelCustomerReminderRequest
	request.AppId = params.AppId
	request.OrgId = s.orgId
	request.Key = params.Key

	if customer.Id != "" {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminder(ctx, &request)
}

func (s *service) CancelCustomerReminderByTag(
	params *CancelCustomerReminderByTagRequest,
) (*hera.TagCommandReply, error) {
	var request hera.CancelCustomerReminderTagRequest
	request.AppId = params.AppId
	request.OrgId = s.orgId
	request.Key = params.Key

	request.Tag = &hera.IndexMapping{
		Key: params.Tag.Key,
		Value: &wrapperspb.StringValue{
			Value: params.Tag.Value,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminderByTag(ctx, &request)
}

func (s *service) UpdateCustomerTag(
	customer *Customer,
	params *UpdateCustomerTagRequest,
) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerTagRequest
	var tags []*hera.CustomerIndex
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	for _, tag := range params.Tags {
		tags = append(tags, &hera.CustomerIndex{
			Expiration: timestamppb.New(tag.Expiration),
			Mapping: &hera.IndexMapping{
				Key: tag.Key,
				Value: &wrapperspb.StringValue{
					Value: tag.Value,
				},
			},
		})
	}
	request.Tags = tags

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.UpdateCustomerTag(ctx, &request)
}

func (s *service) DeleteCustomerTag(
	customer *Customer,
	params *DeleteCustomerTagRequest,
) (*hera.UpdateCustomerStateReply, error) {

	var request hera.DeleteCustomerTagRequest
	request.Keys = params.Keys
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerTag(ctx, &request)
}

func (s *service) UpdateCustomerSecondaryId(
	customer *Customer,
	params *UpdateCustomerSecondaryIdRequest,
) (*hera.UpdateCustomerStateReply, error) {

	var secondaryIds []*hera.CustomerIndex
	var request hera.UpdateCustomerSecondaryIdRequest
	request.OrgId = s.orgId

	if customer.Id != "" {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	for _, secondaryId := range params.SecondaryIds {
		secondaryIds = append(secondaryIds, &hera.CustomerIndex{
			Expiration: timestamppb.New(secondaryId.Expiration),
			Mapping: &hera.IndexMapping{
				Key: secondaryId.Key,
				Value: &wrapperspb.StringValue{
					Value: secondaryId.Value,
				},
			},
		})
	}
	request.SecondaryIds = secondaryIds

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.UpdateCustomerSecondaryId(ctx, &request)
}

func (s *service) DeleteCustomerSecondaryId(
	customer *Customer,
	params *DeleteCustomerSecondaryIdRequest,
) (*hera.UpdateCustomerStateReply, error) {

	var request hera.DeleteCustomerSecondaryIdRequest
	var secondaryIds []*hera.IndexMapping

	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.
			DeleteCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	for _, secondaryId := range params.SecondaryIds {
		secondaryIds = append(secondaryIds, &hera.IndexMapping{
			Key: secondaryId.Key,
			Value: &wrapperspb.StringValue{
				Value: secondaryId.Value,
			},
		})
	}
	request.OrgId = s.orgId
	request.Mappings = secondaryIds

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerSecondaryId(ctx, &request)
}

func (s *service) LeaseCustomerMetaData(
	customer *Customer,
	params *LeaseCustomerMetadataRequest,
) (*hera.LeaseCustomerMetadataReply, error) {

	var request hera.LeaseCustomerMetadataRequest
	request.OrgId = s.orgId
	request.Key = params.Key

	if customer.Id != "" {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.LeaseCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.LeaseCustomerMetadata(ctx, &request)
}

func (s *service) UpdateCustomerMetaData(
	customer *Customer,
	params *UpdateCustomerMetadataRequest,
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
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	request.Metadata = map[string]*hera.DataMapValue{}

	for key, value := range params.Metadata {
		request.Metadata[key] = &hera.DataMapValue{
			Value: &hera.DataMapValue_StringVal{
				StringVal: value,
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.UpdateCustomerMetadata(ctx, &request)
}

func (s *service) DeleteCustomerMetaData(
	customer *Customer,
	params *DeleteCustomerMetadataRequest,
) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerMetadataRequest
	request.OrgId = s.orgId
	request.Keys = params.Metadata

	if customer.Id != "" {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerId{
			CustomerId: customer.Id,
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: s.setCustomerNumber(customer),
		}
	}
	if !reflect.ValueOf(customer.SecondaryId).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_SecondaryId{
			SecondaryId: s.setCustomerSecondaryId(customer),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerMetadata(ctx, &request)
}

func (s *service) NewCustomer(params *CreateCustomerParams) (*Customer, error) {
	var customer Customer
	customer.Id = params.Id
	customer.CustomerNumber = params.CustomerNumber
	customer.service = s

	return &customer, nil
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
	params *UpdateCustomerTagRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerTag(c, params)
}

// DeleteTag disaccosiates a tag from a customer
func (c *Customer) DeleteTag(
	params *DeleteCustomerTagRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerTag(c, params)
}

// UpdateSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryId(
	params *UpdateCustomerSecondaryIdRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerSecondaryId(c, params)
}

// DeleteSecondaryId deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryId(
	params *DeleteCustomerSecondaryIdRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerSecondaryId(c, params)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(
	params *CustomerReminderRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AddCustomerReminder(c, params)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(
	params *CancelCustomerReminderRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.CancelCustomerReminder(c, params)
}

// UpdateMetaData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateMetaData(
	params *UpdateCustomerMetadataRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerMetaData(c, params)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(
	params *DeleteCustomerMetadataRequest,
) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerMetaData(c, params)
}

// LeaseMetaData leases customer metadata
func (c *Customer) LeaseMetaData(
	params *LeaseCustomerMetadataRequest,
) (*hera.LeaseCustomerMetadataReply, error) {
	return c.service.LeaseCustomerMetaData(c, params)
}
