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
		AppID         string                 `json:"appId,omitempty"`
		ProductID     string                 `json:"productId,omitempty"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
		Body          MessageBody            `json:"body"`
	}

	// SendMessageByTagRequest struct
	SendMessageByTagRequest struct {
		AppID         string                 `json:"appId,omitempty"`
		ProductID     string                 `json:"productId,omitempty"`
		Tag           Tag                    `json:"tag"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
		Body          MessageBody            `json:"body"`
	}

	// ReplyToMessageRequest struct
	ReplyToMessageRequest struct {
		AppID            string      `json:"appId,omitempty"`
		ProductID        string      `json:"productId,omitempty"`
		ReplyToMessageID string      `json:"customerId,omitempty"`
		Body             MessageBody `json:"body"`
	}

	// MessagingConsentRequest struct
	MessagingConsentRequest struct {
		AppID         string                 `json:"appId,omitempty"`
		ProductID     string                 `json:"productId,omitempty"`
		ChannelNumber MessagingChannelNumber `json:"channelNumber"`
	}
)

const (
	// MessagingChannelUnspecified is a type of MessagingChannel
	MessagingChannelUnspecified MessagingChannel = iota
	// MessagingChannelGoogleRCS is a type of MessagingChannel
	MessagingChannelGoogleRCS
	// MessagingChannelFaceBookMessanger is a type of MessagingChannel
	MessagingChannelFaceBookMessanger
	// MessagingChannelSMS is a type of MessagingChannel
	MessagingChannelSMS
	// MessagingChannelTelegram is a type of MessagingChannel
	MessagingChannelTelegram
	// MessagingChannelWhatsapp is a type of MessagingChannel
	MessagingChannelWhatsapp
)

func (s *service) SendMessage(customer *Customer, params *SendMessageRequest) (*hera.SendMessageReply, error) {
	var request hera.SendMessageRequest

	request.AppId = params.AppID
	request.ProductId = params.ProductID

	if customer.ID != "" {
		request.Customer = &hera.SendMessageRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(params.ChannelNumber).IsZero() {
		request.Customer = &hera.SendMessageRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	if params.Body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Text: &wrapperspb.StringValue{
						Value: params.Body.Text,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Template: &hera.TextMessageTemplate{
						Name:   params.Body.Template.Name,
						Params: params.Body.Template.Params,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Location{
				Location: &hera.LocationMessageBody{
					Latitude:  params.Body.Location.Latitude,
					Longitude: params.Body.Location.Longitude,
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Media{
				Media: &hera.MediaMessageBody{
					Url:   params.Body.Media.URL,
					Media: params.Body.Media.Type,
				},
			},
		}
	}

	request.ChannelNumber = &hera.MessagingChannelNumber{
		Channel: hera.MessagingChannel(params.ChannelNumber.Channel),
		Number:  params.ChannelNumber.Number,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessage(ctx, &request)
}

func (s *service) SendMessageByTag(params *SendMessageByTagRequest) (*hera.TagCommandReply, error) {
	var request hera.SendMessageTagRequest

	request.AppId = params.AppID
	request.ProductId = params.ProductID

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
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Text: &wrapperspb.StringValue{
						Value: params.Body.Text,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Template: &hera.TextMessageTemplate{
						Name:   params.Body.Template.Name,
						Params: params.Body.Template.Params,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Location{
				Location: &hera.LocationMessageBody{
					Latitude:  params.Body.Location.Latitude,
					Longitude: params.Body.Location.Longitude,
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Media{
				Media: &hera.MediaMessageBody{
					Url:   params.Body.Media.URL,
					Media: params.Body.Media.Type,
				},
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendMessageByTag(ctx, &request)
}

func (s *service) ReplyToMessage(customer *Customer, params *ReplyToMessageRequest) (*hera.SendMessageReply, error) {
	var request hera.ReplyToMessageRequest
	request.AppId = params.AppID
	request.ProductId = params.ProductID
	request.CustomerId = customer.ID
	request.ReplyToMessageId = params.ReplyToMessageID
	request.Body = &hera.CustomerMessageBody{}

	if params.Body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Text: &wrapperspb.StringValue{
						Value: params.Body.Text,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Text{
				Text: &hera.TextMessageBody{
					Template: &hera.TextMessageTemplate{
						Name:   params.Body.Template.Name,
						Params: params.Body.Template.Params,
					},
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Location{
				Location: &hera.LocationMessageBody{
					Latitude:  params.Body.Location.Latitude,
					Longitude: params.Body.Location.Longitude,
				},
			},
		}
	}

	if !reflect.ValueOf(params.Body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: &hera.CustomerMessageBody_Media{
				Media: &hera.MediaMessageBody{
					Url:   params.Body.Media.URL,
					Media: params.Body.Media.Type,
				},
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.ReplyToMessage(ctx, &request)
}

func (s *service) MessagingConsent(customer *Customer, params *MessagingConsentRequest) (*hera.MessagingConsentReply, error) {
	var request hera.MessagingConsentRequest
	request.AppId = params.AppID

	request.Action = hera.MessagingConsentAction_MESSAGING_CONSENT_ACTION_OPT_IN

	if customer.ID != "" {
		request.Customer = &hera.MessagingConsentRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(params.ChannelNumber).IsZero() {
		request.Customer = &hera.MessagingConsentRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}
	if !reflect.ValueOf(params.ChannelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(params.ChannelNumber.Channel),
			Number:  params.ChannelNumber.Number,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.MessagingConsent(ctx, &request)
}
