package elarian

import (
	"context"
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *elarian) NewCustomer(params *CreateCustomer) *Customer {
	var customer Customer
	customer.ID = params.ID
	customer.CustomerNumber = params.CustomerNumber
	customer.service = s
	return &customer
}

// GetState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
func (c *Customer) GetState(ctx context.Context) (*hera.GetCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.GetCustomerState(ctx, c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.GetCustomerState(ctx, c.SecondaryID)
	}
	return c.service.GetCustomerState(ctx, CustomerID(c.ID))
}

// AdoptState copies the state of the second customer to this customer
func (c *Customer) AdoptState(ctx context.Context, otherCustomer IsCustomer) (*UpdateCustomerStateReply, error) {
	return c.service.AdoptCustomerState(ctx, c.ID, otherCustomer)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(ctx context.Context, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	return c.service.SendMessage(ctx, c.CustomerNumber, channelNumber, body)
}

// ReplyToMessage replys to a message sent by the customer
func (c *Customer) ReplyToMessage(ctx context.Context, messageID string, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	return c.service.ReplyToMessage(ctx, c.ID, messageID, body)
}

// UpdateActivity func
func (c *Customer) UpdateActivity(ctx context.Context, channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*CustomerActivityReply, error) {
	return c.service.UpdateCustomerActivity(ctx, c.CustomerNumber, channel, sessionID, key, properties)
}

// UpdateMesssagingConsent func
func (c *Customer) UpdateMesssagingConsent(ctx context.Context, channel *MessagingChannelNumber, action MessagingConsentUpdate) (*UpdateMessagingConsentReply, error) {
	return c.service.UpdateMessagingConsent(ctx, c.CustomerNumber, channel, action)
}

// LeaseAppData leases customer metadata
func (c *Customer) LeaseAppData(ctx context.Context) (*LeaseCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.LeaseCustomerAppData(ctx, c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.LeaseCustomerAppData(ctx, c.SecondaryID)
	}
	return c.service.LeaseCustomerAppData(ctx, CustomerID(c.ID))
}

// UpdateAppData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateAppData(ctx context.Context, appdata *Appdata) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerAppData(ctx, c.CustomerNumber, appdata)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerAppData(ctx, c.SecondaryID, appdata)
	}
	return c.service.UpdateCustomerAppData(ctx, CustomerID(c.ID), appdata)
}

// DeleteAppData removes a customers metadata
func (c *Customer) DeleteAppData(ctx context.Context) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerAppData(ctx, c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerAppData(ctx, c.SecondaryID)
	}
	return c.service.DeleteCustomerAppData(ctx, CustomerID(c.ID))
}

// UpdateMetaData adds abitrary information you want to tie to a customer
func (c *Customer) UpdateMetaData(ctx context.Context, metadata ...*Metadata) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerMetaData(ctx, c.CustomerNumber, metadata...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerMetaData(ctx, c.SecondaryID, metadata...)
	}
	return c.service.UpdateCustomerMetaData(ctx, CustomerID(c.ID), metadata...)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(ctx context.Context, keys ...string) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerMetaData(ctx, c.CustomerNumber, keys...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerMetaData(ctx, c.SecondaryID, keys...)
	}
	return c.service.DeleteCustomerMetaData(ctx, CustomerID(c.ID), keys...)
}

// UpdateTags is used to add more tags to a customer
func (c *Customer) UpdateTags(ctx context.Context, tags ...*Tag) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerTag(ctx, c.CustomerNumber, tags...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerTag(ctx, c.SecondaryID, tags...)
	}
	return c.service.UpdateCustomerTag(ctx, CustomerID(c.ID), tags...)
}

// DeleteTags disaccosiates a tag from a customer
func (c *Customer) DeleteTags(ctx context.Context, keys ...string) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerTag(ctx, c.CustomerNumber, keys...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerTag(ctx, c.SecondaryID, keys...)
	}
	return c.service.DeleteCustomerTag(ctx, CustomerID(c.ID), keys...)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(ctx context.Context, reminder *Reminder) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.AddCustomerReminder(ctx, c.CustomerNumber, reminder)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.AddCustomerReminder(ctx, c.SecondaryID, reminder)
	}
	return c.service.AddCustomerReminder(ctx, CustomerID(c.ID), reminder)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(ctx context.Context, key string) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.CancelCustomerReminder(ctx, c.CustomerNumber, key)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.CancelCustomerReminder(ctx, c.SecondaryID, key)
	}
	return c.service.CancelCustomerReminder(ctx, CustomerID(c.ID), key)
}

// GetCustomerActivity returns a customers activity
func (c *Customer) GetCustomerActivity(ctx context.Context, channelNumber *ActivityChannelNumber, sessionID string) (*CustomerActivityReply, error) {
	return c.service.GetCustomerActivity(ctx, c.CustomerNumber, channelNumber, sessionID)
}

// UpdateSecondaryID adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryID(ctx context.Context, secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerSecondaryID(ctx, c.CustomerNumber, secondaryIds...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerSecondaryID(ctx, c.SecondaryID, secondaryIds...)
	}
	return c.service.UpdateCustomerSecondaryID(ctx, CustomerID(c.ID), secondaryIds...)
}

// DeleteSecondaryID deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryID(ctx context.Context, secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerSecondaryID(ctx, c.CustomerNumber, secondaryIds...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerSecondaryID(ctx, c.SecondaryID, secondaryIds...)
	}
	return c.service.DeleteCustomerSecondaryID(ctx, CustomerID(c.ID), secondaryIds...)
}

// GetMetadata returns customer metadata
func (c *Customer) GetMetadata(ctx context.Context) (map[string]*Metadata, error) {
	var (
		state *hera.GetCustomerStateReply
		err   error
	)
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		state, err = c.service.GetCustomerState(ctx, c.CustomerNumber)
	} else if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		state, err = c.service.GetCustomerState(ctx, c.SecondaryID)
	} else {
		state, err = c.service.GetCustomerState(ctx, CustomerID(c.ID))
	}
	if err != nil {
		return nil, err
	}
	metaMap := make(map[string]*Metadata)
	if reflect.ValueOf(state).IsZero() || reflect.ValueOf(state.Data).IsZero() || reflect.ValueOf(state.Data.IdentityState).IsZero() || reflect.ValueOf(state.Data.IdentityState.Metadata).IsZero() {
		return metaMap, nil
	}
	metadata := state.Data.IdentityState.Metadata
	for key, value := range metadata {
		meta := &Metadata{Key: key}
		if value, ok := value.Value.(*hera.DataMapValue_StringVal); ok {
			meta.Value = value.StringVal
		}
		if value, ok := value.Value.(*hera.DataMapValue_BytesVal); ok {
			meta.BytesValue = value.BytesVal
		}
		metaMap[key] = meta
	}
	return metaMap, nil
}
