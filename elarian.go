package elarian

import (
	"context"
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

	// DataValue interface is implemented by metadata and appdata as with both you can store data as either a string or an array of bytes
	DataValue interface {
		isDataValue()
		String() string
		Bytes() []byte
	}

	// StringDataValue implements the DataValue interface represents a string
	StringDataValue string

	// BinaryDataValue implements the DataValue Interface and represents an array of bytes
	BinaryDataValue []byte

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

	// Tag defines a customer tag
	Tag struct {
		Key        string    `json:"key,omitempty"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// Metadata defines a customer's metadata oneOf value of bytes value should be provided
	Metadata struct {
		Key   string    `json:"key,omitempty"`
		Value DataValue `json:"value,omitempty"`
	}

	// Appdata defines a customer's metadata oneOf value of bytes value should be provided
	Appdata struct {
		Value DataValue `json:"value,omitempty"`
	}

	// ActivityChannelNumber defines an activity channel
	ActivityChannelNumber struct {
		Number  string          `json:"number,omitempty"`
		Channel ActivityChannel `json:"activityChannel,omitempty"`
	}

	// CustomerActivity struct
	CustomerActivity struct {
		Key        string            `json:"key,omitempty"`
		CreatedAt  time.Time         `json:"createdAt,omitempty"`
		Properties map[string]string `json:"properties,omitempty"`
	}

	// CustomerActivityReply struct
	CustomerActivityReply struct {
		Status      bool   `json:"status,omitempty"`
		Description string `json:"description,omitempty"`
		CustomerID  string `json:"customerId,omitempty"`
	}

	// UpdateCustomerStateReply struct
	UpdateCustomerStateReply struct {
		Status      bool   `json:"status,omitempty"`
		Description string `json:"description,omitempty"`
		CustomerID  string `json:"customerId,omitempty"`
	}

	// UpdateCustomerAppDataReply struct
	UpdateCustomerAppDataReply struct {
		Status      bool   `json:"status,omitempty"`
		Description string `json:"description,omitempty"`
		CustomerID  string `json:"customerId,omitempty"`
	}

	// TagCommandReply struct
	TagCommandReply struct {
		Status      bool   `json:"status,omitempty"`
		Description string `json:"description,omitempty"`
		WorkID      string `json:"workId,omitempty"`
	}

	// LeaseCustomerAppDataReply struct
	LeaseCustomerAppDataReply struct {
		Status      bool     `json:"status,omitempty"`
		Description string   `json:"description,omitempty"`
		CustomerID  string   `json:"customerId,omitempty"`
		Appdata     *Appdata `json:"appdata,omitempty"`
	}

	// UpdateMessagingConsentReply struct
	UpdateMessagingConsentReply struct {
		Status      MessagingConsentUpdateStatus `json:"status,omitempty"`
		Description string                       `json:"description,omitempty"`
		CustomerID  string                       `json:"customerId,omitempty"`
	}

	// SimulatorToServerCommandReply struct
	SimulatorToServerCommandReply struct {
		Status      bool             `json:"status,omitempty"`
		Message     *OutBoundMessage `json:"message,omitempty"`
		Description string           `json:"description,omitempty"`
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

func (StringDataValue) isDataValue() {}

func (s StringDataValue) String() string {
	return string(s)
}

// Bytes method returns a string an array of bytes
func (s StringDataValue) Bytes() []byte {
	return []byte(s)
}

func (BinaryDataValue) isDataValue() {}
func (b BinaryDataValue) String() string {
	return string(b)
}

// Bytes method returns a BinaryDataValue as an array of bytes
func (b BinaryDataValue) Bytes() []byte {
	return b
}

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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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

func (s *service) GetCustomerActivity(customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*CustomerActivityReply, error) {
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
		return nil, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.AppToServerCommandReply)
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}

	return &CustomerActivityReply{
		Status:      reply.GetCustomerActivity().Status,
		Description: reply.GetCustomerActivity().Description,
		CustomerID:  reply.GetCustomerActivity().CustomerId.Value,
	}, err
}

func (s *service) UpdateCustomerActivity(customerNumber *CustomerNumber, channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*CustomerActivityReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.CustomerActivityCommand)
	if !reflect.ValueOf(customerNumber).IsZero() {
		command.CustomerNumber = &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		}
	}
	if !reflect.ValueOf(channel).IsZero() {
		command.ChannelNumber = &hera.ActivityChannelNumber{
			Channel: hera.ActivityChannel(channel.Channel),
			Number:  channel.Number,
		}
	}
	command.SessionId = sessionID
	command.Key = key
	command.Properties = properties
	req.Entry = &hera.AppToServerCommand_CustomerActivity{
		CustomerActivity: command,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err := proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &CustomerActivityReply{
		Status:      reply.GetCustomerActivity().Status,
		Description: reply.GetCustomerActivity().Description,
		CustomerID:  reply.GetCustomerActivity().CustomerId.Value,
	}, err
}

func (s *service) AdoptCustomerState(customerID string, otherCustomer *Customer) (*UpdateCustomerStateReply, error) {
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
			OtherCustomerNumber: s.customerNumber(otherCustomer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err := proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) AddCustomerReminder(customer *Customer, reminder *Reminder) (*UpdateCustomerAppDataReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
		}
	}
	if customer.ID != "" {
		command.AddCustomerReminder.Customer = &hera.AddCustomerReminderCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}
	command.AddCustomerReminder.Reminder = &hera.CustomerReminder{
		Key:      reminder.Key,
		Payload:  wrapperspb.String(reminder.Payload),
		RemindAt: timestamppb.New(reminder.RemindAt),
	}
	if reminder.Interval != 0 {
		command.AddCustomerReminder.Reminder.Interval = durationpb.New(reminder.Interval)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	d, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(d, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}

	return &UpdateCustomerAppDataReply{
		Status:      reply.GetUpdateCustomerAppData().Status,
		Description: reply.GetUpdateCustomerAppData().Description,
		CustomerID:  reply.GetUpdateCustomerAppData().CustomerId.Value,
	}, nil
}

func (s *service) AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*TagCommandReply, error) {
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &TagCommandReply{
		Status:      reply.GetTagCommand().Status,
		Description: reply.GetTagCommand().Description,
		WorkID:      reply.GetTagCommand().WorkId.Value,
	}, nil
}

func (s *service) CancelCustomerReminder(customer *Customer, key string) (*UpdateCustomerAppDataReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err

	}
	return &UpdateCustomerAppDataReply{
		Status:      reply.GetUpdateCustomerAppData().Status,
		Description: reply.GetUpdateCustomerAppData().Description,
		CustomerID:  reply.GetUpdateCustomerAppData().CustomerId.Value,
	}, nil
}

func (s *service) CancelCustomerReminderByTag(tag *Tag, key string) (*TagCommandReply, error) {
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &TagCommandReply{
		Status:      reply.GetTagCommand().Status,
		Description: reply.GetTagCommand().Description,
	}, nil
}

func (s *service) UpdateCustomerTag(customer *Customer, tags ...*Tag) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}

	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) DeleteCustomerTag(customer *Customer, keys ...string) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, err
}

func (s *service) UpdateCustomerSecondaryID(customer *Customer, secondaryIDs ...*SecondaryID) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) DeleteCustomerSecondaryID(customer *Customer, secondaryIDs ...*SecondaryID) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}

	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) LeaseCustomerAppData(customer *Customer) (*LeaseCustomerAppDataReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte(""))).Block(ctx)
	if err != nil {
		return nil, err
	}
	commandReply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), commandReply); err != nil {
		return nil, err
	}

	reply := &LeaseCustomerAppDataReply{
		Status:      commandReply.GetLeaseCustomerAppData().Status,
		Description: commandReply.GetLeaseCustomerAppData().Description,
		CustomerID:  commandReply.GetLeaseCustomerAppData().CustomerId.Value,
		Appdata:     &Appdata{},
	}
	if val, ok := commandReply.GetLeaseCustomerAppData().Value.Value.(*hera.DataMapValue_StringVal); ok {
		reply.Appdata.Value = StringDataValue(val.StringVal)
	}
	if val, ok := commandReply.GetLeaseCustomerAppData().Value.Value.(*hera.DataMapValue_BytesVal); ok {
		byteArr := make(BinaryDataValue, len(val.BytesVal))
		for _, byteVal := range val.BytesVal {
			byteArr = append(byteArr, byteVal)
		}
		reply.Appdata.Value = byteArr
	}
	return reply, err
}

func (s *service) UpdateCustomerAppData(customer *Customer, appdata *Appdata) (*UpdateCustomerAppDataReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
		}
	}
	if customer.ID != "" {
		command.UpdateCustomerAppData.Customer = &hera.UpdateCustomerAppDataCommand_CustomerId{
			CustomerId: customer.ID,
		}
	}

	command.UpdateCustomerAppData.Update = &hera.DataMapValue{}

	if stringValue, ok := appdata.Value.(StringDataValue); ok {
		command.UpdateCustomerAppData.Update.Value = &hera.DataMapValue_StringVal{
			StringVal: string(stringValue),
		}
	}
	if binaryValue, ok := appdata.Value.(BinaryDataValue); ok {
		command.UpdateCustomerAppData.Update.Value = &hera.DataMapValue_BytesVal{
			BytesVal: binaryValue,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerAppDataReply{
		Status:      reply.GetUpdateCustomerAppData().Status,
		Description: reply.GetUpdateCustomerAppData().Description,
		CustomerID:  reply.GetUpdateCustomerAppData().CustomerId.Value,
	}, nil
}

func (s *service) DeleteCustomerAppData(customer *Customer) (*UpdateCustomerAppDataReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}

	return &UpdateCustomerAppDataReply{
		Status:      reply.GetUpdateCustomerAppData().Status,
		Description: reply.GetUpdateCustomerAppData().Description,
		CustomerID:  reply.GetUpdateCustomerAppData().CustomerId.Value,
	}, nil
}

func (s *service) UpdateCustomerMetaData(customer *Customer, metadata ...*Metadata) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		if binaryValue, ok := val.Value.(BinaryDataValue); ok {
			mapValue.Value = &hera.DataMapValue_BytesVal{
				BytesVal: binaryValue,
			}
		}
		if stringValue, ok := val.Value.(StringDataValue); ok {
			mapValue.Value = &hera.DataMapValue_StringVal{
				StringVal: string(stringValue),
			}
		}
		meta[val.Key] = mapValue
	}
	command.UpdateCustomerMetadata.Updates = meta

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}

	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) DeleteCustomerMetaData(customer *Customer, keys ...string) (*UpdateCustomerStateReply, error) {
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
			CustomerNumber: s.customerNumber(customer.CustomerNumber),
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *service) UpdateMessagingConsent(customerNumber *CustomerNumber, channelNumber *MessagingChannelNumber, update MessagingConsentUpdate) (*UpdateMessagingConsentReply, error) {
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
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateMessagingConsentReply{
		Status:      MessagingConsentUpdateStatus(reply.GetUpdateMessagingConsent().Status),
		Description: reply.GetUpdateMessagingConsent().Description,
		CustomerID:  reply.GetUpdateMessagingConsent().CustomerId.Value,
	}, err
}
