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
		AddCustomerReminder(customer *Customer, reminder *Reminder) (*hera.UpdateCustomerStateReply, error)

		// AddCustomerReminderByTag sets a reminder on elarian for a group of customers identified by the tag on trigger is pushed through the notification stream.
		AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*hera.TagCommandReply, error)

		// CancelCustomerReminder cancels a set reminder
		CancelCustomerReminder(customer *Customer, key string) (*hera.UpdateCustomerStateReply, error)

		// CancelCustomerReminderByTag cancels a reminder set on a customer tag.
		CancelCustomerReminderByTag(tag *Tag, key string) (*hera.TagCommandReply, error)

		// UpdateCustomerTag is used to add more tags to a customer
		UpdateCustomerTag(customer *Customer, tags []Tag) (*hera.UpdateCustomerStateReply, error)

		// DeleteCustomerTag disaccosiates a tag from a customer
		DeleteCustomerTag(customer *Customer, keys []string) (*hera.UpdateCustomerStateReply, error)

		// UpdateSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
		UpdateCustomerSecondaryId(customer *Customer, secondaryIds []SecondaryId) (*hera.UpdateCustomerStateReply, error)

		// DeleteSecondaryId deletes an associated secondary id from a customer
		DeleteCustomerSecondaryId(customer *Customer, secondaryIds []SecondaryId) (*hera.UpdateCustomerStateReply, error)

		// UpdateCustomerMetaData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerMetaData(customer *Customer, metadata map[string]string) (*hera.UpdateCustomerStateReply, error)

		LeaseCustomerMetaData(customer *Customer, key string) (*hera.LeaseCustomerMetadataReply, error)

		// DeleteCustomerMetaData removes a customers metadata.
		DeleteCustomerMetaData(customer *Customer, keys []string) (*hera.UpdateCustomerStateReply, error)

		// GetAuthToken returns an authentication token for your application.
		GetAuthToken() (*hera.AuthTokenReply, error)

		// StreamNotifications acts as a pipe of information from elarian example of this would be reminders and payment info
		// InitializeNotificationStream() error

		AddNotificationSubscriber(
			notification ElarianNotification,
			handler func(svc Service, cust *Customer, data interface{}),
		) error

		RemoveNotificationSubscriber(
			notification ElarianNotification,
			handler func(svc Service, cust *Customer, data interface{}),
		) error

		// InitiatePayment requires a wallet setup and involves the transfer of funds to a customer
		InitiatePayment(customer *Customer, params *Paymentrequest) (*hera.InitiatePaymentReply, error)

		// MakeVoiceCall allows you to make a voice call to a customer
		MakeVoiceCall(customer *Customer, channel *VoiceChannelNumber) (*hera.MakeVoiceCallReply, error)

		// ReplyToVoiceCall allows you to reply to a voice call
		ReplyToVoiceCall(sessionId string, actions []interface{}) (*hera.WebhookResponseReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(
			customer *Customer,
			channelNumber *MessagingChannelNumber,
			body *MessageBody,
		) (*hera.SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(
			tag *Tag,
			channelNumber *MessagingChannelNumber,
			body *MessageBody,
		) (*hera.TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(
			customer *Customer,
			messageId string,
			body *MessageBody,
		) (*hera.SendMessageReply, error)

		ReplyToUssdSession(
			sessionId string,
			ussdMenu *UssdMenu,
		) (*hera.WebhookResponseReply, error)

		// MessagingConsent func
		MessagingConsent(
			customer *Customer,
			channelNumber *MessagingChannelNumber,
			action MessagingConsentAction,
		) (*hera.MessagingConsentReply, error)

		// NewCustomer func creates and Returns a customer instance for functionality consumable from a customer's perspective
		NewCustomer(params *CreateCustomer) *Customer
	}

	service struct {
		client hera.GrpcWebServiceClient
		orgId  string
		appId  string
		bus    EventBus.Bus
	}
)

// NewService Creates a new Elarian service
func NewService(client *hera.GrpcWebServiceClient, orgId string, appId string) (Service, error) {
	srvc := service{
		client: *client,
		orgId:  orgId,
		appId:  appId,
		bus:    EventBus.New(),
	}
	errchan := srvc.initializeNotificationStream()
	err := <-errchan
	if err != nil {
		return &service{}, err
	}
	return &srvc, nil
}
