package elarian

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// NumberProvider is an enum that defines a type of customer number provider. it could be a telco, facebook, telegram or unspecified
	NumberProvider int32

	// ActivityChannel is an enum that defines a type of activity  channel. it could be a web and mobile
	ActivityChannel int32

	// MessagingConsentUpdate Enum
	MessagingConsentUpdate int32

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

	// Customer struct defines the parameters required to make any command involving a customer. Note: in every scenario either the Id or the phoneNumber is required but not  both unless otherwise specified
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
		Interval time.Duration `json:"interval,omitempty"`
		Key      string        `json:"key,omitempty"`
		Payload  string        `json:"payload,omitempty"`
		RemindAt time.Time     `json:"expiration,omitempty"`
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

	// Metadata defines a customer's metadata oneOf value of bytes value should be provided
	Metadata struct {
		Key        string `json:"key,omitempty"`
		Value      string `json:"value,omitempty"`
		BytesValue []byte `json:"bytesValue,omitempty"`
	}

	// ActivityChannelNumber defines an activity channel
	ActivityChannelNumber struct {
		Number  string          `json:"number,omitempty"`
		Channel ActivityChannel `json:"activityChannel,omitempty"`
	}
)

// CustomerNumberProvider Enums
const (
	CustomerNumberProviderUnspecified NumberProvider = iota
	CustomerNumberProviderFacebook
	CustomerNumberProviderCellular
	CustomerNumberProviderTelegram
)

// ActivityChannel Enums
const (
	ActivityChannelUnspecified ActivityChannel = iota
	ActivityChannelWeb
	ActivityChannelMobile
)

// MessagingConsent Enums
const (
	MessagingConsentUpdateUnspecified MessagingConsentUpdate = iota
	MessagingConsentUpdateAllow
	MessagingConsentUpdateBlock
)

func (s *service) GetCustomerState(customer *Customer) (*hera.GetCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_GetCustomerState)
	command.GetCustomerState = &hera.GetCustomerStateCommand{}
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.GetCustomerState.Customer = &hera.GetCustomerStateCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: &wrapperspb.StringValue{
				Value: customer.SecondaryID.Key,
			}},
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.GetCustomerState.Customer = &hera.GetCustomerStateCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.GetCustomerState.Customer = &hera.GetCustomerStateCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.GetCustomerStateReply{}, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.AppToServerCommandReply)
	if err != nil {
		return &hera.GetCustomerStateReply{}, err
	}
	err = proto.Unmarshal(payload.Data(), reply)
	return reply.GetGetCustomerState(), err
}

func (s *service) GetCustomerActivity(customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*hera.CustomerActivityReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_CustomerActivity)
	command.CustomerActivity = &hera.CustomerActivityCommand{}
	req.Entry = command

	if !reflect.ValueOf(customerNumber).IsZero() {
		command.CustomerActivity.CustomerNumber = &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		}
	}
	command.CustomerActivity.ChannelNumber = &hera.ActivityChannelNumber{
		Channel: hera.ActivityChannel(channelNumber.Channel),
		Number:  channelNumber.Number,
	}
	command.CustomerActivity.SessionId = sessionID

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.CustomerActivityReply{}, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.AppToServerCommandReply)
	if err != nil {
		return &hera.CustomerActivityReply{}, err
	}
	err = proto.Unmarshal(payload.Data(), reply)
	return reply.GetCustomerActivity(), err
}

func (s *service) AdoptCustomerState(customerID string, otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_AdoptCustomerState)
	command.AdoptCustomerState = &hera.AdoptCustomerStateCommand{}
	req.Entry = command

	command.AdoptCustomerState.CustomerId = customerID

	if !reflect.ValueOf(otherCustomer.SecondaryID).IsZero() {
		command.AdoptCustomerState.OtherCustomer = &hera.
			AdoptCustomerStateCommand_OtherSecondaryId{
			OtherSecondaryId: s.secondaryID(otherCustomer),
		}
	}
	if !reflect.ValueOf(otherCustomer.CustomerNumber).IsZero() {
		command.AdoptCustomerState.OtherCustomer = &hera.AdoptCustomerStateCommand_OtherCustomerNumber{
			OtherCustomerNumber: s.customerNumber(otherCustomer),
		}
	}
	if otherCustomer.ID != "" {
		command.AdoptCustomerState.OtherCustomer = &hera.
			AdoptCustomerStateCommand_OtherCustomerId{
			OtherCustomerId: otherCustomer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) AddCustomerReminder(customer *Customer, reminder *Reminder) (*hera.UpdateCustomerAppDataReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_AddCustomerReminder)
	command.AddCustomerReminder = new(hera.AddCustomerReminderCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.AddCustomerReminder.Customer = &hera.AddCustomerReminderCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.AddCustomerReminder.Customer = &hera.AddCustomerReminderCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.AddCustomerReminder.Customer = &hera.AddCustomerReminderCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	command.AddCustomerReminder.Reminder = &hera.CustomerReminder{
		Key:      reminder.Key,
		Interval: durationpb.New(reminder.Interval),
		Payload:  wrapperspb.String(reminder.Payload),
		RemindAt: timestamppb.New(reminder.RemindAt),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	d, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(d, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerAppData(), err
}

func (s *service) AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*hera.TagCommandReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_AddCustomerReminderTag)
	command.AddCustomerReminderTag = new(hera.AddCustomerReminderTagCommand)
	req.Entry = command

	command.AddCustomerReminderTag.Tag = &hera.IndexMapping{
		Key:   tag.Key,
		Value: wrapperspb.String(tag.Value),
	}
	command.AddCustomerReminderTag.Reminder = &hera.CustomerReminder{
		Key:      reminder.Key,
		Interval: durationpb.New(reminder.Interval),
		Payload:  wrapperspb.String(reminder.Payload),
		RemindAt: timestamppb.New(reminder.RemindAt),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetTagCommand(), err
}

func (s *service) CancelCustomerReminder(customer *Customer, key string) (*hera.UpdateCustomerAppDataReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_CancelCustomerReminder)
	command.CancelCustomerReminder = new(hera.CancelCustomerReminderCommand)
	req.Entry = command

	command.CancelCustomerReminder.Key = key
	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.CancelCustomerReminder.Customer = &hera.CancelCustomerReminderCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.CancelCustomerReminder.Customer = &hera.CancelCustomerReminderCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.CancelCustomerReminder.Customer = &hera.CancelCustomerReminderCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerAppData(), err
}

func (s *service) CancelCustomerReminderByTag(tag *Tag, key string) (*hera.TagCommandReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_CancelCustomerReminderTag)
	command.CancelCustomerReminderTag = new(hera.CancelCustomerReminderTagCommand)
	req.Entry = command

	command.CancelCustomerReminderTag.Key = key
	command.CancelCustomerReminderTag.Tag = &hera.IndexMapping{
		Key:   tag.Key,
		Value: wrapperspb.String(tag.Value),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	var reply = new(hera.AppToServerCommandReply)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetTagCommand(), err
}

func (s *service) UpdateCustomerTag(customer *Customer, tags ...*Tag) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_UpdateCustomerTag)
	command.UpdateCustomerTag = new(hera.UpdateCustomerTagCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.UpdateCustomerTag.Customer = &hera.UpdateCustomerTagCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.UpdateCustomerTag.Customer = &hera.UpdateCustomerTagCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.UpdateCustomerTag.Customer = &hera.UpdateCustomerTagCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	heraTags := []*hera.CustomerIndex{}

	for _, tag := range tags {
		heraTags = append(heraTags, &hera.CustomerIndex{
			ExpiresAt: timestamppb.New(tag.Expiration),
			Mapping: &hera.IndexMapping{
				Key:   tag.Key,
				Value: wrapperspb.String(tag.Value),
			},
		})
	}
	command.UpdateCustomerTag.Updates = heraTags

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}

	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) DeleteCustomerTag(customer *Customer, keys ...string) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_DeleteCustomerTag)
	command.DeleteCustomerTag = new(hera.DeleteCustomerTagCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.DeleteCustomerTag.Customer = &hera.DeleteCustomerTagCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.DeleteCustomerTag.Customer = &hera.DeleteCustomerTagCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.DeleteCustomerTag.Customer = &hera.DeleteCustomerTagCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	command.DeleteCustomerTag.Deletions = keys

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) UpdateCustomerSecondaryID(customer *Customer, secondaryIDs ...*SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	heraSecIDs := []*hera.CustomerIndex{}
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_UpdateCustomerSecondaryId)
	command.UpdateCustomerSecondaryId = new(hera.UpdateCustomerSecondaryIdCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.UpdateCustomerSecondaryId.Customer = &hera.UpdateCustomerSecondaryIdCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.UpdateCustomerSecondaryId.Customer = &hera.UpdateCustomerSecondaryIdCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.UpdateCustomerSecondaryId.Customer = &hera.UpdateCustomerSecondaryIdCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	for _, secondaryID := range secondaryIDs {
		heraSecIDs = append(heraSecIDs, &hera.CustomerIndex{
			ExpiresAt: timestamppb.New(secondaryID.Expiration),
			Mapping: &hera.IndexMapping{
				Key:   secondaryID.Key,
				Value: wrapperspb.String(secondaryID.Value),
			},
		})
	}
	command.UpdateCustomerSecondaryId.Updates = heraSecIDs
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) DeleteCustomerSecondaryID(customer *Customer, secondaryIDs ...*SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_DeleteCustomerSecondaryId)
	command.DeleteCustomerSecondaryId = new(hera.DeleteCustomerSecondaryIdCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.DeleteCustomerSecondaryId.Customer = &hera.DeleteCustomerSecondaryIdCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.DeleteCustomerSecondaryId.Customer = &hera.
			DeleteCustomerSecondaryIdCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}

	if customer.ID != "" {
		command.DeleteCustomerSecondaryId.Customer = &hera.DeleteCustomerSecondaryIdCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}

	heraSecIDs := []*hera.IndexMapping{}
	for _, secondaryID := range secondaryIDs {
		heraSecIDs = append(heraSecIDs, &hera.IndexMapping{
			Key:   secondaryID.Key,
			Value: wrapperspb.String(secondaryID.Value),
		})
	}
	command.DeleteCustomerSecondaryId.Deletions = heraSecIDs

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) LeaseCustomerAppData(customer *Customer) (*hera.LeaseCustomerAppDataReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_LeaseCustomerAppData)
	command.LeaseCustomerAppData = new(hera.LeaseCustomerAppDataCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.LeaseCustomerAppData.Customer = &hera.LeaseCustomerAppDataCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.LeaseCustomerAppData.Customer = &hera.LeaseCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.LeaseCustomerAppData.Customer = &hera.LeaseCustomerAppDataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.LeaseCustomerAppDataReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte(""))).Block(ctx)
	if err != nil {
		return &hera.LeaseCustomerAppDataReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetLeaseCustomerAppData(), err
}

func (s *service) UpdateCustomerAppData(customer *Customer, appdata map[string]string) (*hera.UpdateCustomerAppDataReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_UpdateCustomerAppData)
	command.UpdateCustomerAppData = new(hera.UpdateCustomerAppDataCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.UpdateCustomerAppData.Customer = &hera.UpdateCustomerAppDataCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.UpdateCustomerAppData.Customer = &hera.UpdateCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.UpdateCustomerAppData.Customer = &hera.UpdateCustomerAppDataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	jsonData, err := json.Marshal(appdata)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	command.UpdateCustomerAppData.Update = &hera.DataMapValue{
		Value: &hera.DataMapValue_StringVal{
			StringVal: string(jsonData),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerAppData(), err
}

func (s *service) DeleteCustomerAppData(customer *Customer) (*hera.UpdateCustomerAppDataReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_DeleteCustomerAppData)
	command.DeleteCustomerAppData = new(hera.DeleteCustomerAppDataCommand)
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.DeleteCustomerAppData.Customer = &hera.DeleteCustomerAppDataCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.DeleteCustomerAppData.Customer = &hera.DeleteCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.DeleteCustomerAppData.Customer = &hera.DeleteCustomerAppDataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerAppDataReply{}, err
	}
	reply := &hera.AppToServerCommandReply{}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerAppData(), err
}

func (s *service) UpdateCustomerMetaData(customer *Customer, metadata ...*Metadata) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_UpdateCustomerMetadata)
	command.UpdateCustomerMetadata = &hera.UpdateCustomerMetadataCommand{}
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.UpdateCustomerMetadata.Customer = &hera.UpdateCustomerMetadataCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.UpdateCustomerMetadata.Customer = &hera.UpdateCustomerMetadataCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.UpdateCustomerMetadata.Customer = &hera.UpdateCustomerMetadataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	meta := map[string]*hera.DataMapValue{}

	for _, val := range metadata {
		mapValue := new(hera.DataMapValue)
		if len(val.BytesValue) > 0 {
			mapValue.Value = &hera.DataMapValue_BytesVal{
				BytesVal: val.BytesValue,
			}
		}
		if val.Value != "" {
			mapValue.Value = &hera.DataMapValue_StringVal{
				StringVal: val.Value,
			}
		}
		meta[val.Key] = mapValue
	}
	command.UpdateCustomerMetadata.Updates = meta

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := &hera.AppToServerCommandReply{}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) DeleteCustomerMetaData(customer *Customer, keys ...string) (*hera.UpdateCustomerStateReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_DeleteCustomerMetadata)
	command.DeleteCustomerMetadata = &hera.DeleteCustomerMetadataCommand{}
	req.Entry = command

	if !reflect.ValueOf(customer.SecondaryID).IsZero() {
		command.DeleteCustomerMetadata.Customer = &hera.DeleteCustomerMetadataCommand_SecondaryId{
			SecondaryId: s.secondaryID(customer),
		}
	}
	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.DeleteCustomerMetadata.Customer = &hera.DeleteCustomerMetadataCommand_CustomerNumber{
			CustomerNumber: s.customerNumber(customer),
		}
	}
	if customer.ID != "" {
		command.DeleteCustomerMetadata.Customer = &hera.DeleteCustomerMetadataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	command.DeleteCustomerMetadata.Deletions = keys

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateCustomerStateReply{}, err
	}
	reply := &hera.AppToServerCommandReply{}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateCustomerState(), err
}

func (s *service) UpdateMessagingConsent(customerNumber *CustomerNumber, channelNumber *MessagingChannelNumber, update MessagingConsentUpdate) (*hera.UpdateMessagingConsentReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_UpdateMessagingConsent)
	command.UpdateMessagingConsent = &hera.UpdateMessagingConsentCommand{}
	req.Entry = command

	if !reflect.ValueOf(customerNumber).IsZero() {
		command.UpdateMessagingConsent.CustomerNumber = &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		}
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		command.UpdateMessagingConsent.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}

	command.UpdateMessagingConsent.Update = hera.MessagingConsentUpdate(update)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.UpdateMessagingConsentReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.UpdateMessagingConsentReply{}, err
	}
	reply := &hera.AppToServerCommandReply{}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetUpdateMessagingConsent(), err
}
