package elarian

import (
	"context"
	"errors"
	"io"
	"reflect"

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
	handler func(data interface{}, customer *Customer, service Service),
) error {
	err := s.bus.SubscribeAsync(string(event), handler, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveNotificationSubscriber(
	event ElarianNotification,
	handler func(data interface{}, customer *Customer, service Service),
) error {
	err := s.bus.Unsubscribe(string(event), handler)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) notificationsHandler(data *hera.WebhookRequest) {
	if ussdSession := data.GetUssdSession(); !reflect.ValueOf(ussdSession).IsZero() {
		newCustomer, _ := s.NewCustomer(&CreateCustomerParams{
			Id: ussdSession.CustomerId,
		})
		session := s.getUssdSessionNotification(ussdSession)
		s.bus.Publish(
			string(ELARIAN_USSD_SESSION_NOTIFICATION),
			session,
			newCustomer,
			s,
		)
	}
}

func (s *service) InitializeNotificationStream() error {
	var request hera.StreamNotificationRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	ctx := context.Background()
	stream, err := s.client.StreamNotifications(ctx, &request)
	errorChannel := make(chan error)
	if err != nil {
		return err
	}
	go func() {
		for {
			in, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(errorChannel)
				return

			}
			if err != nil {
				errorChannel <- err
				close(errorChannel)
				return
			}
			s.notificationsHandler(in)
		}
	}()
	err = <-errorChannel
	if err != nil {
		return err
	}
	return nil
}
