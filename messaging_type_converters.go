package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *service) messageBodyAsText(text string) *hera.CustomerMessageBody_Text {
	return &hera.CustomerMessageBody_Text{
		Text: &hera.TextMessageBody{
			Text: wrapperspb.String(text),
		},
	}
}

func (s *service) messageBodyAsTemplate(template *Template) *hera.CustomerMessageBody_Text {
	return &hera.CustomerMessageBody_Text{
		Text: &hera.TextMessageBody{
			Template: &hera.TextMessageTemplate{
				Name:   template.Name,
				Params: template.Params,
			},
		},
	}
}

func (s *service) messageBodyAsLocation(location *Location) *hera.CustomerMessageBody_Location {
	return &hera.CustomerMessageBody_Location{
		Location: &hera.LocationMessageBody{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}
}

func (s *service) messageBodyAsMedia(media *Media) *hera.CustomerMessageBody_Media {
	return &hera.CustomerMessageBody_Media{
		Media: &hera.MediaMessageBody{
			Url:   media.URL,
			Media: hera.MediaType(media.Type),
		},
	}
}

func (s *service) messageStatusNotf(notf *hera.MessageStatusNotification) *MessageStatusNotification {
	return &MessageStatusNotification{
		CustomerID: notf.CustomerId,
		MessageID:  notf.MessageId,
		Status:     MessageDeliveryStatus(notf.Status),
	}
}

func (s *service) messageSessionStatusNotf(notf *hera.MessagingSessionStatusNotification) *MessageSessionStatusNotification {
	return &MessageSessionStatusNotification{
		CustomerID: notf.CustomerId,
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		Expiration: notf.Expiration.Seconds,
		Status:     MessagingSessionStatus(notf.Status),
	}
}

func (s *service) messagingConsentStatusNotf(notf *hera.MessagingConsentStatusNotification) *MessagingConsentStatusNotification {
	return &MessagingConsentStatusNotification{
		CustomerID: notf.CustomerId,
		CustomerNumber: &CustomerNumber{
			Number:    notf.CustomerNumber.Number,
			Partition: notf.CustomerNumber.Partition.Value,
			Provider:  NumberProvider(notf.CustomerNumber.Provider),
		},
		ChannelNumber: &MessagingChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: MessagingChannel(notf.ChannelNumber.Channel),
		},
		Status: MessagingConsentStatus(notf.Status),
	}
}

func (s *service) recievedMessageNotification(notf *hera.ReceivedMessageNotification) *RecievedMessageNotification {
	var notification *RecievedMessageNotification

	notification.CustomerID = notf.CustomerId
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

	if !reflect.ValueOf(notf.Location).IsZero() {
		notification.Location = &Location{
			Latitude:  notf.Location.Latitude,
			Longitude: notf.Location.Longitude,
		}
	}

	var notificationMedia []*Media

	for _, notfMedia := range notf.Media {
		notificationMedia = append(notificationMedia, &Media{
			URL: notfMedia.Url,
		})
	}

	notification.Text = notf.Text.Value
	notification.MessageID = notf.MessageId
	notification.Media = notificationMedia
	return notification
}
