package elarian

import (
	"github.com/asaskevich/EventBus"
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

/*
	commands

	RecieveMessage
	ReceivePayment
	UpdatePaymentStatus



	Notifications

	Reminder
	MessagingSessionStarted
	MessagingSessionRenewed
	MessagingSessionEnded
	MessagingConsentUpdate
	ReceivedMessage
	MessageStatus
	SendMessageReaction
	ReceivedPayment
	PaymentStatus
	WalletPaymentStatus
	CustomerActivity


	SIM
	SendMessage
	MakeVoiceCall
	SendCustomerPayment
	SendChannelPayment
	CheckoutPayment




*/

type (
	notificationHandler func(msg payload.Payload)
	// Service interface exposes high level consumable elarian functionality
	Service interface {
		GenerateAuthToken() (*hera.GenerateAuthTokenReply, error)

		// GetCustomerState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
		GetCustomerState(customer *Customer) (*hera.GetCustomerStateReply, error)

		// AdoptCustomerState copies the state of the second customer to the first customer. note for the first customer a customer id is required
		AdoptCustomerState(customerID string, otherCustomer *Customer) (*hera.UpdateCustomerStateReply, error)

		// AddCustomerReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
		AddCustomerReminder(customer *Customer, reminder *Reminder) (*hera.UpdateCustomerAppDataReply, error)

		// AddCustomerReminderByTag sets a reminder on elarian for a group of customers identified by the tag on trigger is pushed through the notification stream.
		AddCustomerReminderByTag(tag *Tag, reminder *Reminder) (*hera.TagCommandReply, error)

		// CancelCustomerReminder cancels a set reminder
		CancelCustomerReminder(customer *Customer, key string) (*hera.UpdateCustomerAppDataReply, error)

		// CancelCustomerReminderByTag cancels a reminder set on a customer tag.
		CancelCustomerReminderByTag(tag *Tag, key string) (*hera.TagCommandReply, error)

		// UpdateCustomerTag is used to add more tags to a customer
		UpdateCustomerTag(customer *Customer, tags ...*Tag) (*hera.UpdateCustomerStateReply, error)

		// DeleteCustomerTag disaccosiates a tag from a customer
		DeleteCustomerTag(customer *Customer, keys ...string) (*hera.UpdateCustomerStateReply, error)

		// UpdateSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
		UpdateCustomerSecondaryID(customer *Customer, secondaryIds ...*SecondaryID) (*hera.UpdateCustomerStateReply, error)

		// DeleteSecondaryId deletes an associated secondary id from a customer
		DeleteCustomerSecondaryID(customer *Customer, secondaryIds ...*SecondaryID) (*hera.UpdateCustomerStateReply, error)

		// UpdateCustomerMetaData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerMetaData(customer *Customer, metadata ...*Metadata) (*hera.UpdateCustomerStateReply, error)

		// DeleteCustomerMetaData removes a customers metadata.
		DeleteCustomerMetaData(customer *Customer, keys ...string) (*hera.UpdateCustomerStateReply, error)

		// LeaseCustomerMetaData removes a customers metadata.
		LeaseCustomerAppData(customer *Customer) (*hera.LeaseCustomerAppDataReply, error)

		// UpdateCustomerAppData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerAppData(customer *Customer, metadata map[string]string) (*hera.UpdateCustomerAppDataReply, error)

		// DeleteCustomerAppData removes a customers metadata.
		DeleteCustomerAppData(customer *Customer) (*hera.UpdateCustomerAppDataReply, error)

		GetCustomerActivity(customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*hera.CustomerActivityReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(customer *Customer, channelNumber *MessagingChannelNumber, body *MessageBody) (*hera.SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(tag *Tag, channelNumber *MessagingChannelNumber, body *MessageBody) (*hera.TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(customer *Customer, messageID string, body *MessageBody) (*hera.SendMessageReply, error)

		// MessagingConsent func
		UpdateMessagingConsent(customer *CustomerNumber, channelNumber *MessagingChannelNumber, action MessagingConsentUpdate) (*hera.UpdateMessagingConsentReply, error)

		// InitiatePayment requires a wallet setup and involves the transfer of funds to a customer
		InitiatePayment(customer *Customer, params *Paymentrequest) (*hera.InitiatePaymentReply, error)

		// NewCustomer func creates and Returns a customer instance for functionality consumable from a customer's perspective
		NewCustomer(params *CreateCustomer) *Customer

		// Disconnect closes the elarian connection
		Disconnect() error

		On(event string, handler notificationHandler)
	}

	service struct {
		client rsocket.Client
		bus    EventBus.Bus
	}
)

func (s *service) Disconnect() error {
	return s.client.Close()
}

func (s *service) On(event string, handler notificationHandler) {
	s.bus.SubscribeAsync(event, handler, false)
}

// NewService Creates a new Elarian service
func NewService(options *Options, connectionOptions *ConnectionOptions) (Service, error) {
	bus := EventBus.New()
	rservice := new(rSocketService)
	rservice.host = "tcp.elarian.dev"
	rservice.port = 8082
	rservice.bus = bus

	client, err := rservice.connect(options, connectionOptions)

	elarianService := new(service)
	if err != nil {
		return &service{}, err
	}
	elarianService.client = client
	elarianService.bus = bus
	return elarianService, nil
}
