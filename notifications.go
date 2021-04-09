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
	// Event enums
	Event int32

	// Notification interface every notification should implement this interface
	Notification interface {
		notification()
	}

	// NotificationCallBack func
	NotificationCallBack func(message *OutBoundMessageBody, appData *Appdata)

	// NotificationHandler func
	NotificationHandler func(notf Notification, appData *Appdata, customer *Customer, cb NotificationCallBack)

	// NotificationPaymentStatus struct
	NotificationPaymentStatus struct {
		TransactionID string
		Status        int // hera.PaymentStatus
	}

	// PurseNotification struct
	PurseNotification struct {
		OrgID         string    `json:"orgId"`
		AppID         string    `json:"appId"`
		CreatedAt     time.Time `json:"createdAt"`
		PurseID       string
		PaymentStatus *NotificationPaymentStatus
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
		Duration       time.Duration
		Reason         MessagingSessionEndReason
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		SessionID      string                  `json:"SessionID,omitempty"`
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
)

// Notification constants
const (
	ElarianReminderNotification Event = iota
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
)

func (*ReminderNotification) notification()               {}
func (*MessageStatusNotification) notification()          {}
func (*MessageSessionStartedNotification) notification()  {}
func (*MessageSessionRenewedNotification) notification()  {}
func (*MessageSessionEndedNotification) notification()    {}
func (*MessagingConsentUpdateNotification) notification() {}
func (*RecievedMessageNotification) notification()        {}
func (*SentMessageReaction) notification()                {}
func (*ReceivedPaymentNotification) notification()        {}
func (*PaymentStatusNotification) notification()          {}
func (*WalletPaymentStatusNotification) notification()    {}
func (*CustomerActivityNotification) notification()       {}
func (*PurseNotification) notification()                  {}
func (*UssdSessionNotification) notification()            {}
func (*Email) notification()                              {}
func (*Voice) notification()                              {}
func (*InBoundMessageBody) notification()                 {}

func (s *service) notificationCallBack(body *OutBoundMessageBody, appData *Appdata) {
	reply := new(hera.ServerToAppNotificationReply)
	if !reflect.ValueOf(appData).IsZero() {
		reply.DataUpdate = &hera.AppDataUpdate{
			Data: &hera.DataMapValue{},
		}
		if !reflect.ValueOf(appData.BytesValue).IsZero() {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_BytesVal{
					BytesVal: appData.BytesValue,
				},
			}
		}
		if !reflect.ValueOf(appData.Value).IsZero() {
			reply.DataUpdate.Data = &hera.DataMapValue{
				Value: &hera.DataMapValue_StringVal{
					StringVal: appData.Value,
				},
			}
		}
	}
	if !reflect.ValueOf(body).IsZero() {
		message := new(hera.OutboundMessage)
		reply.Message = message
		if body.Text != "" {
			message.Body = s.textMessage(body.Text)
		}
		if !reflect.ValueOf(body.Template).IsZero() {
			message.Body = s.templateMesage(body.Template)
		}
		if !reflect.ValueOf(body.Location).IsZero() {
			message.Body = s.locationMessage(body.Location)
		}
		if !reflect.ValueOf(body.Media).IsZero() {
			message.Body = s.mediaMessage(body.Media)
		}
		if !reflect.ValueOf(body.Ussd).IsZero() {
			message.Body = s.ussdMessage(body.Ussd)
		}
		if !reflect.ValueOf(body.VoiceActions).IsZero() {
			message.Body = s.voiceMessage(body.VoiceActions)
		}
		if !reflect.ValueOf(body.Email).IsZero() {
			message.Body = s.email(body.Email)
		}
	}
	s.replyChannel <- reply
}

func (s *service) On(event Event, handler NotificationHandler) {
	s.bus.SubscribeAsync(string(event), handler, false)
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
		msg, err := <-s.msgChannel, <-s.errChannel
		if errors.Is(err, io.EOF) {
			close(s.errChannel)
		}
		if err != nil {
			log.Println(err)
		}
		s.handleNotifications(msg)
	}
}
