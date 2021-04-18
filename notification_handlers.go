package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *elarian) reminderNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_Reminder); ok {
		reminder := s.reminderNotification(entry.Reminder)
		customer := &Customer{ID: notf.CustomerId, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianReminderNotification), s, reminder, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) messageStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessageStatus); ok {
		statusNotification := s.messageStatusNotf(entry.MessageStatus)
		customer := &Customer{ID: notf.CustomerId, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianMessageStatusNotification), s, statusNotification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) messagingSessionStartedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionStarted); ok {
		notification := s.messageSessionStartedNotf(entry.MessagingSessionStarted)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianMessagingSessionStartedNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) messagingSessionRenewedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionRenewed); ok {
		notification := s.messageSessionRenewedNotf(entry.MessagingSessionRenewed)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianMessagingSessionRenewedNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) messagingSessionEndedNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingSessionEnded); ok {
		notification := s.MessageSessionEndedNotf(entry.MessagingSessionEnded)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianMessagingSessionEndedNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) messagingConsentUpdateNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_MessagingConsentUpdate); ok {
		notification := s.messagingConsentUpdateNotf(entry.MessagingConsentUpdate)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianMessagingConsentUpdateNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) recievedMessageNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_ReceivedMessage); ok {
		notification := s.recievedMessageNotification(entry.ReceivedMessage)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}

		for _, part := range notification.Parts {
			if notification.ChannelNumber.Channel == MessagingChannelUssd {
				s.bus.Publish(string(ElarianReceivedUssdSessionNotification), s, part.Ussd, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelEmail {
				s.bus.Publish(string(ElarianReceivedEmailNotification), s, part.Email, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelVoice {
				s.bus.Publish(string(ElarianReceivedVoiceCallNotification), s, part.Voice, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelSms {
				s.bus.Publish(string(ElarianReceivedSmsNotification), s, part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelTelegram {
				s.bus.Publish(string(ElarianReceivedTelegramNotification), s, part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelWhatsapp {
				s.bus.Publish(string(ElarianReceivedWhatsappNotification), s, part, appData, customer, s.notificationCallBack)
			}
			if notification.ChannelNumber.Channel == MessagingChannelFBMessanger {
				s.bus.Publish(string(ElarianReceivedFbMessengerNotification), s, part, appData, customer, s.notificationCallBack)
			}
		}
	}
}

func (s *elarian) sentMesssageNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_SentMessageReaction); ok {
		notification := s.sentMessageReaction(entry.SentMessageReaction)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianSentMessageReactionNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) receivedPaymentNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_ReceivedPayment); ok {
		notification := s.recievedPaymentNotf(entry.ReceivedPayment)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianReceivedPaymentNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) paymentStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_PaymentStatus); ok {
		notification := s.paymentStatusNotf(entry.PaymentStatus)
		customer := &Customer{ID: notf.CustomerId, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianPaymentStatusNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) walletPaymentStatusNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_WalletPaymentStatus); ok {
		notification := s.walletPaymentStatusNotf(entry.WalletPaymentStatus)
		customer := &Customer{ID: notf.CustomerId, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianWalletPaymentStatusNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) customerActivityNotificationHandler(notf *hera.ServerToAppCustomerNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToAppCustomerNotification_CustomerActivity); ok {
		notification := s.customerActivity(entry.CustomerActivity)
		customer := &Customer{ID: notf.CustomerId, CustomerNumber: notification.CustomerNumber, service: s}
		appData := &Appdata{}
		if !reflect.ValueOf(notf.AppData).IsZero() {
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_StringVal); ok {
				appData.Value = val.StringVal
			}
			if val, ok := notf.AppData.Value.(*hera.DataMapValue_BytesVal); ok {
				appData.BytesValue = val.BytesVal
			}
		}
		s.bus.Publish(string(ElarianCustomerActivityNotification), s, notification, appData, customer, s.notificationCallBack)
	}
}

func (s *elarian) paymentPurseStatusNotificationHandler(notf *hera.ServerToAppNotification_Purse) {
	if reflect.ValueOf(notf).IsZero() || reflect.ValueOf(notf.Purse).IsZero() {
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
			Status:        PaymentStatus(entry.PaymentStatus.Status),
		}
		s.bus.Publish(string(ElarianPaymentPurseNotifiication), s, notification, nil, nil, s.notificationCallBack)
	}
}

func (s *elarian) SendChannelPaymentSimulatorNotificationHandler(notf *hera.ServerToSimulatorNotification) {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
		return
	}
	if entry, ok := notf.Entry.(*hera.ServerToSimulatorNotification_SendChannelPayment); ok {
		notification := &SendChannelPaymentSimulatorNotification{
			OrgID:         entry.SendChannelPayment.OrgId,
			AppID:         entry.SendChannelPayment.AppId,
			TransactionID: entry.SendChannelPayment.TransactionId,
			Account:       entry.SendChannelPayment.Account.Value,
			Value: &Cash{
				CurrencyCode: entry.SendChannelPayment.Value.CurrencyCode,
				Amount:       entry.SendChannelPayment.Value.Amount,
			},
			ChannelNumber: &PaymentChannelNumber{
				Number:  entry.SendChannelPayment.ChannelNumber.Number,
				Channel: PaymentChannel(entry.SendChannelPayment.ChannelNumber.Channel),
			},
		}
		notification.DebitParty = &PaymentParty{}
		if purse, ok := entry.SendChannelPayment.DebitParty.(*hera.SendChannelPaymentSimulatorNotification_Purse); ok {
			notification.DebitParty.Purse = &Purse{
				PurseID: purse.Purse.PurseId,
			}
		}
		if wallet, ok := entry.SendChannelPayment.DebitParty.(*hera.SendChannelPaymentSimulatorNotification_Wallet); ok {
			notification.DebitParty.Wallet = &Wallet{
				WalletID:   wallet.Wallet.WalletId,
				CustomerID: wallet.Wallet.CustomerId,
			}
		}
		s.bus.Publish(string(ElarianSendChannelPaymentSimulatorNotification), s, notification, nil, nil, s.notificationCallBack)
	}
}

func (s *elarian) CheckoutPaymentSimulatorNotificationHandler(notf *hera.ServerToSimulatorNotification) {
	if entry, ok := notf.Entry.(*hera.ServerToSimulatorNotification_CheckoutPayment); ok {
		customer := &Customer{
			ID: entry.CheckoutPayment.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    entry.CheckoutPayment.CustomerNumber.Number,
				Provider:  NumberProvider(entry.CheckoutPayment.CustomerNumber.Provider),
				Partition: entry.CheckoutPayment.CustomerNumber.Partition.Value,
			},
			service: s,
		}
		notification := &CheckoutPaymentSimulatorNotification{
			OrgID:         entry.CheckoutPayment.OrgId,
			AppID:         entry.CheckoutPayment.AppId,
			TransactionID: entry.CheckoutPayment.TransactionId,
			Value: &Cash{
				CurrencyCode: entry.CheckoutPayment.Value.CurrencyCode,
				Amount:       entry.CheckoutPayment.Value.Amount,
			},
			ChannelNumber: &PaymentChannelNumber{
				Number:  entry.CheckoutPayment.ChannelNumber.Number,
				Channel: PaymentChannel(entry.CheckoutPayment.ChannelNumber.Channel),
			},
		}
		notification.CreditParty = &PaymentParty{}
		if purse, ok := entry.CheckoutPayment.CreditParty.(*hera.CheckoutPaymentSimulatorNotification_Purse); ok {
			notification.CreditParty.Purse = &Purse{
				PurseID: purse.Purse.PurseId,
			}
		}
		if wallet, ok := entry.CheckoutPayment.CreditParty.(*hera.CheckoutPaymentSimulatorNotification_Wallet); ok {
			notification.CreditParty.Wallet = &Wallet{
				CustomerID: wallet.Wallet.CustomerId,
				WalletID:   wallet.Wallet.WalletId,
			}
		}
		s.bus.Publish(string(ElarianCheckoutPaymentSimulatorNotification), s, notification, nil, customer, s.notificationCallBack)
	}
}
func (s *elarian) SendCustomerPaymentSimulatorNotificationHandler(notf *hera.ServerToSimulatorNotification) {
	if entry, ok := notf.Entry.(*hera.ServerToSimulatorNotification_SendCustomerPayment); ok {
		customer := &Customer{
			ID: entry.SendCustomerPayment.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    entry.SendCustomerPayment.CustomerNumber.Number,
				Provider:  NumberProvider(entry.SendCustomerPayment.CustomerNumber.Provider),
				Partition: entry.SendCustomerPayment.CustomerNumber.Partition.Value,
			},
			service: s,
		}
		notification := &SendCustomerPaymentSimulatorNotification{
			OrgID:         entry.SendCustomerPayment.OrgId,
			AppID:         entry.SendCustomerPayment.AppId,
			TransactionID: entry.SendCustomerPayment.TransactionId,
			ChannelNumber: &PaymentChannelNumber{
				Number:  entry.SendCustomerPayment.ChannelNumber.Number,
				Channel: PaymentChannel(entry.SendCustomerPayment.ChannelNumber.Channel),
			},
			Value: &Cash{
				CurrencyCode: entry.SendCustomerPayment.Value.CurrencyCode,
				Amount:       entry.SendCustomerPayment.Value.Amount,
			},
		}
		notification.DebitParty = &PaymentParty{}
		if purse, ok := entry.SendCustomerPayment.DebitParty.(*hera.SendCustomerPaymentSimulatorNotification_Purse); ok {
			notification.DebitParty.Purse = &Purse{
				PurseID: purse.Purse.PurseId,
			}
		}
		if wallet, ok := entry.SendCustomerPayment.DebitParty.(*hera.SendCustomerPaymentSimulatorNotification_Wallet); ok {
			notification.DebitParty.Wallet = &Wallet{
				CustomerID: wallet.Wallet.CustomerId,
				WalletID:   wallet.Wallet.WalletId,
			}
		}
		s.bus.Publish(string(ElarianSendCustomerPaymentSimulatorNotification), s, notification, nil, customer, s.notificationCallBack)
	}
}
func (s *elarian) MakeVoiceCallSimulatorNotificationHandler(notf *hera.ServerToSimulatorNotification) {
	if entry, ok := notf.Entry.(*hera.ServerToSimulatorNotification_MakeVoiceCall); ok {
		customer := &Customer{
			ID: entry.MakeVoiceCall.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    entry.MakeVoiceCall.CustomerNumber.Number,
				Provider:  NumberProvider(entry.MakeVoiceCall.CustomerNumber.Provider),
				Partition: entry.MakeVoiceCall.CustomerNumber.Partition.Value,
			},
			service: s,
		}
		notification := &MakeVoiceCallSimulatorNotification{
			OrgID:     entry.MakeVoiceCall.OrgId,
			SessionID: entry.MakeVoiceCall.SessionId,
			ChannelNumber: &MessagingChannelNumber{
				Number:  entry.MakeVoiceCall.ChannelNumber.Number,
				Channel: MessagingChannel(entry.MakeVoiceCall.ChannelNumber.Channel),
			},
		}
		s.bus.Publish(string(ElarianMakeVoiceCallSimulatorNotification), s, notification, nil, customer, s.notificationCallBack)
	}
}

func (s *elarian) SendMessageSimulatorNotificationHandler(notf *hera.ServerToSimulatorNotification) {
	if entry, ok := notf.Entry.(*hera.ServerToSimulatorNotification_SendMessage); ok {
		customer := &Customer{
			ID: entry.SendMessage.CustomerId,
			CustomerNumber: &CustomerNumber{
				Number:    entry.SendMessage.CustomerNumber.Number,
				Provider:  NumberProvider(entry.SendMessage.CustomerNumber.Provider),
				Partition: entry.SendMessage.CustomerNumber.Partition.Value,
			},
			service: s,
		}
		notification := &SendMessageSimulatorNotification{
			OrgID:     entry.SendMessage.OrgId,
			MessageID: entry.SendMessage.MessageId,
			ChannelNumber: &MessagingChannelNumber{
				Number:  entry.SendMessage.ChannelNumber.Number,
				Channel: MessagingChannel(entry.SendMessage.ChannelNumber.Channel),
			},
		}
		notification.Message = s.OutboundMessage(entry.SendMessage.Message)
		s.bus.Publish(string(ElarianSendMessageSimulatorNotification), s, notification, nil, customer, s.notificationCallBack)
	}
}
