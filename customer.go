package elarian

import (
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
func (c *Customer) GetState() (*hera.GetCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.GetCustomerState(c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.GetCustomerState(c.SecondaryID)
	}
	return c.service.GetCustomerState(CustomerID(c.ID))
}

// AdoptState copies the state of the second customer to this customer
func (c *Customer) AdoptState(otherCustomer IsCustomer) (*UpdateCustomerStateReply, error) {
	return c.service.AdoptCustomerState(c.ID, otherCustomer)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	return c.service.SendMessage(c.CustomerNumber, channelNumber, body)
}

// ReplyToMessage replys to a message sent by the customer
func (c *Customer) ReplyToMessage(messageID string, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	return c.service.ReplyToMessage(c.ID, messageID, body)
}

// UpdateActivity func
func (c *Customer) UpdateActivity(channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*CustomerActivityReply, error) {
	return c.service.UpdateCustomerActivity(c.CustomerNumber, channel, sessionID, key, properties)
}

// UpdateMesssagingConsent func
func (c *Customer) UpdateMesssagingConsent(channel *MessagingChannelNumber, action MessagingConsentUpdate) (*UpdateMessagingConsentReply, error) {
	return c.service.UpdateMessagingConsent(c.CustomerNumber, channel, action)
}

// LeaseAppData leases customer metadata
func (c *Customer) LeaseAppData() (*LeaseCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.LeaseCustomerAppData(c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.LeaseCustomerAppData(c.SecondaryID)
	}
	return c.service.LeaseCustomerAppData(CustomerID(c.ID))
}

// UpdateAppData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateAppData(appdata *Appdata) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerAppData(c.CustomerNumber, appdata)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerAppData(c.SecondaryID, appdata)
	}
	return c.service.UpdateCustomerAppData(CustomerID(c.ID), appdata)
}

// DeleteAppData removes a customers metadata
func (c *Customer) DeleteAppData() (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerAppData(c.CustomerNumber)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerAppData(c.SecondaryID)
	}
	return c.service.DeleteCustomerAppData(CustomerID(c.ID))
}

// UpdateMetaData adds abitrary information you want to tie to a customer
func (c *Customer) UpdateMetaData(metadata ...*Metadata) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerMetaData(c.CustomerNumber, metadata...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerMetaData(c.SecondaryID, metadata...)
	}
	return c.service.UpdateCustomerMetaData(CustomerID(c.ID), metadata...)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(keys ...string) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerMetaData(c.CustomerNumber, keys...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerMetaData(c.SecondaryID, keys...)
	}
	return c.service.DeleteCustomerMetaData(CustomerID(c.ID), keys...)
}

// UpdateTags is used to add more tags to a customer
func (c *Customer) UpdateTags(tags ...*Tag) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerTag(c.CustomerNumber, tags...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerTag(c.SecondaryID, tags...)
	}
	return c.service.UpdateCustomerTag(CustomerID(c.ID), tags...)
}

// DeleteTags disaccosiates a tag from a customer
func (c *Customer) DeleteTags(keys ...string) (*UpdateCustomerStateReply, error) {
	if !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerTag(c.CustomerNumber, keys...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerTag(c.SecondaryID, keys...)
	}
	return c.service.DeleteCustomerTag(CustomerID(c.ID), keys...)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(reminder *Reminder) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.AddCustomerReminder(c.CustomerNumber, reminder)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.AddCustomerReminder(c.SecondaryID, reminder)
	}
	return c.service.AddCustomerReminder(CustomerID(c.ID), reminder)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(key string) (*UpdateCustomerAppDataReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.CancelCustomerReminder(c.CustomerNumber, key)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.CancelCustomerReminder(c.SecondaryID, key)
	}
	return c.service.CancelCustomerReminder(CustomerID(c.ID), key)
}

// GetCustomerActivity returns a customers activity
func (c *Customer) GetCustomerActivity(channelNumber *ActivityChannelNumber, sessionID string) (*CustomerActivityReply, error) {
	return c.service.GetCustomerActivity(c.CustomerNumber, channelNumber, sessionID)
}

// UpdateSecondaryID adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryID(secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.UpdateCustomerSecondaryID(c.CustomerNumber, secondaryIds...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.UpdateCustomerSecondaryID(c.SecondaryID, secondaryIds...)
	}
	return c.service.UpdateCustomerSecondaryID(CustomerID(c.ID), secondaryIds...)
}

// DeleteSecondaryID deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryID(secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error) {
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		return c.service.DeleteCustomerSecondaryID(c.CustomerNumber, secondaryIds...)
	}
	if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		return c.service.DeleteCustomerSecondaryID(c.SecondaryID, secondaryIds...)
	}
	return c.service.DeleteCustomerSecondaryID(CustomerID(c.ID), secondaryIds...)
}

// GetMetadata returns customer metadata
func (c *Customer) GetMetadata() (map[string]*Metadata, error) {
	var (
		state *hera.GetCustomerStateReply
		err   error
	)
	if c.CustomerNumber != nil && !reflect.ValueOf(c.CustomerNumber).IsZero() {
		state, err = c.service.GetCustomerState(c.CustomerNumber)
	} else if c.SecondaryID != nil && !reflect.ValueOf(c.SecondaryID).IsZero() {
		state, err = c.service.GetCustomerState(c.SecondaryID)
	} else {
		state, err = c.service.GetCustomerState(CustomerID(c.ID))
	}
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(state.Data.IdentityState).IsZero() {
		return nil, err
	}
	metaMap := make(map[string]*Metadata)
	metadata := state.Data.IdentityState.Metadata
	for key, value := range metadata {
		meta := &Metadata{Key: key}
		if value, ok := value.Value.(*hera.DataMapValue_StringVal); ok {
			meta.Value = StringDataValue(value.StringVal)
		}
		if value, ok := value.Value.(*hera.DataMapValue_BytesVal); ok {
			bytesArr := make(BinaryDataValue, len(value.BytesVal))
			for _, byteval := range value.BytesVal {
				bytesArr = append(bytesArr, byteval)
			}
			meta.Value = bytesArr
		}
		metaMap[key] = meta
	}
	return metaMap, nil
}
