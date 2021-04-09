package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) NewCustomer(params *CreateCustomer) *Customer {
	var customer Customer
	customer.ID = params.ID
	customer.CustomerNumber = params.CustomerNumber
	customer.service = s
	return &customer
}

// GetState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
func (c *Customer) GetState() (*hera.GetCustomerStateReply, error) {
	return c.service.GetCustomerState(c)
}

// AdoptState copies the state of the second customer to this customer
func (c *Customer) AdoptState(otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error) {
	return c.service.AdoptCustomerState(c.ID, otherCustomer)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(channelNumber *MessagingChannelNumber, body *OutBoundMessageBody) (*hera.SendMessageReply, error) {
	return c.service.SendMessage(c, channelNumber, body)
}

// ReplyToMessage replys to a message sent by the customer
func (c *Customer) ReplyToMessage(messageID string, body *OutBoundMessageBody) (*hera.SendMessageReply, error) {
	return c.service.ReplyToMessage(c, messageID, body)
}

// UpdateActivity func
func (c *Customer) UpdateActivity(channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*hera.CustomerActivityReply, error) {
	return c.service.UpdateCustomerActivity(c.CustomerNumber, channel, sessionID, key, properties)
}

// UpdateMesssagingConsent func
func (c *Customer) UpdateMesssagingConsent(channel *MessagingChannelNumber, action MessagingConsentUpdate) (*hera.UpdateMessagingConsentReply, error) {
	return c.service.UpdateMessagingConsent(c.CustomerNumber, channel, action)
}

// LeaseAppData leases customer metadata
func (c *Customer) LeaseAppData() (*hera.LeaseCustomerAppDataReply, error) {
	return c.service.LeaseCustomerAppData(c)
}

// UpdateAppData adds abitrary or application specific information that you may want to tie to a customer.
func (c *Customer) UpdateAppData(appdata *Appdata) (*hera.UpdateCustomerAppDataReply, error) {
	return c.service.UpdateCustomerAppData(c, appdata)
}

// DeleteAppData removes a customers metadata
func (c *Customer) DeleteAppData() (*hera.UpdateCustomerAppDataReply, error) {
	return c.service.DeleteCustomerAppData(c)
}

// UpdateMetaData adds abitrary information you want to tie to a customer
func (c *Customer) UpdateMetaData(metadata ...*Metadata) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerMetaData(c, metadata...)
}

// DeleteMetaData removes a customers metadata
func (c *Customer) DeleteMetaData(keys ...string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerMetaData(c, keys...)
}

// UpdateTags is used to add more tags to a customer
func (c *Customer) UpdateTags(tags ...*Tag) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerTag(c, tags...)
}

// DeleteTags disaccosiates a tag from a customer
func (c *Customer) DeleteTags(keys ...string) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerTag(c, keys...)
}

// AddReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
func (c *Customer) AddReminder(reminder *Reminder) (*hera.UpdateCustomerAppDataReply, error) {
	return c.service.AddCustomerReminder(c, reminder)
}

// CancelReminder cancels a set reminder
func (c *Customer) CancelReminder(key string) (*hera.UpdateCustomerAppDataReply, error) {
	return c.service.CancelCustomerReminder(c, key)
}

// GetCustomerActivity returns a customers activity
func (c *Customer) GetCustomerActivity(channelNumber *ActivityChannelNumber, sessionID string) (*hera.CustomerActivityReply, error) {
	return c.service.GetCustomerActivity(c.CustomerNumber, channelNumber, sessionID)
}

// UpdateSecondaryID adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
func (c *Customer) UpdateSecondaryID(secondaryIds ...*SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	return c.service.UpdateCustomerSecondaryID(c, secondaryIds...)
}

// DeleteSecondaryID deletes an associated secondary id from a customer
func (c *Customer) DeleteSecondaryID(secondaryIds ...*SecondaryID) (*hera.UpdateCustomerStateReply, error) {
	return c.service.DeleteCustomerSecondaryID(c, secondaryIds...)
}

// GetMetadata returns customer metadata
func (c *Customer) GetMetadata() (map[string]*Metadata, error) {
	state, err := c.service.GetCustomerState(c)
	metaMap := make(map[string]*Metadata)

	if err != nil {
		return metaMap, err
	}
	if reflect.ValueOf(state.Data.IdentityState).IsZero() {
		return metaMap, err
	}
	metadata := state.Data.IdentityState.Metadata
	for key, value := range metadata {
		metaMap[key] = &Metadata{
			Key:        key,
			Value:      value.GetStringVal(),
			BytesValue: value.GetBytesVal(),
		}
	}
	return metaMap, nil
}
