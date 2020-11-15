package elarian

import (
	"context"
	"errors"
	"io"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	ElarianNotification int32
)

const (
	ELARIAN_REMINDER_NOTIFICATION ElarianNotification = iota
	ELARIAN_VOICE_CALL_NOTIFICATION
	ELARIAN_USSD_SESSION_NOTIFICATION
	ELARIAN_PAYMENT_STATUS_NOTIFICATION
	ELARIAN_MESSAGE_STATUS_NOTIFICATION
	ELARIAN_RECIEVED_MESSAGE_NOTIFICATION
	ELARIAN_RECIEVED_PAYMENT_NOTIFICATION
	ELARIAN_WALLET_PAYMENT_STATUS_NOTIFICATION
	ELARIAN_MESSAGING_SESSION_STATUS_NOTIFICATION
	ELARIAN_MESSAGING_CONSENT_STATUS_NOTIFICATION
)

func (s *service) AddNotificationSubscriber(
	event ElarianNotification,
	handler func(svc Service, cust *Customer, data interface{}),
) error {
	err := s.bus.SubscribeAsync(string(event), handler, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveNotificationSubscriber(
	event ElarianNotification,
	handler func(svc Service, cust *Customer, data interface{}),
) error {
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
	request.AppId = s.appId
	request.OrgId = s.orgId
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
