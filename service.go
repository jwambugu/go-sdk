package elarian

import (
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// Cash defines a cash object
	Cash struct {
		CurrencyCode string  `json:"currencyCode"`
		Amount       float64 `json:"amount"`
	}

	// Service interface exposes high level consumable elarian functionality
	Service interface {
		// GetCustomerState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIDs, payments etc.
		GetCustomerState(customer *Customer, params *CustomerStateRequest) (*hera.CustomerStateReplyData, error)

		// AdoptCustomerState copies the state of the second customer to the first customer. note for the first customer a customer id is required
		AdoptCustomerState(customer *Customer, otherCustomer *Customer, params *AdoptCustomerStateRequest) (*hera.UpdateCustomerStateReply, error)

		// AddCustomerReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
		AddCustomerReminder(customer *Customer, params *CustomerReminderRequest) (*hera.UpdateCustomerStateReply, error)

		// AddCustomerReminderByTag sets a reminder on elarian for a group of customers identified by the tag on trigger is pushed through the notification stream.
		AddCustomerReminderByTag(params *CustomerReminderByTagRequest) (*hera.TagCommandReply, error)

		// CancelCustomerReminder cancels a set reminder
		CancelCustomerReminder(customer *Customer, params *CancelCustomerReminderRequest) (*hera.UpdateCustomerStateReply, error)

		// CancelCustomerReminderByTag cancels a reminder set on a customer tag.
		CancelCustomerReminderByTag(params *CancelCustomerReminderByTagRequest) (*hera.TagCommandReply, error)

		// UpdateCustomerTag is used to add more tags to a customer
		UpdateCustomerTag(customer *Customer, params *UpdateCustomerTagRequest) (*hera.UpdateCustomerStateReply, error)

		// DeleteCustomerTag disaccosiates a tag from a customer
		DeleteCustomerTag(customer *Customer, params *DeleteCustomerTagRequest) (*hera.UpdateCustomerStateReply, error)

		// UpdateCustomerSecondaryID adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
		UpdateCustomerSecondaryID(customer *Customer, params *UpdateCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error)
		// DeleteCustomerSecondaryID deletes an associated secondary id from a customer
		DeleteCustomerSecondaryID(customer *Customer, params *DeleteCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error)

		// UpdateCustomerMetaData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerMetaData(customer *Customer, params *UpdateCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)

		// DeleteCustomerMetaData removes a customers metadata.
		DeleteCustomerMetaData(customer *Customer, params *DeleteCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)

		// GetAuthToken returns an authentication token for your application.
		GetAuthToken(appID string) (*hera.AuthTokenReply, error)

		// SendWebhookResponse func
		SendWebhookResponse(params *WebhookRequest) (*hera.WebhookResponseReply, error)

		// StreamNotifications acts as a pipe of information from elarian example of this would be reminders and payment info
		StreamNotifications(appID string) (chan *hera.WebhookRequest, chan error)

		// SendPayment requires a wallet setup and involves the transfer of funds to a customer
		SendPayment(customer *Customer, params *PaymentRequest) (*hera.SendPaymentReply, error)

		// CheckoutPayment func
		CheckoutPayment(customer *Customer, params *PaymentRequest) (*hera.CheckoutPaymentReply, error)

		// MakeVoiceCall allows you to make a voice call to a customer
		MakeVoiceCall(customer *Customer, params *VoiceCallRequest) (*hera.MakeVoiceCallReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(customer *Customer, params *SendMessageRequest) (*hera.SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(params *SendMessageByTagRequest) (*hera.TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(customer *Customer, params *ReplyToMessageRequest) (*hera.SendMessageReply, error)

		// MessagingConsent func
		MessagingConsent(customer *Customer, params *MessagingConsentRequest) (*hera.MessagingConsentReply, error)
	}

	service struct {
		client hera.GrpcWebServiceClient
	}
)

// NewService Creates a new Elarian service
func NewService(client *hera.GrpcWebServiceClient) Service {
	return &service{
		client: *client,
	}
}
