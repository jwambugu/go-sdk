package elarian

import (
	"context"
	"errors"
	"io"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// Notification int
	Notification int32

	// NotificationHandler func
	NotificationHandler func(svc Service, cust *Customer, data interface{})
)

// Notification constants
const (
	ElarianReminderNotification Notification = iota
	ElarianVoiceCallNotification
	ElarianUssdSessionNotification
	ElarianPaymentStatusNotification
	ElarianMessageStatusNotification
	ElarianRecievedMessageNotification
	ElarianRecievedPaymentNotification
	ElarianWalletPaymentStatusNotification
	ElarianMessagingSessionStatusNotification
	ElarianMessagingConsentStatusNotification
)

func (s *service) AddNotificationSubscriber(event Notification, handler NotificationHandler) error {
	err := s.bus.SubscribeAsync(string(event), handler, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveNotificationSubscriber(event Notification, handler NotificationHandler) error {
	err := s.bus.Unsubscribe(string(event), handler)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) notificationsHandler(data *hera.WebhookRequest) {
	s.ussdSessionNotificationHandler(data.GetUssdSession())
	s.reminderNotificationHandler(data.GetReminder())
	s.paymentStatusNotificationHandler(data.GetPaymentStatus())
	s.receivedPaymentNotificationHandler(data.GetReceivedPayment())
	s.walletPaymentStatusNotificationHandler(data.GetWalletPaymentStatus())
	s.voiceCallNotificationHandler(data.GetVoiceCall())
	s.messageStatusNotificationHandler(data.GetMessageStatus())
	s.recievedMessageNotificationHandler(data.GetReceivedMessage())
	s.voiceCallNotificationHandler(data.GetVoiceCall())
	s.messagingSesssionStatusNotificationHandler(
		data.GetMessagingSessionStatus())
	s.messagingConsentStatusNotificationHandler(
		data.GetMessagingConsentStatus())

}

func (s *service) initializeNotificationStream() chan error {
	var request hera.StreamNotificationRequest
	request.AppId = s.appID
	request.OrgId = s.orgID
	errorChannel := make(chan error)

	ctx := context.Background()
	stream, err := s.client.StreamNotifications(ctx, &request)
	if err != nil {
		errorChannel <- err
	}
	go func() {
		for {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(errorChannel)
				return
			}
			if err != nil {
				errorChannel <- err
				return
			}
			s.notificationsHandler(data)
		}
	}()
	return errorChannel
}
