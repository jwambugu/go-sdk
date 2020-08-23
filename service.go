package elariango

import (
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (

	// Reminder defines the composition of a reminder. The key is an identifier property. The payload is also a string.
	Reminder struct {
		ProductID  string    `json:"productId,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
		Key        string    `json:"key,omitempty"`
		Payload    string    `json:"payload,omitempty"`
	}

	// Tag defines a customer tag
	Tag struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}

	// Cash defines a cash object
	Cash struct {
		CurrencyCode string  `json:"currencyCode"`
		Amount       float64 `json:"amount"`
	}

	// Elarian interface exposes high level consumable elarian functionality
	Elarian interface {
		GetCustomerState(customer *Customer, params *CustomerStateRequest) (*hera.CustomerStateReplyData, error)
		AdoptCustomerState(customer *Customer, otherCustomer *Customer, params *AdoptCustomerStateRequest) (*hera.UpdateCustomerStateReply, error)
		AddCustomerReminder(customer *Customer, params *CustomerReminderRequest) (*hera.UpdateCustomerStateReply, error)
		AddCustomerReminderByTag(params *CustomerReminderByTagRequest) (*hera.TagCommandReply, error)
		CancelCustomerReminder(customer *Customer, params *CancelCustomerReminderRequest) (*hera.UpdateCustomerStateReply, error)
		CancelCustomerReminderByTag(params *CancelCustomerReminderByTagRequest) (*hera.TagCommandReply, error)
		UpdateCustomerTag(customer *Customer, params *UpdateCustomerTagRequest) (*hera.UpdateCustomerStateReply, error)
		DeleteCustomerTag(customer *Customer, params *DeleteCustomerTagRequest) (*hera.UpdateCustomerStateReply, error)
		UpdateCustomerSecondaryID(customer *Customer, params *UpdateCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error)
		DeleteCustomerSecondaryID(customer *Customer, params *DeleteCustomerSecondaryIDRequest) (*hera.UpdateCustomerStateReply, error)
		UpdateCustomerMetaData(customer *Customer, params *UpdateCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)
		DeleteCustomerMetaData(customer *Customer, params *DeleteCustomerMetadataRequest) (*hera.UpdateCustomerStateReply, error)
		GetAuthToken(appID string) (*hera.AuthTokenReply, error)
		SendWebhookResponse(params *WebhookRequest) (*hera.WebhookResponseReply, error)
		StreamNotifications(appID string) (chan *hera.WebhookRequest, chan error)
		SendPayment(customer *Customer, params *PaymentRequest) (*hera.SendPaymentReply, error)
		CheckoutPayment(customer *Customer, params *PaymentRequest) (*hera.CheckoutPaymentReply, error)
		MakeVoiceCall(customer *Customer, params *VoiceCallRequest) (*hera.MakeVoiceCallReply, error)
		SendMessage(customer *Customer, params *SendMessageRequest) (*hera.SendMessageReply, error)
		SendMessageByTag(params *SendMessageByTagRequest) (*hera.TagCommandReply, error)
		ReplyToMessage(customer *Customer, params *ReplyToMessageRequest) (*hera.SendMessageReply, error)
		MessagingConsent(customer *Customer, params *MessagingConsentRequest) (*hera.MessagingConsentReply, error)
	}
	elarian struct {
		client hera.GrpcWebServiceClient
	}
)

// NewService Creates a new Elarian service
func NewService(client *hera.GrpcWebServiceClient) Elarian {
	return &elarian{
		client: *client,
	}
}
