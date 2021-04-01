package elarian

import (
	"context"
	"log"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

// type (
// 	// Notification int
// 	Notification int32

// 	// NotificationHandler func
// 	NotificationHandler func(svc Service, cust *Customer, data interface{})
// )

// // Notification constants
// const (
// 	ElarianReminderNotification Notification = iota
// 	ElarianVoiceCallNotification
// 	ElarianUssdSessionNotification
// 	ElarianPaymentStatusNotification
// 	ElarianMessageStatusNotification
// 	ElarianRecievedMessageNotification
// 	ElarianRecievedPaymentNotification
// 	ElarianWalletPaymentStatusNotification
// 	ElarianMessagingSessionStatusNotification
// 	ElarianMessagingConsentStatusNotification
// )

// func (s *service) AddNotificationSubscriber(event Notification, handler NotificationHandler) error {
// 	err := s.bus.SubscribeAsync(string(event), handler, false)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *service) RemoveNotificationSubscriber(event Notification, handler NotificationHandler) error {
// 	err := s.bus.Unsubscribe(string(event), handler)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // func (s *service) notificationsHandler(data *hera.ServerToAppNotificationReply) {
// // 	s.ussdSessionNotificationHandler(data.GetUssdSession())
// // 	s.reminderNotificationHandler(data.GetReminder())
// // 	s.paymentStatusNotificationHandler(data.GetPaymentStatus())
// // 	s.receivedPaymentNotificationHandler(data.GetReceivedPayment())
// // 	s.walletPaymentStatusNotificationHandler(data.GetWalletPaymentStatus())
// // 	s.voiceCallNotificationHandler(data.GetVoiceCall())
// // 	s.messageStatusNotificationHandler(data.GetMessageStatus())
// // 	s.recievedMessageNotificationHandler(data.GetReceivedMessage())
// // 	s.voiceCallNotificationHandler(data.GetVoiceCall())
// // 	s.messagingSesssionStatusNotificationHandler(
// // 		data.GetMessagingSessionStatus())
// // 	s.messagingConsentStatusNotificationHandler(
// // 		data.GetMessagingConsentStatus())

// // }

func (s *service) notificationHandler() {
	req := new(hera.ServerToAppNotification)

	data, err := proto.Marshal(req)
	if err != nil {
	}

	s.client.RequestResponse(payload.New(data, []byte{})).
		Subscribe(
			context.Background(),
			rx.OnNext(func(input payload.Payload) error {
				log.Println("SOME PAYLOAD", input)
				return nil
			}),
			rx.OnComplete(func() {
				//
			}),
			rx.OnError(func(e error) {
				log.Println("ERROR: ", err)
			}),
		)
}
