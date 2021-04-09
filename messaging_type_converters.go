package elarian

import (
	"reflect"

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

func (s *service) messageStatusNotf(notf *hera.MessageStatusNotification) *MessageStatusNotification {
	return &MessageStatusNotification{
		MessageID: notf.MessageId,
		Status:    MessageDeliveryStatus(notf.Status),
	}
}

func (s *service) messageSessionStartedNotf(notf *hera.MessagingSessionStartedNotification) *MessageSessionStartedNotification {
	return &MessageSessionStartedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		Expiration: notf.ExpiresAt.Seconds,
		SessionID:  notf.SessionId,
	}
}

func (s *service) messageSessionRenewedNotf(notf *hera.MessagingSessionRenewedNotification) *MessageSessionRenewedNotification {
	return &MessageSessionRenewedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		Expiration: notf.ExpiresAt.Seconds,
		SessionID:  notf.SessionId,
	}
}

func (s *service) MessageSessionEndedNotf(notf *hera.MessagingSessionEndedNotification) *MessageSessionEndedNotification {
	return &MessageSessionEndedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		Duration:  notf.Duration.AsDuration(),
		SessionID: notf.SessionId,
		Reason:    MessagingSessionEndReason(notf.Reason),
	}
}

func (s *service) messagingConsentUpdateNotf(notf *hera.MessagingConsentUpdateNotification) *MessagingConsentUpdateNotification {
	return &MessagingConsentUpdateNotification{
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		Status: MessagingConsentUpdateStatus(notf.Status),
	}
}

func (s *service) sentMessageReaction(notf *hera.SentMessageReactionNotification) *SentMessageReaction {
	return &SentMessageReaction{
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		MessageID: notf.MessageId,
		Reaction:  MessageReaction(notf.Reaction),
	}
}

func (s *service) recievedMessageNotification(notf *hera.ReceivedMessageNotification) *RecievedMessageNotification {
	var notification *RecievedMessageNotification

	notification.MessageID = notf.MessageId

	if !reflect.ValueOf(notf.CustomerNumber).IsZero() {
		notification.CustomerNumber = &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		}
	}

	if !reflect.ValueOf(notf.ChannelNumber).IsZero() {
		notification.ChannelNumber = &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		}
	}
	notification.MessageID = notf.MessageId
	notification.SessionID = notf.SessionId.Value
	notification.InReplyTo = notf.InReplyTo.Value
	notification.Parts = []*InBoundMessageBody{}

	for _, part := range notf.Parts {
		if email, ok := part.Entry.(*hera.InboundMessageBody_Email); ok {
			notification.Parts = append(
				notification.Parts,
				&InBoundMessageBody{
					Email: &Email{
						Subject:     email.Email.Subject,
						Body:        email.Email.BodyPlain,
						HTML:        email.Email.BodyHtml,
						CcList:      email.Email.CcList,
						BccList:     email.Email.BccList,
						Attachments: email.Email.Attachments,
					},
				},
			)
			continue
		}

		if media, ok := part.Entry.(*hera.InboundMessageBody_Media); ok {
			notification.Parts = append(notification.Parts, &InBoundMessageBody{
				Media: &Media{
					URL:  media.Media.Url,
					Type: MediaType(media.Media.Media),
				},
			})
			continue
		}

		if location, ok := part.Entry.(*hera.InboundMessageBody_Location); ok {
			notification.Parts = append(notification.Parts, &InBoundMessageBody{
				Location: &Location{
					Latitude:  location.Location.Latitude,
					Longitude: location.Location.Longitude,
					Address:   location.Location.Address.Value,
					Label:     location.Location.Label.Value,
				},
			})
			continue
		}
		if text, ok := part.Entry.(*hera.InboundMessageBody_Text); ok {
			notification.Parts = append(notification.Parts, &InBoundMessageBody{Text: text.Text})
			continue

		}
		if ussd, ok := part.Entry.(*hera.InboundMessageBody_Ussd); ok {
			notification.Parts = append(notification.Parts,
				&InBoundMessageBody{
					Ussd: &UssdSessionNotification{
						SessionID: notf.SessionId.Value,
						Input:     ussd.Ussd.Value,
					},
				})

		}
		if voice, ok := part.Entry.(*hera.InboundMessageBody_Voice); ok {
			notification.Parts = append(notification.Parts, &InBoundMessageBody{
				Voice: s.voiceCallNotification(voice),
			})
		}
	}
	return notification
}
