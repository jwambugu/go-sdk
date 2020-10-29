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

	// SendMessageRequest struct
	SendMessageRequest struct {
		AppId         string                 `json:"appId,omitempty"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
		Body          MessageBody            `json:"body"`
	}

	// SendMessageByTagRequest struct
	SendMessageByTagRequest struct {
		AppId         string                 `json:"appId,omitempty"`
		Tag           Tag                    `json:"tag"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
		Body          MessageBody            `json:"body"`
	}

	// ReplyToMessageRequest struct
	ReplyToMessageRequest struct {
		AppId            string      `json:"appId,omitempty"`
		ReplyToMessageId string      `json:"customerId,omitempty"`
		Body             MessageBody `json:"body"`
	}

	// MessagingConsentRequest struct
	MessagingConsentRequest struct {
		AppId         string                 `json:"appId,omitempty"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
		Action        MessagingConsentAction `json:"consentAction"`
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

func (s *service) SendMessage(customer *Customer, params *SendMessageRequest) (*hera.SendMessageReply, error) {
	var request hera.SendMessageRequest
	request.AppId = params.AppId
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.setCustomerNumber(customer)
	}
	if !reflect.ValueOf(params.ChannelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(params.ChannelNumber.Channel),
			Number:  params.ChannelNumber.Number,
		}
	}

	if params.Body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(params.Body.Text),
		}
	}
	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&params.Body.Template),
		}
	}
	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&params.Body.Location),
		}
	}
	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&params.Body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessage(ctx, &request)
}

func (s *service) SendMessageByTag(
	params *SendMessageByTagRequest,
) (*hera.TagCommandReply, error) {
	var request hera.SendMessageTagRequest
	request.AppId = params.AppId
	request.OrgId = s.orgId

	if !reflect.ValueOf(params.Tag).IsZero() {
		request.Tag = &hera.IndexMapping{
			Key: params.Tag.Key,
			Value: &wrapperspb.StringValue{
				Value: params.Tag.Value,
			},
		}
	}

	if params.Body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(params.Body.Text),
		}
	}
	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&params.Body.Template),
		}
	}
	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&params.Body.Location),
		}
	}
	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&params.Body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessageByTag(ctx, &request)
}

func (s *service) ReplyToMessage(customer *Customer, params *ReplyToMessageRequest) (*hera.SendMessageReply, error) {
	var request hera.ReplyToMessageRequest
	request.AppId = params.AppId
	request.OrgId = s.orgId
	request.CustomerId = customer.Id
	request.ReplyToMessageId = params.ReplyToMessageId

	if params.Body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsText(params.Body.Text),
		}
	}
	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsTemplate(&params.Body.Template),
		}
	}
	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsLocation(&params.Body.Location),
		}
	}
	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.setMessageBodyAsMedia(&params.Body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.ReplyToMessage(ctx, &request)
}

func (s *service) MessagingConsent(customer *Customer, params *MessagingConsentRequest) (*hera.MessagingConsentReply, error) {
	var request hera.MessagingConsentRequest

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.setCustomerNumber(customer)
	}
	if !reflect.ValueOf(params.ChannelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(params.ChannelNumber.Channel),
			Number:  params.ChannelNumber.Number,
		}
	}
	request.Action = hera.MessagingConsentAction(params.Action)
	request.AppId = params.AppId
	request.OrgId = s.orgId

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.MessagingConsent(ctx, &request)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(
	params *SendMessageRequest,
) (*hera.SendMessageReply, error) {
	return c.service.SendMessage(c, params)
}
