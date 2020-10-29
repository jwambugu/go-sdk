package elarian

import (
	"context"
	"errors"
	"io"
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	ElarianEvent int32
)

const (
	ElarianReminderEvent ElarianEvent = iota
	ElarianVoiceCallEvent
	ElarianUSSDSessionEvent
	ElarianPaymentStatusEvent
	ElarianMessageStatusEvent
	ElarianReceivedMessageEvent
	ElarianReceivedPaymentEvent
	ElarianWalletPaymentStatusEvent
	ElarianMessagingSessionStatusEvent
	ElarianMessagingConsentStatusEvent
)

func (s *service) AddNotificationSubscriber(
	event ElarianEvent,
	handler func(data interface{}, customer *Customer),
) error {
	err := s.bus.SubscribeAsync(string(event), handler, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveNotificationSubscriber(
	event ElarianEvent,
	handler func(data interface{}, customer *Customer),
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
		session := &USSDSessionNotification{
			SessionId:  ussdSession.SessionId,
			CustomerId: ussdSession.CustomerId,
			Input:      ussdSession.Input.Value,
			CustomerNumber: &CustomerNumber{
				Number:    ussdSession.CustomerNumber.Number,
				Partition: ussdSession.CustomerNumber.Partition.Value,
				Provider:  NumberProvider(ussdSession.CustomerNumber.Provider),
			},
			ChannelNumber: USSDChannelNumber{
				Channel: USSDChannel(ussdSession.ChannelNumber.Channel),
				Number:  ussdSession.ChannelNumber.Number,
			},
		}
		s.bus.Publish(string(ElarianUSSDSessionEvent), session, newCustomer)
	}
}

func (s *service) InitializeNotificationStream(appId string) error {
	var request hera.StreamNotificationRequest
	request.AppId = appId
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
