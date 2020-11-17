package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) ussdSessionNotificationHandler(notf *hera.UssdSessionNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{ID: notf.CustomerId})
		sessionNotf := s.ussdSessionNotification(notf)
		s.bus.Publish(
			string(ElarianUssdSessionNotification),
			s,
			cust,
			sessionNotf,
		)
	}
}

func (s *service) reminderNotificationHandler(notf *hera.ReminderNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{ID: notf.CustomerId})
		reminderNotf := s.reminderNotification(notf)
		s.bus.Publish(
			string(ElarianReminderNotification),
			s,
			cust,
			reminderNotf,
		)
	}
}

func (s *service) paymentStatusNotificationHandler(notf *hera.PaymentStatusNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{ID: notf.CustomerId.Value})
		statusNotf := s.paymentStatusNotf(notf)
		s.bus.Publish(
			string(ElarianPaymentStatusNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) receivedPaymentNotificationHandler(notf *hera.ReceivedPaymentNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{
			ID: notf.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    notf.CustomerNumber.Number,
				Provider:  NumberProvider(notf.CustomerNumber.Provider),
				Partition: notf.CustomerNumber.Partition.Value,
			},
		})
		statusNotf := s.recievedPaymentNotf(notf)
		s.bus.Publish(
			string(ElarianRecievedPaymentNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) walletPaymentStatusNotificationHandler(notf *hera.WalletPaymentStatusNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{ID: notf.CustomerId})
		statusNotf := s.walletPaymentStatusNotf(notf)
		s.bus.Publish(
			string(ElarianWalletPaymentStatusNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) messageStatusNotificationHandler(notf *hera.MessageStatusNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{ID: notf.CustomerId})
		statusNotf := s.messageStatusNotf(notf)
		s.bus.Publish(
			string(ElarianMessageStatusNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) messagingSesssionStatusNotificationHandler(notf *hera.MessagingSessionStatusNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{
			ID: notf.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    notf.CustomerNumber.Number,
				Provider:  NumberProvider(notf.CustomerNumber.Provider),
				Partition: notf.CustomerNumber.Partition.Value,
			},
		})
		statusNotf := s.messageSessionStatusNotf(notf)
		s.bus.Publish(
			string(ElarianMessagingSessionStatusNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) messagingConsentStatusNotificationHandler(notf *hera.MessagingConsentStatusNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{
			ID: notf.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    notf.CustomerNumber.Number,
				Provider:  NumberProvider(notf.CustomerNumber.Provider),
				Partition: notf.CustomerNumber.Partition.Value,
			},
		})
		statusNotf := s.messagingConsentStatusNotf(notf)
		s.bus.Publish(
			string(ElarianMessagingConsentStatusNotification),
			s,
			cust,
			statusNotf,
		)
	}
}

func (s *service) recievedMessageNotificationHandler(notf *hera.ReceivedMessageNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{
			ID: notf.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    notf.CustomerNumber.Number,
				Provider:  NumberProvider(notf.CustomerNumber.Provider),
				Partition: notf.CustomerNumber.Partition.Value,
			},
		})
		statusNoft := s.recievedMessageNotification(notf)
		s.bus.Publish(
			string(ElarianRecievedMessageNotification),
			s,
			cust,
			statusNoft,
		)
	}
}

func (s *service) voiceCallNotificationHandler(notf *hera.VoiceCallNotification) {
	if !reflect.ValueOf(notf).IsZero() {
		cust := s.NewCustomer(&CreateCustomer{
			ID: notf.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    notf.CustomerNumber.Number,
				Provider:  NumberProvider(notf.CustomerNumber.Provider),
				Partition: notf.CustomerNumber.Partition.Value,
			},
		})
		statusNotf := s.voiceCallNotification(notf)
		s.bus.Publish(
			string(ElarianVoiceCallNotification),
			s,
			cust,
			statusNotf,
		)
	}
}
