package elarian

import (
	"github.com/asaskevich/EventBus"
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go"
)

type (
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
		UpdateCustomerAppData(customer *Customer, appdata *Appdata) (*hera.UpdateCustomerAppDataReply, error)

		// DeleteCustomerAppData removes a customers metadata.
		DeleteCustomerAppData(customer *Customer) (*hera.UpdateCustomerAppDataReply, error)

		// UpdateCustomerActivity func
		UpdateCustomerActivity(customerNumber *CustomerNumber, channel *ActivityChannelNumber, sessionID, key string, properties map[string]string) (*hera.CustomerActivityReply, error)

		// GetCustomerActivity func
		GetCustomerActivity(customerNumber *CustomerNumber, channelNumber *ActivityChannelNumber, sessionID string) (*hera.CustomerActivityReply, error)

		// SendMessage transmits a message to a customer the message body can be of different types including text, location, media and template
		SendMessage(customer *CustomerNumber, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*hera.SendMessageReply, error)

		// SendMessageByTag transmits a message to customers with the given tag. The message body can be of different types including text, location, media and template
		SendMessageByTag(tag *Tag, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*hera.TagCommandReply, error)

		// ReplyToMessage transmits a message to a customer and creates a link of two way communication with a customer that can act as a conversation history. The message body can be of different types including text, location, media and template
		ReplyToMessage(customerID, messageID string, body IsOutBoundMessageBody) (*hera.SendMessageReply, error)

		// MessagingConsent func
		UpdateMessagingConsent(customer *CustomerNumber, channelNumber *MessagingChannelNumber, action MessagingConsentUpdate) (*hera.UpdateMessagingConsentReply, error)

		// InitiatePayment requires a wallet setup and involves the transfer of funds to a customer
		InitiatePayment(customer *Customer, params *Paymentrequest) (*hera.InitiatePaymentReply, error)

		// NewCustomer func creates and Returns a customer instance for functionality consumable from a customer's perspective
		NewCustomer(params *CreateCustomer) *Customer

		// Disconnect closes the elarian connection
		Disconnect() error

		// InitializeNotificationStream starts listening for notifications if notifications are enabled
		InitializeNotificationStream()

		// On registers an event to a notification handler
		On(event Notification, handler NotificationHandler)

		// ReceiveMessage is a simulator method that can be used to ReceiveMessage messages from a custom simulator
		ReceiveMessage(customerNumber string, channel *MessagingChannelNumber, parts []*InBoundMessageBody) (*hera.SimulatorToServerCommandReply, error)

		// ReceivePayment is a simulator method that can be used to ReceivePayment messages from a custom simulator
		ReceivePayment(channel *PaymentChannelNumber, customerNumber, transactionID string) (*hera.SimulatorToServerCommandReply, error)

		// UpdatePaymentStatus is a simulator method that can be used to update a payment's status from a custom simulator
		UpdatePaymentStatus(transactionID string, paymentStatus PaymentStatus) (*hera.SimulatorToServerCommandReply, error)
	}

	service struct {
		client                       rsocket.Client
		bus                          EventBus.Bus
		errorChannel                 chan error
		replyChannel                 chan *hera.ServerToAppNotificationReply
		notificationChannel          chan *hera.ServerToAppNotification
		simulatorNotificationChannel chan *hera.ServerToSimulatorNotification
	}
)

func (s *service) Disconnect() error {
	return s.client.Close()
}

// NewService Creates a new Elarian service
func NewService(options *Options, connectionOptions *ConnectionOptions) (Service, error) {
	rservice := new(rSocketService)
	rservice.host = "tcp.elarian.dev"
	rservice.port = 8082

	errorChan := make(chan error)
	replyChan := make(chan *hera.ServerToAppNotificationReply)
	notificationChannel := make(chan *hera.ServerToAppNotification)
	simulatorNotificationChannel := make(chan *hera.ServerToSimulatorNotification)

	rservice.errorChannel = errorChan
	rservice.replyChannel = replyChan
	rservice.notificationChannel = notificationChannel
	rservice.simulatorNotificationChannel = simulatorNotificationChannel

	client, err := rservice.connect(options, connectionOptions)

	bus := EventBus.New()
	elarianService := new(service)
	if err != nil {
		return &service{}, err
	}
	elarianService.client = client
	elarianService.bus = bus

	elarianService.errorChannel = errorChan
	elarianService.replyChannel = replyChan
	elarianService.notificationChannel = notificationChannel
	elarianService.simulatorNotificationChannel = simulatorNotificationChannel

	return elarianService, nil
}
