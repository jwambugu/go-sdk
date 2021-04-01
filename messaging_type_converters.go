package elarian

import (
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *service) textMessage(text string) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Text{
			Text: text,
		},
	}
}

func (s *service) templateMesage(template *Template) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Template{
			Template: &hera.TemplateMessageBody{
				Id:     template.ID,
				Params: template.Params,
			},
		},
	}
}

func (s *service) mediaMessage(media *Media) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Media{
			Media: &hera.MediaMessageBody{
				Url:   media.URL,
				Media: hera.MediaType(media.Type),
			},
		},
	}
}

func (s *service) locationMessage(location *Location) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Location{
			Location: &hera.LocationMessageBody{
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
				Label:     wrapperspb.String(location.Label),
				Address:   wrapperspb.String(location.Address),
			},
		},
	}
}

func (s *service) ussdMessage(ussdMenu *UssdMenu) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Ussd{
			Ussd: &hera.UssdMenuMessageBody{
				Text:       ussdMenu.Text,
				IsTerminal: ussdMenu.IsTerminal,
			},
		},
	}
}

func (s *service) voiceMessage(actions []VoiceAction) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Voice{
			Voice: &hera.VoiceCallDialplanMessageBody{
				Actions: s.transformVoiceCallActions(actions),
			},
		},
	}
}

func (s *service) email(email *Email) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Email{
			Email: &hera.EmailMessageBody{
				Subject:     email.Subject,
				BodyPlain:   email.Body,
				BodyHtml:    email.HTML,
				CcList:      email.CcList,
				BccList:     email.BccList,
				Attachments: email.Attachments,
			},
		},
	}
}

// func (s *service) messageStatusNotf(notf *hera.MessageStatusNotification) *MessageStatusNotification {
// 	return &MessageStatusNotification{
// 		CustomerID: notf.CustomerId,
// 		MessageID:  notf.MessageId,
// 		Status:     MessageDeliveryStatus(notf.Status),
// 	}
// }

// func (s *service) messageSessionStatusNotf(notf *hera.MessagingSessionStatusNotification) *MessageSessionStatusNotification {
// 	return &MessageSessionStatusNotification{
// 		CustomerID: notf.CustomerId,
// 		ChannelNumber: &MessagingChannelNumber{
// 			Number:  notf.ChannelNumber.Number,
// 			Channel: MessagingChannel(notf.ChannelNumber.Channel),
// 		},
// 		CustomerNumber: &CustomerNumber{
// 			Number:    notf.CustomerNumber.Number,
// 			Partition: notf.CustomerNumber.Partition.Value,
// 			Provider:  NumberProvider(notf.CustomerNumber.Provider),
// 		},
// 		Expiration: notf.Expiration.Seconds,
// 		Status:     MessagingSessionStatus(notf.Status),
// 	}
// }

// func (s *service) messagingConsentStatusNotf(notf *hera.MessagingConsentStatusNotification) *MessagingConsentStatusNotification {
// 	return &MessagingConsentStatusNotification{
// 		CustomerID: notf.CustomerId,
// 		CustomerNumber: &CustomerNumber{
// 			Number:    notf.CustomerNumber.Number,
// 			Partition: notf.CustomerNumber.Partition.Value,
// 			Provider:  NumberProvider(notf.CustomerNumber.Provider),
// 		},
// 		ChannelNumber: &MessagingChannelNumber{
// 			Number:  notf.ChannelNumber.Number,
// 			Channel: MessagingChannel(notf.ChannelNumber.Channel),
// 		},
// 		Status: MessagingConsentStatus(notf.Status),
// 	}
// }

// func (s *service) recievedMessageNotification(notf *hera.ReceivedMessageNotification) *RecievedMessageNotification {
// 	var notification *RecievedMessageNotification

// 	notification.CustomerID = notf.CustomerId
// 	notification.MessageID = notf.MessageId

// 	if !reflect.ValueOf(notf.CustomerNumber).IsZero() {
// 		notification.CustomerNumber = &CustomerNumber{
// 			Number:    notf.CustomerNumber.Number,
// 			Partition: notf.CustomerNumber.Partition.Value,
// 			Provider:  NumberProvider(notf.CustomerNumber.Provider),
// 		}
// 	}

// 	if !reflect.ValueOf(notf.ChannelNumber).IsZero() {
// 		notification.ChannelNumber = &MessagingChannelNumber{
// 			Number:  notf.ChannelNumber.Number,
// 			Channel: MessagingChannel(notf.ChannelNumber.Channel),
// 		}
// 	}

// 	if !reflect.ValueOf(notf.Location).IsZero() {
// 		notification.Location = &Location{
// 			Latitude:  notf.Location.Latitude,
// 			Longitude: notf.Location.Longitude,
// 		}
// 	}

// 	var notificationMedia []*Media

// 	for _, notfMedia := range notf.Media {
// 		notificationMedia = append(notificationMedia, &Media{
// 			URL: notfMedia.Url,
// 		})
// 	}

// 	notification.Text = notf.Text.Value
// 	notification.MessageID = notf.MessageId
// 	notification.Media = notificationMedia
// 	return notification
// }
