package elarian

import (
	"errors"
	"io"
	"log"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// Notification enum
	Notification int32

	// IsNotification interface every notification should implement this interface
	IsNotification interface {
		notification()
	}

	// NotificationCallBack func
	NotificationCallBack func(message IsOutBoundMessageBody, appData *Appdata)

	// NotificationHandler func
	NotificationHandler func(notification IsNotification, appData *Appdata, customer *Customer, cb NotificationCallBack)

	// NotificationPaymentStatus struct
	NotificationPaymentStatus struct {
		TransactionID string        `json:"transactionId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
	}

	// PurseNotification struct
	PurseNotification struct {
		OrgID         string                     `json:"orgId,omitempty"`
		AppID         string                     `json:"appId,omitempty"`
		CreatedAt     time.Time                  `json:"createdAt,omitempty"`
		PurseID       string                     `json:"purseId,omitempty"`
		PaymentStatus *NotificationPaymentStatus `json:"paymentStatus,omitempty"`
	}

	// MessageStatusNotification struct
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

func (s *service) notificationCallBack(body IsOutBoundMessageBody, appData *Appdata) {
	reply := new(hera.ServerToAppNotificationReply)
	if !reflect.ValueOf(appData).IsZero() {
		reply.DataUpdate = &hera.AppDataUpdate{
			Data: &hera.DataMapValue{},
		}

		if stringVal, ok := appData.Value.(StringDataValue); ok {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_StringVal{
					StringVal: string(stringVal),
				},
			}
		}
		if bytesVal, ok := appData.Value.(StringDataValue); ok {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_BytesVal{
					BytesVal: []byte(bytesVal),
				},
			}
		}
	}

	if !reflect.ValueOf(body).IsZero() {
		message := new(hera.OutboundMessage)
		reply.Message = message

		if entry, ok := body.(TextMessage); ok {
			message.Body = s.textMessage(string(entry))
		}
		if entry, ok := body.(*Template); ok {
			message.Body = s.templateMesage(entry)
		}
		if entry, ok := body.(*Location); ok {
			message.Body = s.locationMessage(entry)
		}
		if entry, ok := body.(*Media); ok {
			message.Body = s.mediaMessage(entry)
		}
		if entry, ok := body.(*UssdMenu); ok {
			message.Body = s.ussdMessage(entry)
		}
		if entry, ok := body.(*Email); ok {
			message.Body = s.email(entry)
		}
		if entry, ok := body.(VoiceCallActions); ok {
			message.Body = s.voiceMessage(entry)
		}

	}
	s.replyChannel <- reply
}

func (s *service) On(notification Notification, handler NotificationHandler) {
	s.bus.SubscribeAsync(string(notification), handler, false)
}

func (s *service) handleSimulatorNotification(notf *hera.ServerToSimulatorNotification) {
	if reflect.ValueOf(notf).IsZero() || reflect.ValueOf(notf.Entry).IsZero() {
		return
	}
	s.SendChannelPaymentSimulatorNotificationHandler(notf)
	s.CheckoutPaymentSimulatorNotificationHandler(notf)
	s.SendCustomerPaymentSimulatorNotificationHandler(notf)
	s.MakeVoiceCallSimulatorNotificationHandler(notf)
	s.SendMessageSimulatorNotificationHandler(notf)
}

func (s *service) handleNotifications(notf *hera.ServerToAppNotification) {
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

func (s *service) InitializeNotificationStream() {
	for {
		select {
		case notification := <-s.notificationChannel:
			s.handleNotifications(notification)
		case simulatorNotification := <-s.simulatorNotificationChannel:
			s.handleSimulatorNotification(simulatorNotification)
		case err := <-s.errorChannel:
			if errors.Is(err, io.EOF) {
				close(s.errorChannel)
				return
			}
			if err != nil {
				log.Println(err)
			}
		}
	}
}
