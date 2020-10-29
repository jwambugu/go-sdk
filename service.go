package elarian

import (
	"github.com/asaskevich/EventBus"
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
		// GetCustomerState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
		GetCustomerState(customer *Customer) (*hera.CustomerStateReplyData, error)

		// AdoptCustomerState copies the state of the second customer to the first customer. note for the first customer a customer id is required
		AdoptCustomerState(customer *Customer, otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error)

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

		// UpdateCustomerSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
		UpdateCustomerSecondaryId(customer *Customer, params *UpdateCustomerSecondaryIdRequest) (*hera.UpdateCustomerStateReply, error)
		// DeleteCustomerSecondaryId deletes an associated secondary id from a customer
		DeleteCustomerSecondaryId(customer *Customer, params *DeleteCustomerSecondaryIdRequest) (*hera.UpdateCustomerStateReply, error)

		// UpdateCustomerMetaData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerMetaData(customer *Customer, params *UpdateCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)

		LeaseCustomerMetaData(customer *Customer, params *LeaseCustomerMetadataRequest) (*hera.LeaseCustomerMetadataReply, error)

		// DeleteCustomerMetaData removes a customers metadata.
		DeleteCustomerMetaData(customer *Customer, params *DeleteCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)

		// GetAuthToken returns an authentication token for your application.
		GetAuthToken() (*hera.AuthTokenReply, error)

		// StreamNotifications acts as a pipe of information from elarian example of this would be reminders and payment info
		InitializeNotificationStream(appId string) error

		AddNotificationSubscriber(
			event ElarianEvent,
			handler func(data interface{}, customer *Customer),
		) error

		RemoveNotificationSubscriber(
			event ElarianEvent,
			handler func(data interface{}, customer *Customer),
		) error

		// InitiatePayment requires a wallet setup and involves the transfer of funds to a customer
		InitiatePayment(customer *Customer, params *PaymentRequest) (*hera.InitiatePaymentReply, error)

		// MakeVoiceCall allows you to make a voice call to a customer
		MakeVoiceCall(customer *Customer, params *VoiceCallRequest) (*hera.MakeVoiceCallReply, error)

		// ReplyToVoiceCall allows you to reply to a voice call
		ReplyToVoiceCall(params *VoiceCallReplyRequest) (*hera.WebhookResponseReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(customer *Customer, params *SendMessageRequest) (*hera.SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(params *SendMessageByTagRequest) (*hera.TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(customer *Customer, params *ReplyToMessageRequest) (*hera.SendMessageReply, error)

		ReplyToUSSDSession(params *USSDRequest) (*hera.WebhookResponseReply, error)

		// MessagingConsent func
		MessagingConsent(customer *Customer, params *MessagingConsentRequest) (*hera.MessagingConsentReply, error)

		// NewCustomer func creates and Returns a customer instance for functionality consumable from a customer's perspective
		NewCustomer(params *CreateCustomerParams) (*Customer, error)
	}

	service struct {
		client hera.GrpcWebServiceClient
		orgId  string
		bus    EventBus.Bus
	}
)

// NewService Creates a new Elarian service
func NewService(client *hera.GrpcWebServiceClient, orgId string) Service {
	return &service{
		client: *client,
		orgId:  orgId,
		bus:    EventBus.New(),
	}
}
