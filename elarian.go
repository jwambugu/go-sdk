package elarian

import (
	"context"
	"errors"
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

	// IsCustomer interface denotes a customer identifier which can be an id, secondaryID or a customerNumber
	IsCustomer interface {
		customer()
	}

	// StringDataValue implements the DataValue interface represents a string
	StringDataValue string

	// BinaryDataValue implements the DataValue Interface and represents an array of bytes
	BinaryDataValue []byte

	// CustomerID implements the IsCustomer interface
	CustomerID string

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
		service        Elarian
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

func (*SecondaryID) customer()    {}
func (*CustomerNumber) customer() {}
func (CustomerID) customer()      {}

func (s *elarian) GetCustomerState(customer IsCustomer) (*hera.GetCustomerStateReply, error) {
	command := &hera.GetCustomerStateCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.GetCustomerStateCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.GetCustomerStateCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.GetCustomerStateCommand_CustomerId{
			CustomerId: string(id),
		}
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_GetCustomerState{GetCustomerState: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return reply.GetGetCustomerState(), err
}

func (s *elarian) GetCustomerActivity(customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*CustomerActivityReply, error) {
	if customerNumber == nil || reflect.ValueOf(customerNumber).IsZero() {
		return nil, errors.New("customerNumber required")
	}
	if channelNumber == nil || reflect.ValueOf(channelNumber).IsZero() {
		return nil, errors.New("channelNumber required")
	}

	command := &hera.CustomerActivityCommand{
		SessionId: sessionID,
		CustomerNumber: &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		},
		ChannelNumber: &hera.ActivityChannelNumber{
			Channel: hera.ActivityChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		},
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_CustomerActivity{CustomerActivity: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &CustomerActivityReply{
		Status:      reply.GetCustomerActivity().Status,
		Description: reply.GetCustomerActivity().Description,
		CustomerID:  reply.GetCustomerActivity().CustomerId.Value,
	}, nil
}

func (s *elarian) UpdateCustomerActivity(customerNumber *CustomerNumber, channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*CustomerActivityReply, error) {
	if customerNumber == nil || reflect.ValueOf(customerNumber).IsZero() {
		return nil, errors.New("CustomerNumber required")
	}

	if channel == nil || reflect.ValueOf(channel).IsZero() {
		return nil, errors.New("ChannelNumber  required")
	}

	command := &hera.CustomerActivityCommand{
		SessionId:  sessionID,
		Key:        key,
		Properties: properties,
		CustomerNumber: &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		},
		ChannelNumber: &hera.ActivityChannelNumber{
			Channel: hera.ActivityChannel(channel.Channel),
			Number:  channel.Number,
		},
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_CustomerActivity{CustomerActivity: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err := proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &CustomerActivityReply{
		Status:      reply.GetCustomerActivity().Status,
		Description: reply.GetCustomerActivity().Description,
		CustomerID:  reply.GetCustomerActivity().CustomerId.Value,
	}, nil
}

func (s *elarian) AdoptCustomerState(customerID string, otherCustomer IsCustomer) (*UpdateCustomerStateReply, error) {
	command := &hera.AdoptCustomerStateCommand{
		CustomerId: customerID,
	}
	if secondaryID, ok := otherCustomer.(*SecondaryID); ok {
		command.OtherCustomer = &hera.AdoptCustomerStateCommand_OtherSecondaryId{
			OtherSecondaryId: &hera.IndexMapping{
				Key:   secondaryID.Key,
				Value: wrapperspb.String(secondaryID.Value),
			},
		}
	}
	if customerNumber, ok := otherCustomer.(*CustomerNumber); ok {
		command.OtherCustomer = &hera.AdoptCustomerStateCommand_OtherCustomerNumber{
			OtherCustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}

	if id, ok := otherCustomer.(CustomerID); ok {
		command.OtherCustomer = &hera.AdoptCustomerStateCommand_OtherCustomerId{
			OtherCustomerId: string(id),
		}
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_AdoptCustomerState{AdoptCustomerState: command},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err := proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &UpdateCustomerStateReply{
		Status:      reply.GetUpdateCustomerState().Status,
		Description: reply.GetUpdateCustomerState().Description,
		CustomerID:  reply.GetUpdateCustomerState().CustomerId.Value,
	}, nil
}

func (s *elarian) AddCustomerReminder(customer IsCustomer, reminder *Reminder) (*UpdateCustomerAppDataReply, error) {
	if reminder == nil || reflect.ValueOf(reminder).IsZero() {
		return nil, errors.New("Reminder Required")
	}

	command := &hera.AddCustomerReminderCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.AddCustomerReminderCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.AddCustomerReminderCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.AddCustomerReminderCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	command.Reminder = &hera.CustomerReminder{
		Key:      reminder.Key,
		Payload:  wrapperspb.String(reminder.Payload),
		RemindAt: timestamppb.New(reminder.RemindAt),
	}
	if reminder.Interval != 0 {
		command.Reminder.Interval = durationpb.New(reminder.Interval)
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_AddCustomerReminder{AddCustomerReminder: command},
	}
	d, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*TagCommandReply, error) {
	if tag == nil || reflect.ValueOf(tag).IsZero() {
		return nil, errors.New("Tag is required")
	}
	if reminder == nil || reflect.ValueOf(reminder).IsZero() {
		return nil, errors.New("Reminder is required")
	}
	command := &hera.AddCustomerReminderTagCommand{
		Tag: &hera.IndexMapping{
			Key:   tag.Key,
			Value: wrapperspb.String(tag.Value),
		},
	}
	command.Reminder = &hera.CustomerReminder{
		Key:      reminder.Key,
		Payload:  wrapperspb.String(reminder.Payload),
		RemindAt: timestamppb.New(reminder.RemindAt),
	}
	if reminder.Interval != 0 {
		command.Reminder.Interval = durationpb.New(reminder.Interval)
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_AddCustomerReminderTag{AddCustomerReminderTag: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

func (s *elarian) CancelCustomerReminder(customer IsCustomer, key string) (*UpdateCustomerAppDataReply, error) {
	command := &hera.CancelCustomerReminderCommand{
		Key: key,
	}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.CancelCustomerReminderCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.CancelCustomerReminderCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.CancelCustomerReminderCommand_CustomerId{
			CustomerId: string(id),
		}
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_CancelCustomerReminder{CancelCustomerReminder: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) CancelCustomerReminderByTag(tag *Tag, key string) (*TagCommandReply, error) {
	if tag == nil || reflect.ValueOf(tag).IsZero() {
		return nil, errors.New("Tag is required")
	}
	command := &hera.CancelCustomerReminderTagCommand{
		Key: key,
		Tag: &hera.IndexMapping{
			Key:   tag.Key,
			Value: wrapperspb.String(tag.Value),
		},
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_CancelCustomerReminderTag{CancelCustomerReminderTag: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) UpdateCustomerTag(customer IsCustomer, tags ...*Tag) (*UpdateCustomerStateReply, error) {
	command := &hera.UpdateCustomerTagCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.UpdateCustomerTagCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.UpdateCustomerTagCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.UpdateCustomerTagCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	heraTags := []*hera.CustomerIndex{}
	for _, tag := range tags {
		t := &hera.CustomerIndex{
			Mapping: &hera.IndexMapping{
				Key:   tag.Key,
				Value: wrapperspb.String(tag.Value),
			},
		}
		if !reflect.ValueOf(tag.Expiration).IsZero() {
			t.ExpiresAt = timestamppb.New(tag.Expiration)
		}
		heraTags = append(heraTags, t)
	}
	command.Updates = heraTags

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_UpdateCustomerTag{UpdateCustomerTag: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) DeleteCustomerTag(customer IsCustomer, keys ...string) (*UpdateCustomerStateReply, error) {
	command := &hera.DeleteCustomerTagCommand{
		Deletions: keys,
	}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.DeleteCustomerTagCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.DeleteCustomerTagCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.DeleteCustomerTagCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_DeleteCustomerTag{DeleteCustomerTag: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) UpdateCustomerSecondaryID(customer IsCustomer, secondaryIDs ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	command := &hera.UpdateCustomerSecondaryIdCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.UpdateCustomerSecondaryIdCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.UpdateCustomerSecondaryIdCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.UpdateCustomerSecondaryIdCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	heraSecIDs := []*hera.CustomerIndex{}
	for _, secondaryID := range secondaryIDs {
		s := &hera.CustomerIndex{
			Mapping: &hera.IndexMapping{
				Key:   secondaryID.Key,
				Value: wrapperspb.String(secondaryID.Value),
			},
		}
		if !reflect.ValueOf(secondaryID.Expiration).IsZero() {
			s.ExpiresAt = timestamppb.New(secondaryID.Expiration)
		}
		heraSecIDs = append(heraSecIDs, s)
	}
	command.Updates = heraSecIDs

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_UpdateCustomerSecondaryId{UpdateCustomerSecondaryId: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) DeleteCustomerSecondaryID(customer IsCustomer, secondaryIDs ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	command := &hera.DeleteCustomerSecondaryIdCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.DeleteCustomerSecondaryIdCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.DeleteCustomerSecondaryIdCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.DeleteCustomerSecondaryIdCommand_CustomerId{
			CustomerId: string(id),
		}
	}

	heraSecIDs := []*hera.IndexMapping{}
	for _, secondaryID := range secondaryIDs {
		heraSecIDs = append(heraSecIDs, &hera.IndexMapping{
			Key:   secondaryID.Key,
			Value: wrapperspb.String(secondaryID.Value),
		})
	}
	command.Deletions = heraSecIDs

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_DeleteCustomerSecondaryId{DeleteCustomerSecondaryId: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) LeaseCustomerAppData(customer IsCustomer) (*LeaseCustomerAppDataReply, error) {
	command := &hera.LeaseCustomerAppDataCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.LeaseCustomerAppDataCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.LeaseCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.LeaseCustomerAppDataCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_LeaseCustomerAppData{LeaseCustomerAppData: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) UpdateCustomerAppData(customer IsCustomer, appdata *Appdata) (*UpdateCustomerAppDataReply, error) {
	command := &hera.UpdateCustomerAppDataCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.UpdateCustomerAppDataCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.UpdateCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.UpdateCustomerAppDataCommand_CustomerId{
			CustomerId: string(id),
		}
	}

	command.Update = &hera.DataMapValue{}
	if stringValue, ok := appdata.Value.(StringDataValue); ok {
		command.Update.Value = &hera.DataMapValue_StringVal{
			StringVal: string(stringValue),
		}
	}
	if binaryValue, ok := appdata.Value.(BinaryDataValue); ok {
		command.Update.Value = &hera.DataMapValue_BytesVal{
			BytesVal: binaryValue,
		}
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_UpdateCustomerAppData{UpdateCustomerAppData: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) DeleteCustomerAppData(customer IsCustomer) (*UpdateCustomerAppDataReply, error) {
	command := &hera.DeleteCustomerAppDataCommand{}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.DeleteCustomerAppDataCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.DeleteCustomerAppDataCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.DeleteCustomerAppDataCommand_CustomerId{
			CustomerId: string(id),
		}
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_DeleteCustomerAppData{DeleteCustomerAppData: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) UpdateCustomerMetaData(customer IsCustomer, metadata ...*Metadata) (*UpdateCustomerStateReply, error) {
	command := &hera.UpdateCustomerMetadataCommand{}

	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.UpdateCustomerMetadataCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.UpdateCustomerMetadataCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.UpdateCustomerMetadataCommand_CustomerId{
			CustomerId: string(id),
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
	command.Updates = meta

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_UpdateCustomerMetadata{UpdateCustomerMetadata: command},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

func (s *elarian) DeleteCustomerMetaData(customer IsCustomer, keys ...string) (*UpdateCustomerStateReply, error) {
	command := &hera.DeleteCustomerMetadataCommand{
		Deletions: keys,
	}
	if secondaryID, ok := customer.(*SecondaryID); ok {
		command.Customer = &hera.DeleteCustomerMetadataCommand_SecondaryId{
			SecondaryId: &hera.IndexMapping{Value: wrapperspb.String(secondaryID.Value), Key: secondaryID.Key},
		}
	}
	if customerNumber, ok := customer.(*CustomerNumber); ok {
		command.Customer = &hera.DeleteCustomerMetadataCommand_CustomerNumber{
			CustomerNumber: s.heraCustomerNumber(customerNumber),
		}
	}
	if id, ok := customer.(CustomerID); ok {
		command.Customer = &hera.DeleteCustomerMetadataCommand_CustomerId{
			CustomerId: string(id),
		}
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_DeleteCustomerMetadata{DeleteCustomerMetadata: command},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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

func (s *elarian) UpdateMessagingConsent(customerNumber *CustomerNumber, channelNumber *MessagingChannelNumber, update MessagingConsentUpdate) (*UpdateMessagingConsentReply, error) {
	command := &hera.UpdateMessagingConsentCommand{}
	if !reflect.ValueOf(customerNumber).IsZero() {
		command.CustomerNumber = &hera.CustomerNumber{
			Provider:  hera.CustomerNumberProvider(customerNumber.Provider),
			Number:    customerNumber.Number,
			Partition: wrapperspb.String(customerNumber.Partition),
		}
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		command.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}
	command.Update = hera.MessagingConsentUpdate(update)

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_UpdateMessagingConsent{UpdateMessagingConsent: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
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
	}, nil
}
