package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// USSDChannel type
	USSDChannel int32

	// USSDMenu struct
	USSDMenu struct {
		IsTerminal bool   `json:"isTerminal,omitempty"`
		Text       string `json:"text,omitempty"`
	}
	// USSDRequest struct
	USSDRequest struct {
		AppId     string   `json:"appId,omitempty"`
		SessionId string   `json:"sessionId,omitempty"`
		USSDMenu  USSDMenu `json:"ussdMenu,omitempty"`
	}

	// USSDChannelNumber struct
	USSDChannelNumber struct {
		Channel USSDChannel `json:"channel,omitempty"`
		Number  string      `json:"number,omitempty"`
	}

	// USSDSessionNotification struct
	USSDSessionNotification struct {
		SessionId      string            `json:"sessionId,omitempty"`
		CustomerId     string            `json:"customerId,omitempty"`
		Input          string            `json:"input,omitempty"`
		CustomerNumber *CustomerNumber   `json:"customerNumber,omitempty"`
		ChannelNumber  USSDChannelNumber `json:"channelNumber,omitempty"`
	}
)

const (
	USSD_CHANNEL_UNSPECIFIED USSDChannel = iota
	USSD_CHANNEL_TELCO
)

func (s *service) ReplyToUSSDSession(params *USSDRequest) (*hera.WebhookResponseReply, error) {
	var request hera.WebhookResponse
	request.AppId = params.AppId
	request.OrgId = s.orgId
	request.SessionId = params.SessionId
	request.UssdMenu = &hera.UssdMenu{
		IsTerminal: params.USSDMenu.IsTerminal,
		Text:       params.USSDMenu.Text,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendWebhookResponse(ctx, &request)
}

// func (s *service) nott() {
// 	var reminderNotification hera.ReminderNotification
// 	reminderNotification.CustomerId = ""
// 	reminderNotification.WorkId = wrapperspb.String("")
// 	reminderNotification.Reminder = &hera.CustomerReminder{
// 		AppId:      "",
// 		Expiration: timestamppb.Now(),
// 		Interval:   durationpb.New(2),
// 		Key:        "",
// 		Payload:    wrapperspb.String(""),
// 	}
// 	reminderNotification.Tag = &hera.CustomerIndex{
// 		Expiration: timestamppb.Now(),
// 		Mapping: &hera.IndexMapping{
// 			Key:   "",
// 			Value: wrapperspb.String(""),
// 		},
// 	}

// 	var voiceCallNotification hera.VoiceCallNotification
// 	voiceCallNotification.ChannelNumber = &hera.VoiceChannelNumber{
// 		Channel: hera.VoiceChannel(1),
// 		Number:  "",
// 	}
// 	voiceCallNotification.Cost = &hera.Cash{
// 		Amount:       float64(303.00),
// 		CurrencyCode: "",
// 	}
// 	voiceCallNotification.CustomerId = ""
// 	voiceCallNotification.CustomerNumber = &hera.CustomerNumber{}
// 	voiceCallNotification.Direction = hera.CustomerEventDirection(1)
// 	voiceCallNotification.Duration = durationpb.New(2)
// 	voiceCallNotification.SessionId = ""
// 	voiceCallNotification.Input = &hera.VoiceCallHopInput{
// 		DialData: &hera.VoiceCallDialInput{
// 			DestinationNumber: "",
// 			Duration:          durationpb.New(2),
// 			StartedAt:         timestamppb.Now(),
// 		},
// 		DtmfDigits: &wrapperspb.StringValue{
// 			Value: "",
// 		},
// 		HangupCause: hera.VoiceCallHangupCause(2),
// 		QueueData: &hera.VoiceCallQueueInput{
// 			DequeuedAt:          timestamppb.Now(),
// 			DequeuedToNumber:    wrapperspb.String("d"),
// 			DequeuedToSessionId: wrapperspb.String(""),
// 			EnqueuedAt:          timestamppb.Now(),
// 			QueueDuration:       durationpb.New(2),
// 		},
// 		RecordingUrl: wrapperspb.String(""),
// 		StartedAt:    timestamppb.Now(),
// 		Status:       hera.VoiceCallStatus(2),
// 	}

// 	var paymentStatusNotification hera.PaymentStatusNotification
// 	paymentStatusNotification.CustomerId = ""
// 	paymentStatusNotification.Status = hera.PaymentStatus(1)
// 	paymentStatusNotification.TransactionId = ""

// 	var messageStatusNotification hera.MessageStatusNotification
// 	messageStatusNotification.CustomerId = ""
// 	messageStatusNotification.MessageId = ""
// 	messageStatusNotification.Status = hera.MessageDeliveryStatus(1)

// 	var recievedMessageNotification hera.ReceivedMessageNotification
// 	recievedMessageNotification.ChannelNumber = &hera.MessagingChannelNumber{}
// 	recievedMessageNotification.CustomerId = ""
// 	recievedMessageNotification.CustomerNumber = &hera.CustomerNumber{}
// 	recievedMessageNotification.Location = &hera.LocationMessageBody{}
// 	recievedMessageNotification.Media = []*hera.MediaMessageBody{
// 		&hera.MediaMessageBody{
// 			Media: hera.MediaType(2),
// 			Url:   "",
// 		},
// 	}
// 	recievedMessageNotification.MessageId = ""
// 	recievedMessageNotification.Text = wrapperspb.String("")

// 	var recievedPaymentNotification hera.ReceivedPaymentNotification

// 	var walletPaymentStatusNotification hera.WalletPaymentStatusNotification
// 	var messageSession hera.MessagingSessionStatusNotification
// 	var messageConsentStatus hera.MessagingConsentStatusNotification
// }
