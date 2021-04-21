package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *elarian) heraOutBoundTextMessage(text string) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Text{
			Text: text,
		},
	}
}

func (s *elarian) heraOutBoundTemplateMesage(template *Template) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Template{
			Template: &hera.TemplateMessageBody{
				Id:     template.ID,
				Params: template.Params,
			},
		},
	}
}

func (s *elarian) heraOutBoundMediaMessage(media *Media) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Media{
			Media: &hera.MediaMessageBody{
				Url:   media.URL,
				Media: hera.MediaType(media.Type),
			},
		},
	}
}

func (s *elarian) heraOutBoundLocationMessage(location *Location) *hera.OutboundMessageBody {
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

func (s *elarian) heraOutBoundUssdMessage(ussdMenu *UssdMenu) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Ussd{
			Ussd: &hera.UssdMenuMessageBody{
				Text:       ussdMenu.Text,
				IsTerminal: ussdMenu.IsTerminal,
			},
		},
	}
}

func (s *elarian) heraOutBoundVoiceMessage(actions []VoiceAction) *hera.OutboundMessageBody {
	return &hera.OutboundMessageBody{
		Entry: &hera.OutboundMessageBody_Voice{
			Voice: &hera.VoiceCallDialplanMessageBody{
				Actions: s.heraVoiceCallActions(actions),
			},
		},
	}
}

func (s *elarian) heraOutBoundEmail(email *Email) *hera.OutboundMessageBody {
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

func (s *elarian) OutboundMessage(message *hera.OutboundMessage) *OutBoundMessage {
	outboundMessage := &OutBoundMessage{}
	outboundMessage.Labels = message.Labels
	outboundMessage.ProviderTag = message.ProviderTag.Value
	outboundMessage.ReplyToken = message.ReplyToken.Value
	outboundMessage.ReplyPrompt = &OutboundMessageReplyPrompt{
		Action: PromptMessageReplyAction(message.ReplyPrompt.Action),
		Menu:   []*PromptMessageMenuItemBody{},
	}
	for _, menuItem := range message.ReplyPrompt.Menu {
		item := &PromptMessageMenuItemBody{}
		if entry, ok := menuItem.Entry.(*hera.PromptMessageMenuItemBody_Text); ok {
			item.Text = entry.Text
		}
		if entry, ok := menuItem.Entry.(*hera.PromptMessageMenuItemBody_Media); ok {
			item.Media = &Media{
				URL:  entry.Media.Url,
				Type: MediaType(entry.Media.Media),
			}
		}
		outboundMessage.ReplyPrompt.Menu = append(outboundMessage.ReplyPrompt.Menu, item)
	}

	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Text); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: TextMessage(entry.Text),
		}
		return outboundMessage
	}
	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Email); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: &Email{
				Subject:     entry.Email.Subject,
				Body:        entry.Email.BodyPlain,
				HTML:        entry.Email.BodyHtml,
				CcList:      entry.Email.CcList,
				BccList:     entry.Email.BccList,
				Attachments: entry.Email.Attachments,
			},
		}
		return outboundMessage
	}
	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Location); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: &Location{
				Latitude:  entry.Location.Latitude,
				Longitude: entry.Location.Longitude,
				Label:     entry.Location.Label.Value,
				Address:   entry.Location.Address.Value,
			},
		}
		return outboundMessage
	}
	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Media); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: &Media{
				URL:  entry.Media.Url,
				Type: MediaType(entry.Media.Media),
			},
		}
		return outboundMessage
	}
	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Template); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: &Template{
				ID:     entry.Template.Id,
				Params: entry.Template.Params,
			},
		}
		return outboundMessage
	}
	if entry, ok := message.Body.Entry.(*hera.OutboundMessageBody_Ussd); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: &UssdMenu{
				IsTerminal: entry.Ussd.IsTerminal,
				Text:       entry.Ussd.Text,
			},
		}
		return outboundMessage
	}
	if value, ok := message.Body.Entry.(*hera.OutboundMessageBody_Url); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: URLMessage(value.Url),
		}
		return outboundMessage

	}
	if value, ok := message.Body.Entry.(*hera.OutboundMessageBody_Voice); ok {
		outboundMessage.Body = &OutBoundMessageBody{
			Entry: s.voiceCallActions(value.Voice.Actions),
		}
		return outboundMessage
	}
	return outboundMessage
}

func (s *elarian) messageStatusNotf(notf *hera.MessageStatusNotification) *MessageStatusNotification {
	return &MessageStatusNotification{
		MessageID: notf.MessageId,
		Status:    MessageDeliveryStatus(notf.Status),
	}
}

func (s *elarian) messageSessionStartedNotf(notf *hera.MessagingSessionStartedNotification) *MessageSessionStartedNotification {
	return &MessageSessionStartedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: s.customerNumber(notf.CustomerNumber),
		Expiration:     notf.ExpiresAt.Seconds,
		SessionID:      notf.SessionId,
	}
}

func (s *elarian) messageSessionRenewedNotf(notf *hera.MessagingSessionRenewedNotification) *MessageSessionRenewedNotification {
	return &MessageSessionRenewedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: s.customerNumber(notf.CustomerNumber),
		Expiration:     notf.ExpiresAt.Seconds,
		SessionID:      notf.SessionId,
	}
}

func (s *elarian) MessageSessionEndedNotf(notf *hera.MessagingSessionEndedNotification) *MessageSessionEndedNotification {
	return &MessageSessionEndedNotification{
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: s.customerNumber(notf.CustomerNumber),
		Duration:       notf.Duration.AsDuration(),
		SessionID:      notf.SessionId,
		Reason:         MessagingSessionEndReason(notf.Reason),
	}
}

func (s *elarian) messagingConsentUpdateNotf(notf *hera.MessagingConsentUpdateNotification) *MessagingConsentUpdateNotification {
	return &MessagingConsentUpdateNotification{
		CustomerNumber: s.customerNumber(notf.CustomerNumber),
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		Status: MessagingConsentUpdateStatus(notf.Status),
	}
}

func (s *elarian) sentMessageReaction(notf *hera.SentMessageReactionNotification) *SentMessageReaction {
	return &SentMessageReaction{
		CustomerNumber: s.customerNumber(notf.CustomerNumber),
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		MessageID: notf.MessageId,
		Reaction:  MessageReaction(notf.Reaction),
	}
}

func (s *elarian) recievedMessageNotification(notf *hera.ReceivedMessageNotification) *RecievedMessageNotification {
	notification := &RecievedMessageNotification{}

	notification.MessageID = notf.MessageId

	if !reflect.ValueOf(notf.CustomerNumber).IsZero() {
		notification.CustomerNumber = s.customerNumber(notf.CustomerNumber)
	}

	if !reflect.ValueOf(notf.ChannelNumber).IsZero() {
		notification.ChannelNumber = &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		}
	}
	if notf.SessionId != nil {
		notification.SessionID = notf.SessionId.Value
	}
	if notf.InReplyTo != nil {
		notification.InReplyTo = notf.InReplyTo.Value
	}

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
						SessionID: notf.GetSessionId().Value,
						Input:     ussd.Ussd.GetValue(),
					},
				})
			continue
		}
		if voice, ok := part.Entry.(*hera.InboundMessageBody_Voice); ok {
			notification.Parts = append(notification.Parts, &InBoundMessageBody{
				Voice: s.voiceCallNotification(voice),
			})
		}
	}
	return notification
}
