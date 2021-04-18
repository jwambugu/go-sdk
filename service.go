package elarian

import (
	"context"

	"github.com/asaskevich/EventBus"
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go"
)

type (
	// Elarian interface exposes high level consumable elarian functionality
	Elarian interface {
		GenerateAuthToken(ctx context.Context) (*GenerateAuthTokenReply, error)

		// GetCustomerState returns a customers state on elarian, the state could me messaging state, metadata, secondaryIds, payments etc.
		GetCustomerState(ctx context.Context, customer IsCustomer) (*hera.GetCustomerStateReply, error)

		// AdoptCustomerState copies the state of the second customer to the first customer. note for the first customer a customer id is required
		AdoptCustomerState(ctx context.Context, customerID string, otherCustomer IsCustomer) (*UpdateCustomerStateReply, error)

		// AddCustomerReminder sets a reminder on elarian for a customer which is triggered on set time. The reminder is push through the notification stream.
		AddCustomerReminder(ctx context.Context, customer IsCustomer, reminder *Reminder) (*UpdateCustomerAppDataReply, error)

		// AddCustomerReminderByTag sets a reminder on elarian for a group of customers identified by the tag on trigger is pushed through the notification stream.
		AddCustomerReminderByTag(ctx context.Context, tag *Tag, reminder *Reminder) (*TagCommandReply, error)

		// CancelCustomerReminder cancels a set reminder
		CancelCustomerReminder(ctx context.Context, customer IsCustomer, key string) (*UpdateCustomerAppDataReply, error)

		// CancelCustomerReminderByTag cancels a reminder set on a customer tag.
		CancelCustomerReminderByTag(ctx context.Context, tag *Tag, key string) (*TagCommandReply, error)

		// UpdateCustomerTag is used to add more tags to a customer
		UpdateCustomerTag(ctx context.Context, customer IsCustomer, tags ...*Tag) (*UpdateCustomerStateReply, error)

		// DeleteCustomerTag disaccosiates a tag from a customer
		DeleteCustomerTag(ctx context.Context, customer IsCustomer, keys ...string) (*UpdateCustomerStateReply, error)

		// UpdateSecondaryId adds secondary ids to a customer, this could be the id you associate the customer with locally on your application.
		UpdateCustomerSecondaryID(ctx context.Context, customer IsCustomer, secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error)

		// DeleteSecondaryId deletes an associated secondary id from a customer
		DeleteCustomerSecondaryID(ctx context.Context, customer IsCustomer, secondaryIds ...*SecondaryID) (*UpdateCustomerStateReply, error)

		// UpdateCustomerMetaData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerMetaData(ctx context.Context, customer IsCustomer, metadata ...*Metadata) (*UpdateCustomerStateReply, error)

		// DeleteCustomerMetaData removes a customers metadata.
		DeleteCustomerMetaData(ctx context.Context, customer IsCustomer, keys ...string) (*UpdateCustomerStateReply, error)

		// LeaseCustomerMetaData removes a customers metadata.
		LeaseCustomerAppData(ctx context.Context, customer IsCustomer) (*LeaseCustomerAppDataReply, error)

		// UpdateCustomerAppData adds abitrary or application specific information that you may want to tie to a customer.
		UpdateCustomerAppData(ctx context.Context, customer IsCustomer, appdata *Appdata) (*UpdateCustomerAppDataReply, error)

		// DeleteCustomerAppData removes a customers metadata.
		DeleteCustomerAppData(ctx context.Context, customer IsCustomer) (*UpdateCustomerAppDataReply, error)

		// UpdateCustomerActivity func
		UpdateCustomerActivity(ctx context.Context, customerNumber *CustomerNumber, channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*CustomerActivityReply, error)

		// GetCustomerActivity func
		GetCustomerActivity(ctx context.Context, customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*CustomerActivityReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(ctx context.Context, customer *CustomerNumber, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(ctx context.Context, tag *Tag, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(ctx context.Context, customerID, messageID string, body IsOutBoundMessageBody) (*SendMessageReply, error)

		// MessagingConsent func
		UpdateMessagingConsent(ctx context.Context, customer *CustomerNumber, channelNumber *MessagingChannelNumber, action MessagingConsentUpdate) (*UpdateMessagingConsentReply, error)

		// InitiatePayment requires a wallet setup and involves the transfer of funds to a customer
		InitiatePayment(ctx context.Context, customer *Customer, params *Paymentrequest) (*InitiatePaymentReply, error)

		// NewCustomer func creates and Returns a customer instance for functionality consumable from a customer's perspective
		NewCustomer(params *CreateCustomer) *Customer

		// Disconnect closes the elarian connection
		Disconnect() error

		// InitializeNotificationStream starts listening for notifications if notifications are enabled
		InitializeNotificationStream() <-chan error

		// On registers an event to a notification handler
		On(event Notification, handler NotificationHandler)

		// ReceiveMessage is a simulator method that can be used to ReceiveMessage messages from a custom simulator
		ReceiveMessage(ctx context.Context, customerNumber string, channel *MessagingChannelNumber, sessionID string, parts []*InBoundMessageBody) (*SimulatorToServerCommandReply, error)

		// ReceivePayment is a simulator method that can be used to ReceivePayment messages from a custom simulator
		ReceivePayment(ctx context.Context, customerNumber, transactionID string, channel *PaymentChannelNumber, cash *Cash, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error)

		// UpdatePaymentStatus is a simulator method that can be used to update a payment's status from a custom simulator
		UpdatePaymentStatus(ctx context.Context, transactionID string, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error)
	}

	elarian struct {
		client                       rsocket.Client
		bus                          EventBus.Bus
		errorChannel                 <-chan error
		replyChannel                 chan<- *hera.ServerToAppNotificationReply
		notificationChannel          <-chan *hera.ServerToAppNotification
		simulatorNotificationChannel <-chan *hera.ServerToSimulatorNotification
	}
)

func (s *elarian) Disconnect() error {
	return s.client.Close()
}

// NewService Creates a new Elarian service
func NewService(options *Options, connectionOptions *ConnectionOptions) (Elarian, error) {
	errorChan := make(chan error)
	replyChan := make(chan *hera.ServerToAppNotificationReply)
	notificationChannel := make(chan *hera.ServerToAppNotification)
	simulatorNotificationChannel := make(chan *hera.ServerToSimulatorNotification)

	srvc := &service{
		host:                         "tcp.elarian.dev",
		port:                         8082,
		errorChannel:                 errorChan,
		replyChannel:                 replyChan,
		notificationChannel:          notificationChannel,
		simulatorNotificationChannel: simulatorNotificationChannel,
	}

	client, err := srvc.connect(options, connectionOptions)
	if err != nil {
		return nil, err
	}
	return &elarian{
		client:                       client,
		bus:                          EventBus.New(),
		errorChannel:                 errorChan,
		replyChannel:                 replyChan,
		notificationChannel:          notificationChannel,
		simulatorNotificationChannel: simulatorNotificationChannel,
	}, nil
}
