package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) reminderNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if reminderNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_Reminder); ok {
		reminder := s.reminderNotification(reminderNotf.Reminder)
		customer := &Customer{ID: notf.CustomerId}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianReminderNotification), reminder, appData, customer, s.notificationCallBack)
	}
}

func (s *service) messageStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if messageStatusNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessageStatus); ok {
		statusNotification := s.messageStatusNotf(messageStatusNotf.MessageStatus)
		customer := &Customer{ID: notf.CustomerId}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianMessageStatusNotification), statusNotification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) messagingSessionStartedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if messagingSessionStartedNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionStarted); ok {
		notification := s.messageSessionStartedNotf(messagingSessionStartedNotf.MessagingSessionStarted)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianMessagingSessionStartedNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) messagingSessionRenewedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if messagingSessionRenewedNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionRenewed); ok {
		notification := s.messageSessionRenewedNotf(messagingSessionRenewedNotf.MessagingSessionRenewed)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianMessagingSessionRenewedNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) messagingSessionEndedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if messagingSessionEndedNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionEnded); ok {
		notification := s.MessageSessionEndedNotf(messagingSessionEndedNotf.MessagingSessionEnded)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianMessagingSessionEndedNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) messagingConsentUpdateNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if messagingConsentUpdateNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingConsentUpdate); ok {
		notification := s.messagingConsentUpdateNotf(messagingConsentUpdateNotf.MessagingConsentUpdate)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianMessagingConsentUpdateNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) recievedMessageNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if recievedMessageNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_ReceivedMessage); ok {
		notification := s.recievedMessageNotification(recievedMessageNotf.ReceivedMessage)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		for _, part := range notification.Parts {
			if notification.ChannelNumber.Channel == MessagingChannelUssd {
				s.bus.Publish(string(ElarianReceivedUssdSessionNotification), part.Ussd, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelEmail {
				s.bus.Publish(string(ElarianReceivedEmailNotification), part.Email, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelVoice {
				s.bus.Publish(string(ElarianReceivedVoiceCallNotification), part.Voice, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelSms {
				s.bus.Publish(string(ElarianReceivedSmsNotification), part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelTelegram {
				s.bus.Publish(string(ElarianReceivedTelegramNotification), part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelWhatsapp {
				s.bus.Publish(string(ElarianReceivedWhatsappNotification), part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelFBMessanger {
				s.bus.Publish(string(ElarianReceivedFbMessengerNotification), part, appData, customer, s.notificationCallBack)
			}
		}
	}
}

func (s *service) sentMesssageNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if sentMessageReactionNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_SentMessageReaction); ok {
		notification := s.sentMessageReaction(sentMessageReactionNotf.SentMessageReaction)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianSentMessageReactionNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) receivedPaymentNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if receivedPaymentNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_ReceivedPayment); ok {
		notification := s.recievedPaymentNotf(receivedPaymentNotf.ReceivedPayment)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianReceivedPaymentNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) paymentStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if paymentStatusNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_PaymentStatus); ok {
		notification := s.paymentStatusNotf(paymentStatusNotf.PaymentStatus)
		customer := &Customer{ID: notf.CustomerId}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianPaymentStatusNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) walletPaymentStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if walletPaymentStatusNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_WalletPaymentStatus); ok {
		notification := s.walletPaymentStatusNotf(walletPaymentStatusNotf.WalletPaymentStatus)
		customer := &Customer{ID: notf.CustomerId}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianWalletPaymentStatusNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) customerActivityNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if customerActivityNotf, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_CustomerActivity); ok {
		notification := s.customerActivity(customerActivityNotf.CustomerActivity)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			appData.Value = notf.AppData.GetStringVal()
			appData.BytesValue = notf.AppData.GetBytesVal()
		}
		s.bus.Publish(string(ElarianCustomerActivityNotification), notification, appData, customer, s.notificationCallBack)
	}
}

func (s *service) paymentPurseStatusNotificationHandler(notf *hera.ServerToAppNotification_Purse) {
	if reflect.ValueOf(notf).IsZero() {
		return
	}
	if reflect.ValueOf(notf.Purse).IsZero() {
		return
	}
	notification := &PurseNotification{
		OrgID:     notf.Purse.OrgId,
		AppID:     notf.Purse.AppId,
		CreatedAt: notf.Purse.CreatedAt.AsTime(),
		PurseID:   notf.Purse.PurseId,
	}
	if entry, ok := notf.Purse.Entry.(*hera.ServerToAppPurseNotification_PaymentStatus); ok {
		notification.PaymentStatus = &NotificationPaymentStatus{
			TransactionID: entry.PaymentStatus.TransactionId,
			Status:        int(entry.PaymentStatus.Status),
		}
		s.bus.Publish(string(ElarianPaymentPurseNotifiication), notification, nil, nil, s.notificationCallBack)
	}
}
