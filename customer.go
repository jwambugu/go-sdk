package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// NumberProvider is an enum that defines a type of customer number provider. it could be a telco, facebook, telegram or unspecified
	NumberProvider int32

	// PhoneNumber struct
	PhoneNumber struct {
		Number   string         `json:"number,omitempty"`
		Provider NumberProvider `json:"provider,omitempty"`
	}

	// Customer struct defines the paramters required to make any request involving a customer. Note: in every scenario either the ID or the phoneNumber is required but not  both unless otherwise specified
	Customer struct {
		ID          string      `json:"customerId,omitempty"`
		PhoneNumber PhoneNumber `json:"phoneNumber"`
	}

	// CustomerStateRequest struct defines the arguments required to get a customers state
	CustomerStateRequest struct {
		AppID string `json:"appId,omitempty"`
	}

	// AdoptCustomerStateRequest defines the arguments required to adopt the state of one customer to another. For the other customer a customer ID or customer phone Number are required.
	AdoptCustomerStateRequest struct {
		AppID string `json:"appId,omitempty"`
	}

	// CustomerReminderRequest defines the arguments required to create a reminder based on a customer.
	CustomerReminderRequest struct {
		AppID    string   `json:"appId,omitempty"`
		Reminder Reminder `json:"reminder"`
	}

	// CancelCustomerReminderRequest defines the arguments required to cancel a reminder set on a customer
	CancelCustomerReminderRequest struct {
		AppID     string `json:"appId,omitempty"`
		Key       string `json:"key,omitempty"`
		ProductID string `json:"reminder"`
	}

	// CustomerReminderByTagRequest defines the arguments required to Create a reminder based on a tag. With this the reminder can applly to a group of customers.
	CustomerReminderByTagRequest struct {
		AppID    string   `json:"appId,omitempty"`
		Reminder Reminder `json:"reminder"`
		Tag      Tag      `json:"tag"`
	}

	// CancelCustomerReminderByTagRequest defines the arguments required to cancel a reminder set on a tag
	CancelCustomerReminderByTagRequest struct {
		AppID     string `json:"appId,omitempty"`
		Key       string `json:"key,omitempty"`
		ProductID string `json:"productId,omitempty"`
		Tag       Tag    `json:"tag"`
	}

	// DeleteCustomerTagRequest defines the arguments required to delete a customer's tags
	DeleteCustomerTagRequest struct {
		AppID string   `json:"appId,omitempty"`
		Tags  []string `json:"tags"`
	}

	// UpdateCustomerTagRequest defines the arguments required to update a customer's tags, you can add one or more tags
	UpdateCustomerTagRequest struct {
		AppID string `json:"appId,omitempty"`
		Tags  []struct {
			Key        string    `json:"key,omitempty"`
			Value      string    `json:"value,omitempty"`
			Expiration time.Time `json:"expiration"`
		} `json:"tags"`
	}

	// UpdateCustomerSecondaryIDRequest defines the arguments required to update a customer's secondary IDs you can add one or more secondary IDs
	UpdateCustomerSecondaryIDRequest struct {
		AppID        string `json:"appId,omitempty"`
		SecondaryIDs []struct {
			Key        string    `json:"key,omitempty"`
			Value      string    `json:"value,omitempty"`
			Expiration time.Time `json:"expiration"`
		} `json:"secondaryIds"`
	}

	// DeleteCustomerSecondaryIDRequest defines the arguments required to delete a customer's secondary Identifiers. You can provide one or more secondary IDs you want to delete
	DeleteCustomerSecondaryIDRequest struct {
		AppID        string `json:"appId,omitempty"`
		SecondaryIDs []struct {
			Key   string `json:"key,omitempty"`
			Value string `json:"value,omitempty"`
		} `json:"secondaryIds"`
	}

	// UpdateCustomerMetadataRequest defines the arguments required to update a customer's metadata
	UpdateCustomerMetadataRequest struct {
		AppID    string            `json:"appId,omitempty"`
		Metadata map[string]string `json:"metadata"`
	}

	// DeleteCustomerMetadataRequest defines the arguments required to delete a customer's metadata
	DeleteCustomerMetadataRequest struct {
		AppID    string   `json:"appId,omitempty"`
		Metadata []string `json:"metadata"`
	}
)

const (
	// CustomerNumberProviderUnspecified is a type of NumberProvider
	CustomerNumberProviderUnspecified NumberProvider = iota
	// CustomerNumberProviderFacebook is a type of NumberProvider
	CustomerNumberProviderFacebook
	// CustomerNumberProviderTelco is a type of NumberProvider
	CustomerNumberProviderTelco
	// CustomerNumberProviderTelegram is a type of NumberProvider
	CustomerNumberProviderTelegram
)

func (s *service) GetCustomerState(customer *Customer, params *CustomerStateRequest) (*hera.CustomerStateReplyData, error) {
	var request hera.GetCustomerStateRequest

	if customer.ID != "" {
		request.Customer = &hera.GetCustomerStateRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.GetCustomerStateRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}
	request.AppId = params.AppID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := s.client.GetCustomerState(ctx, &request)
	return result.GetData(), err
}

func (s *service) AdoptCustomerState(customer *Customer, otherCustomer *Customer, params *AdoptCustomerStateRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.AdoptCustomerStateRequest
	request.AppId = params.AppID

	request.CustomerId = customer.ID

	if otherCustomer.ID != "" {
		request.OtherCustomer = &hera.AdoptCustomerStateRequest_OtherCustomerId{
			OtherCustomerId: otherCustomer.ID,
		}
	}
	if otherCustomer.PhoneNumber.Number != "" {
		request.OtherCustomer = &hera.AdoptCustomerStateRequest_OtherCustomerNumber{
			OtherCustomerNumber: &hera.CustomerNumber{
				Number:   otherCustomer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(otherCustomer.PhoneNumber.Provider),
			},
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AdoptCustomerState(ctx, &request)
}

func (s *service) AddCustomerReminder(customer *Customer, params *CustomerReminderRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.AddCustomerReminderRequest
	request.AppId = params.AppID

	if customer.ID != "" {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.AddCustomerReminderRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(params.Reminder.Expiration),
		ProductId:  params.Reminder.ProductID,
		Key:        params.Reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: params.Reminder.Payload,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AddCustomerReminder(ctx, &request)
}

func (s *service) AddCustomerReminderByTag(params *CustomerReminderByTagRequest) (*hera.TagCommandReply, error) {
	var request hera.AddCustomerReminderTagRequest
	request.AppId = params.AppID
	request.Tag = &hera.IndexMapping{
		Key: params.Tag.Key,
		Value: &wrapperspb.StringValue{
			Value: params.Tag.Value,
		},
	}
	request.Reminder = &hera.CustomerReminder{
		Expiration: timestamppb.New(params.Reminder.Expiration),
		ProductId:  params.Reminder.ProductID,
		Key:        params.Reminder.Key,
		Payload: &wrapperspb.StringValue{
			Value: params.Reminder.Payload,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AddCustomerReminderByTag(ctx, &request)
}

func (s *service) CancelCustomerReminder(customer *Customer, params *CancelCustomerReminderRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.CancelCustomerReminderRequest

	if customer.ID != "" {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.CancelCustomerReminderRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	request.AppId = params.AppID
	request.Key = params.Key
	request.ProductId = params.ProductID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.CancelCustomerReminder(ctx, &request)
}

func (s *service) CancelCustomerReminderByTag(params *CancelCustomerReminderByTagRequest) (*hera.TagCommandReply, error) {
	var request hera.CancelCustomerReminderTagRequest

	request.AppId = params.AppID
	request.Key = params.Key
	request.ProductId = params.ProductID
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

func (s *service) UpdateCustomerTag(customer *Customer, params *UpdateCustomerTagRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerTagRequest
	var tags []*hera.CustomerIndex
	request.AppId = params.AppID

	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerTagRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
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

func (s *service) DeleteCustomerTag(customer *Customer, params *DeleteCustomerTagRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerTagRequest
	request.Tags = []string{"hello"}
	request.AppId = params.AppID
	request.Tags = params.Tags

	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerTagRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerTag(ctx, &request)
}

func (s *service) UpdateCustomerSecondaryID(customer *Customer, params *UpdateCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerSecondaryIdRequest
	var secondaryIDs []*hera.CustomerIndex

	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	for _, secondaryID := range params.SecondaryIDs {
		secondaryIDs = append(secondaryIDs, &hera.CustomerIndex{
			Expiration: timestamppb.New(secondaryID.Expiration),
			Mapping: &hera.IndexMapping{
				Key: secondaryID.Key,
				Value: &wrapperspb.StringValue{
					Value: secondaryID.Value,
				},
			},
		})
	}
	request.SecondaryIds = secondaryIDs
	request.AppId = params.AppID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.UpdateCustomerSecondaryId(ctx, &request)
}

func (s *service) DeleteCustomerSecondaryID(customer *Customer, params *DeleteCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerSecondaryIdRequest
	var secondaryIDs []*hera.IndexMapping

	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerSecondaryIdRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	for _, secondaryID := range params.SecondaryIDs {
		secondaryIDs = append(secondaryIDs, &hera.IndexMapping{
			Key: secondaryID.Key,
			Value: &wrapperspb.StringValue{
				Value: secondaryID.Value,
			},
		})
	}
	request.AppId = params.AppID
	request.SecondaryIds = secondaryIDs

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerSecondaryId(ctx, &request)
}

func (s *service) UpdateCustomerMetaData(customer *Customer, params *UpdateCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.UpdateCustomerMetadataRequest

	if customer.ID != "" {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}

	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.UpdateCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	request.AppId = params.AppID
	request.Metadata = params.Metadata

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.UpdateCustomerMetadata(ctx, &request)
}

func (s *service) DeleteCustomerMetaData(customer *Customer, params *DeleteCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error) {
	var request hera.DeleteCustomerMetadataRequest
	request.AppId = params.AppID
	request.Metadata = params.Metadata

	if customer.ID != "" {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.DeleteCustomerMetadataRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.DeleteCustomerMetadata(ctx, &request)
}
