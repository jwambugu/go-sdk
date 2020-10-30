package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// MessagingChannel is an enum
	MessagingChannel int32

	// MessagingConsentAction is an enum
	MessagingConsentAction int32

	// MessagingChannelNumber struct
	MessagingChannelNumber struct {
		Number  string           `json:"number"`
		Channel MessagingChannel `json:"channel"`
	}

	// Media defines the necessary attributes required to send a file as a message
	Media struct {
		URL  string         `json:"url"`
		Type hera.MediaType `json:"type"`
	}

	// Location defines a set of latitude and longitude that can be communicated as a message
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	// Template This refers to a predefined template for your message, the name of the template is used as the identifier and the params should be added in their logical order
	Template struct {
		Name   string   `json:"name"`
		Params []string `json:"params"`
	}

	// MessageBody defines how the message body should look like Note all the options are optional and the consturction of this struct depends on your needs.
	MessageBody struct {
		Text     string   `json:"text"`
		Media    Media    `json:"media"`
		Location Location `json:"location"`
		Template Template `json:"template"`
	}
)

const (
	MESSAGING_CHANNEL_UNSPECIFIED MessagingChannel = iota
	MESSAGING_CHANNEL_GOOGLE_RCS
	MESSAGING_CHANNEL_FB_MESSENGER
	MESSAGING_CHANNEL_SMS
	MESSAGING_CHANNEL_TELEGRAM
	MESSAGING_CHANNEL_WHATSAPP
)

const (
	MESSAGING_CONSENT_ACTION_UNSPECIFIED MessagingConsentAction = iota
	MESSAGING_CONSENT_ACTION_OPT_IN
	MESSAGING_CONSENT_ACTION_OPT_OUT
)

func (s *service) setMessageBodyAsText(text string) *hera.CustomerMessageBody_Text {
	return &hera.CustomerMessageBody_Text{
		Text: &hera.TextMessageBody{
			Text: &wrapperspb.StringValue{
				Value: text,
			},
		},
	}
}
func (s *service) setMessageBodyAsTemplate(template *Template) *hera.CustomerMessageBody_Text {
	return &hera.CustomerMessageBody_Text{
		Text: &hera.TextMessageBody{
			Template: &hera.TextMessageTemplate{
				Name:   template.Name,
				Params: template.Params,
			},
		},
	}
}
func (s *service) setMessageBodyAsLocation(location *Location) *hera.CustomerMessageBody_Location {
	return &hera.CustomerMessageBody_Location{
		Location: &hera.LocationMessageBody{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}
}
func (s *service) setMessageBodyAsMedia(media *Media) *hera.CustomerMessageBody_Media {
	return &hera.CustomerMessageBody_Media{
		Media: &hera.MediaMessageBody{
			Url:   media.URL,
			Media: media.Type,
		},
	}
}

func (s *service) SendMessage(
	customer *Customer,
	body *MessageBody,
	channelNumber *MessagingChannelNumber,
) (*hera.SendMessageReply, error) {
	var request hera.SendMessageRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.setCustomerNumber(customer)
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}

	if body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessage(ctx, &request)
}

func (s *service) SendMessageByTag(
	tag *Tag,
	body *MessageBody,
	channelNumber *MessagingChannelNumber,
) (*hera.TagCommandReply, error) {
	var request hera.SendMessageTagRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	if !reflect.ValueOf(tag).IsZero() {
		request.Tag = &hera.IndexMapping{
			Key: tag.Key,
			Value: &wrapperspb.StringValue{
				Value: tag.Value,
			},
		}
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}

	if body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessageByTag(ctx, &request)
}

func (s *service) ReplyToMessage(
	customer *Customer,
	messageId string,
	body *MessageBody,
) (*hera.SendMessageReply, error) {
	var request hera.ReplyToMessageRequest
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.CustomerId = customer.Id
	request.ReplyToMessageId = messageId

	if body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.ReplyToMessage(ctx, &request)
}

func (s *service) MessagingConsent(
	customer *Customer,
	channelNumber *MessagingChannelNumber,
	action MessagingConsentAction,
) (*hera.MessagingConsentReply, error) {
	var request hera.MessagingConsentRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.setCustomerNumber(customer)
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}
	request.Action = hera.MessagingConsentAction(action)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.MessagingConsent(ctx, &request)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(
	body *MessageBody,
	channelNumber *MessagingChannelNumber,
) (*hera.SendMessageReply, error) {
	return c.service.SendMessage(c, body, channelNumber)
}
