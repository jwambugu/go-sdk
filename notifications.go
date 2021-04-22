package elarian

import (
	"errors"
	"io"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// Notification is a constant that encapsulates all kinds of notifications you can listen for.
	// Its the first argument in the On method defined by the service.
	Notification int32

	// IsNotification defines one method notification() that all notifications implement.
	// There are 18 types of notifications relating to messaging, payments and simulator.
	// They are all declared below.
	IsNotification interface {
		notification()
	}

	// NotificationCallBack is part of the notification handler and should be called after a notification has been handled.
	// It takes in a message or appData and both can be nil.
	// These are sent back to elarian. if you the callback is not called in 15 seconds an empty reply is sent back to elarian.
	NotificationCallBack func(message IsOutBoundMessageBody, appData *Appdata)

	// NotificationHandler type is a handler function for all notifications. it provides the service, the notification, appdata, customer and the callback handler defined above.
	NotificationHandler func(service Elarian, notification IsNotification, appData *Appdata, customer *Customer, cb NotificationCallBack)

	// NotificationPaymentStatus defines a structure for a payment status it has a transaction id and a status which is of type payment status
	NotificationPaymentStatus struct {
		TransactionID string        `json:"transactionId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
	}

	// PurseNotification defines a structure for a purse notification.
	// A purse notification is generated in context of a purse and has a purse ID and a paymentStatus. other fields include an OrgID, AppID and a CreatedAt
	PurseNotification struct {
		OrgID         string                     `json:"orgId,omitempty"`
		AppID         string                     `json:"appId,omitempty"`
		CreatedAt     time.Time                  `json:"createdAt,omitempty"`
		PurseID       string                     `json:"purseId,omitempty"`
		PaymentStatus *NotificationPaymentStatus `json:"paymentStatus,omitempty"`
	}

	// MessageStatusNotification defines a structure for a message status.
	// When a message is sent a message delivery status event is emitted. This is send to the SDK as a notification.
	MessageStatusNotification struct {
		CustomerID string                `json:"customerId,omitempty"`
		Status     MessageDeliveryStatus `json:"status,omitempty"`
		MessageID  string                `json:"messageId,omitempty"`
	}

	// MessageSessionStartedNotification struct
	MessageSessionStartedNotification struct {
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		SessionID      string                  `json:"SessionID,omitempty"`
		Expiration     int64                   `json:"expiration,omitempty"`
	}

	// MessageSessionRenewedNotification struct
	MessageSessionRenewedNotification struct {
		Expiration     int64                   `json:"expiration,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		SessionID      string                  `json:"SessionID,omitempty"`
	}

	// MessageSessionEndedNotification struct
	MessageSessionEndedNotification struct {
		Duration       time.Duration             `json:"duration,omitempty"`
		Reason         MessagingSessionEndReason `json:"reason,omitempty"`
		CustomerNumber *CustomerNumber           `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber   `json:"channelNumber,omitempty"`
		SessionID      string                    `json:"SessionID,omitempty"`
	}

	// MessagingConsentUpdateNotification struct
	MessagingConsentUpdateNotification struct {
		CustomerID     string                       `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber              `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber      `json:"channelNumber,omitempty"`
		Status         MessagingConsentUpdateStatus `json:"status,omitempty"`
	}

	// RecievedMessageNotification struct
	RecievedMessageNotification struct {
		MessageID      string                  `json:"messageId,omitempty"`
		SessionID      string                  `json:"sessionId,omitempty"`
		InReplyTo      string                  `json:"inReplyTo,omitempty"`
		Parts          []*InBoundMessageBody   `json:"parts,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
	}

	// SentMessageReaction struct
	SentMessageReaction struct {
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		MessageID      string                  `json:"messageId,omitempty"`
		Reaction       MessageReaction         `json:"MessageReaction,omitempty"`
	}

	// PaymentStatusNotification struct
	PaymentStatusNotification struct {
		TransactionID string        `json:"transactionId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
	}

	// ReceivedPaymentNotification struct
	ReceivedPaymentNotification struct {
		CustomerNumber *CustomerNumber       `json:"customerNumber,omitempty"`
		ChannelNumber  *PaymentChannelNumber `json:"channelNumber,omitempty"`
		PurseID        string                `json:"purseId,omitempty"`
		TransactionID  string                `json:"transactionId,omitempty"`
		Value          *Cash                 `json:"value,omitempty"`
		Status         PaymentStatus         `json:"status,omitempty"`
	}

	// WalletPaymentStatusNotification struct
	WalletPaymentStatusNotification struct {
		CustomerID    string        `json:"customerId,omitempty"`
		TransactionID string        `json:"transactionId,omitempty"`
		WalletID      string        `json:"walletId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
	}

	// CustomerActivityNotification struct
	CustomerActivityNotification struct {
		SessionID       string                 `json:"sessionId,omitempty"`
		Activity        *CustomerActivity      `json:"activity,omitempty"`
		CustomerNumber  *CustomerNumber        `json:"customerNumber,omitempty"`
		ActivityChannel *ActivityChannelNumber `json:"activityChannel,omitempty"`
	}

	// ReminderNotification struct
	ReminderNotification struct {
		Reminder *Reminder `json:"reminder,omitempty"`
		WorkID   string    `json:"workId,omitempty"`
		Tag      *Tag      `json:"tag,omitempty"`
	}

	// SendChannelPaymentSimulatorNotification struct
	SendChannelPaymentSimulatorNotification struct {
		OrgID         string                `json:"orgId,omitempty"`
		AppID         string                `json:"appId,omitempty"`
		TransactionID string                `json:"transactionId,omitempty"`
		Account       string                `json:"account,omitempty"`
		DebitParty    *PaymentParty         `json:"debitPart,omitempty"`
		ChannelNumber *PaymentChannelNumber `json:"channelNumber,omitempty"`
		Value         *Cash                 `json:"value,omitempty"`
	}

	// CheckoutPaymentSimulatorNotification struct
	CheckoutPaymentSimulatorNotification struct {
		OrgID         string                `json:"orgId,omitempty"`
		AppID         string                `json:"appId,omitempty"`
		TransactionID string                `json:"transactionId,omitempty"`
		CreditParty   *PaymentParty         `json:"creditParty,omitempty"`
		ChannelNumber *PaymentChannelNumber `json:"channelNumber,omitempty"`
		Value         *Cash                 `json:"value,omitempty"`
	}

	// SendCustomerPaymentSimulatorNotification struct
	SendCustomerPaymentSimulatorNotification struct {
		OrgID         string                `json:"orgId,omitempty"`
		AppID         string                `json:"appId,omitempty"`
		TransactionID string                `json:"transactionId,omitempty"`
		DebitParty    *PaymentParty         `json:"creditParty,omitempty"`
		ChannelNumber *PaymentChannelNumber `json:"channelNumber,omitempty"`
		Value         *Cash                 `json:"value,omitempty"`
	}

	// MakeVoiceCallSimulatorNotification struct
	MakeVoiceCallSimulatorNotification struct {
		OrgID         string                  `json:"orgId,omitempty"`
		SessionID     string                  `json:"sessionID,omitempty"`
		ChannelNumber *MessagingChannelNumber `json:"channelNumber,omitempty"`
	}

	// SendMessageSimulatorNotification struct
	SendMessageSimulatorNotification struct {
		OrgID         string                  `json:"orgId,omitempty"`
		MessageID     string                  `json:"messageId,omitempty"`
		Message       *OutBoundMessage        `json:"message,omitempty"`
		ChannelNumber *MessagingChannelNumber `json:"channelNumber,omitempty"`
	}
)

// Notification constants
const (
	ElarianReminderNotification Notification = iota
	ElarianMessageStatusNotification
	ElarianMessagingSessionStartedNotification
	ElarianMessagingSessionRenewedNotification
	ElarianMessagingSessionEndedNotification
	ElarianMessagingConsentUpdateNotification
	ElarianReceivedEmailNotification
	ElarianReceivedUssdSessionNotification
	ElarianReceivedVoiceCallNotification
	ElarianReceivedSmsNotification
	ElarianReceivedFbMessengerNotification
	ElarianReceivedTelegramNotification
	ElarianReceivedWhatsappNotification
	ElarianSentMessageReactionNotification
	ElarianReceivedPaymentNotification
	ElarianPaymentStatusNotification
	ElarianWalletPaymentStatusNotification
	ElarianCustomerActivityNotification
	ElarianPaymentPurseNotifiication
	ElarianSendChannelPaymentSimulatorNotification
	ElarianCheckoutPaymentSimulatorNotification
	ElarianSendCustomerPaymentSimulatorNotification
	ElarianMakeVoiceCallSimulatorNotification
	ElarianSendMessageSimulatorNotification
)

func (*ReminderNotification) notification()                     {}
func (*MessageStatusNotification) notification()                {}
func (*MessageSessionStartedNotification) notification()        {}
func (*MessageSessionRenewedNotification) notification()        {}
func (*MessageSessionEndedNotification) notification()          {}
func (*MessagingConsentUpdateNotification) notification()       {}
func (*RecievedMessageNotification) notification()              {}
func (*SentMessageReaction) notification()                      {}
func (*ReceivedPaymentNotification) notification()              {}
func (*PaymentStatusNotification) notification()                {}
func (*WalletPaymentStatusNotification) notification()          {}
func (*CustomerActivityNotification) notification()             {}
func (*PurseNotification) notification()                        {}
func (*UssdSessionNotification) notification()                  {}
func (*Email) notification()                                    {}
func (*Voice) notification()                                    {}
func (*InBoundMessageBody) notification()                       {}
func (*SendChannelPaymentSimulatorNotification) notification()  {}
func (*CheckoutPaymentSimulatorNotification) notification()     {}
func (*SendCustomerPaymentSimulatorNotification) notification() {}
func (*MakeVoiceCallSimulatorNotification) notification()       {}
func (*SendMessageSimulatorNotification) notification()         {}

func (s *elarian) notificationCallBack(body IsOutBoundMessageBody, appData *Appdata) {
	reply := new(hera.ServerToAppNotificationReply)
	if appData != nil && !reflect.ValueOf(appData).IsZero() {
		reply.DataUpdate = &hera.AppDataUpdate{
			Data: &hera.DataMapValue{},
		}

		if appData.Value != "" {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_StringVal{
					StringVal: appData.Value,
				},
			}
		}
		if len(appData.BytesValue) > 0 {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_BytesVal{
					BytesVal: appData.BytesValue,
				},
			}
		}
	}

	if body != nil && !reflect.ValueOf(body).IsZero() {
		message := &hera.OutboundMessage{}
		reply.Message = message
		if entry, ok := body.(TextMessage); ok {
			message.Body = s.heraOutBoundTextMessage(string(entry))
		}
		if entry, ok := body.(*Template); ok {
			message.Body = s.heraOutBoundTemplateMesage(entry)
		}
		if entry, ok := body.(*Location); ok {
			message.Body = s.heraOutBoundLocationMessage(entry)
		}
		if entry, ok := body.(*Media); ok {
			message.Body = s.heraOutBoundMediaMessage(entry)
		}
		if entry, ok := body.(*UssdMenu); ok {
			message.Body = s.heraOutBoundUssdMessage(entry)
		}
		if entry, ok := body.(*Email); ok {
			message.Body = s.heraOutBoundEmail(entry)
		}
		if entry, ok := body.(VoiceCallActions); ok {
			message.Body = s.heraOutBoundVoiceMessage(entry)
		}

	}
	s.replyChannel <- reply
}

func (s *elarian) On(notification Notification, handler NotificationHandler) {
	s.bus.Subscribe(string(notification), handler)
}

func (s *elarian) handleSimulatorNotification(notf *hera.ServerToSimulatorNotification) {
	if reflect.ValueOf(notf).IsZero() || reflect.ValueOf(notf.Entry).IsZero() {
		return
	}
	s.SendChannelPaymentSimulatorNotificationHandler(notf)
	s.CheckoutPaymentSimulatorNotificationHandler(notf)
	s.SendCustomerPaymentSimulatorNotificationHandler(notf)
	s.MakeVoiceCallSimulatorNotificationHandler(notf)
	s.SendMessageSimulatorNotificationHandler(notf)
}

func (s *elarian) handleNotifications(notf *hera.ServerToAppNotification) {
	if reflect.ValueOf(notf).IsZero() || reflect.ValueOf(notf.Entry).IsZero() {
		return
	}
	if customerNotf, ok := notf.Entry.(*hera.ServerToAppNotification_Customer); ok {
		if reflect.ValueOf(customerNotf.Customer).IsZero() {
			return
		}
		s.reminderNotificationHandler(customerNotf.Customer)
		s.messageStatusNotificationHandler(customerNotf.Customer)
		s.messagingSessionStartedNotificationHandler(customerNotf.Customer)
		s.messagingSessionRenewedNotificationHandler(customerNotf.Customer)
		s.messagingSessionEndedNotificationHandler(customerNotf.Customer)
		s.messagingConsentUpdateNotificationHandler(customerNotf.Customer)
		s.recievedMessageNotificationHandler(customerNotf.Customer)
		s.sentMesssageNotificationHandler(customerNotf.Customer)
		s.receivedPaymentNotificationHandler(customerNotf.Customer)
		s.paymentStatusNotificationHandler(customerNotf.Customer)
		s.walletPaymentStatusNotificationHandler(customerNotf.Customer)
		s.customerActivityNotificationHandler(customerNotf.Customer)
		return
	}

	if purseNotification, ok := notf.Entry.(*hera.ServerToAppNotification_Purse); ok {
		s.paymentPurseStatusNotificationHandler(purseNotification)
		return
	}
}

func (s *elarian) InitializeNotificationStream() <-chan error {
	errorChan := make(chan error)
	go func() {
		for {
			select {
			case notification := <-s.notificationChannel:
				s.handleNotifications(notification)
			case simulatorNotification := <-s.simulatorNotificationChannel:
				s.handleSimulatorNotification(simulatorNotification)
			case err := <-s.errorChannel:
				if errors.Is(err, io.EOF) {
					close(s.replyChannel)
					return
				}
				if err != nil {
					errorChan <- err
				}
			}
		}
	}()
	return errorChan
}
