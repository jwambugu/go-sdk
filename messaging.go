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

	// MessagingConsentStatus int
	MessagingConsentStatus int32

	// MessageDeliveryStatus int
	MessageDeliveryStatus int32

	// MessagingSessionStatus int
	MessagingSessionStatus int32

	// MediaType int
	MediaType int32

	// MessagingChannelNumber struct
	MessagingChannelNumber struct {
		Number  string           `json:"number,omitempty"`
		Channel MessagingChannel `json:"channel,omitempty"`
	}

	// Media defines the necessary attributes required to send a file as a message
	Media struct {
		URL  string    `json:"url,omitempty"`
		Type MediaType `json:"type,omitempty"`
	}

	// Location defines a set of latitude and longitude that can be communicated as a message
	Location struct {
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	}

	// Template This refers to a predefined template for your message, the name of the template is used as the identifier and the params should be added in their logical order
	Template struct {
		Name   string   `json:"name,omitempty"`
		Params []string `json:"params,omitempty"`
	}

	// MessageBody defines how the message body should look like Note all the options are optional and the construction of this struct depends on your needs.
	MessageBody struct {
		Text     string    `json:"text,omitempty"`
		Media    *Media    `json:"media,omitempty"`
		Location *Location `json:"location,omitempty"`
		Template *Template `json:"template,omitempty"`
	}

	// MessageStatusNotification struct
	MessageStatusNotification struct {
		CustomerID string                `json:"customerId,omitempty"`
		Status     MessageDeliveryStatus `json:"status,omitempty"`
		MessageID  string                `json:"messageId,omitempty"`
	}

	// MessageSessionStatusNotification struct
	MessageSessionStatusNotification struct {
		CustomerID     string                  `json:"customerId,omitempty"`
		Expiration     int64                   `json:"expiration,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		Status         MessagingSessionStatus  `json:"status,omitempty"`
	}

	// MessagingConsentStatusNotification struct
	MessagingConsentStatusNotification struct {
		CustomerID     string                  `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		Status         MessagingConsentStatus  `json:"status,omitempty"`
	}

	// RecievedMessageNotification struct
	RecievedMessageNotification struct {
		CustomerID     string                  `json:"customerId,omitempty"`
		MessageID      string                  `json:"messageId,omitempty"`
		Text           string                  `json:"text,omitempty"`
		Media          []*Media                `json:"media,omitempty"`
		Location       *Location               `json:"location,omitempty"`
		Template       *Template               `json:"template,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
	}
)

// MediaType constants
const (
	MediaTypeUnspecified MediaType = iota
	MediaTypeImage
	MediaTypeAudio
	MediaTypeVideo
	MediaTypeDocument
	MediaTypeVoice
	MediaTypeSticker
)

// MessagingChannel constants
const (
	MessagingChannelUnspecified MessagingChannel = iota
	MessagingChannelGoogleRcs
	MessagingChannelFbMessenger
	MessagingChannelSms
	MessagingChannelTelegram
	MessagingChannelWhatsapp
)

// MessagingConsentAction constants
const (
	MessagingConsentActionUnspecified MessagingConsentAction = iota
	MessagingConsentActionOptIn
	MessagingConsentActionOptOut
)

// MessagingSessionStatus constants
const (
	MessagingSessionStatusnUnspecified MessagingSessionStatus = 0
	MessagingSessionStatusnActive      MessagingSessionStatus = 100
	MessagingSessionStatusnExpired     MessagingSessionStatus = 200
)

// MessageDeliveryStatus constants
const (
	MessageDeliveryStatusUnsepcified              MessageDeliveryStatus = 0
	MessageDeliveryStatusSent                     MessageDeliveryStatus = 101
	MessageDeliveryStatusDelivered                MessageDeliveryStatus = 300
	MessageDeliveryStatusRead                     MessageDeliveryStatus = 301
	MessageDeliveryStatusReceived                 MessageDeliveryStatus = 302
	MessageDeliveryStatusFailed                   MessageDeliveryStatus = 400
	MessageDeliveryStatusNoConsent                MessageDeliveryStatus = 401
	MessageDeliveryStatusNoCapability             MessageDeliveryStatus = 402
	MessageDeliveryStatusExpired                  MessageDeliveryStatus = 403
	MessageDeliveryStatusOnlyTemplateAllowed      MessageDeliveryStatus = 404
	MessageDeliveryStatusInvalidChannelNumber     MessageDeliveryStatus = 405
	MessageDeliveryStatusNotSupported             MessageDeliveryStatus = 406
	MessageDeliveryStatusInvalidReplyToMessageID  MessageDeliveryStatus = 407
	MessageDeliveryStatusInvalidCustomerID        MessageDeliveryStatus = 408
	MessageDeliveryStatusDuplicateRequest         MessageDeliveryStatus = 409
	MessageDeliveryStatusTagNotFound              MessageDeliveryStatus = 410
	MessageDeliveryStatusCustomerNumberNotFound   MessageDeliveryStatus = 411
	MessageDeliveryStatusDecommissionedCustomerid MessageDeliveryStatus = 412
	MessageDeliveryStatusInvalidRequest           MessageDeliveryStatus = 413
	MessageDeliveryStatusApplicationError         MessageDeliveryStatus = 501
)

// MessagingConsentStatus constants
const (
	MessagingConsentStatusUnspecified              MessagingConsentStatus = 0
	MessagingConsentStatusOptInRequestSent         MessagingConsentStatus = 101
	MessagingConsentStatusOptInCompleted           MessagingConsentStatus = 300
	MessagingConsentStatusOptOutCompleted          MessagingConsentStatus = 301
	MessagingConsentStatusInvalidChannelNumber     MessagingConsentStatus = 401
	MessagingConsentStatusDecommissionedCustomerID MessagingConsentStatus = 402
	MessagingConsentStatusApplicationError         MessagingConsentStatus = 501
)

func (s *service) SendMessage(customer *Customer, channelNumber *MessagingChannelNumber, body *MessageBody) (*hera.SendMessageReply, error) {
	var request hera.SendMessageRequest
	request.AppId = s.appID
	request.OrgId = s.orgID

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.customerNumber(customer)
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}
	if body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsTemplate(body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(body.Media),
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.SendMessage(ctx, &request)
}

func (s *service) SendMessageByTag(tag *Tag, channelNumber *MessagingChannelNumber, body *MessageBody) (*hera.TagCommandReply, error) {
	var request hera.SendMessageTagRequest
	request.AppId = s.appID
	request.OrgId = s.orgID

	if !reflect.ValueOf(tag).IsZero() {
		request.Tag = &hera.IndexMapping{
			Key:   tag.Key,
			Value: wrapperspb.String(tag.Value),
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
			Entry: s.messageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsTemplate(body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(body.Media),
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.SendMessageByTag(ctx, &request)
}

func (s *service) ReplyToMessage(customer *Customer, messageID string, body *MessageBody) (*hera.SendMessageReply, error) {
	var request hera.ReplyToMessageRequest
	request.AppId = s.appID
	request.OrgId = s.orgID
	request.CustomerId = customer.ID
	request.ReplyToMessageId = messageID

	if body.Text != "" {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsTemplate(body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(body.Media),
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.ReplyToMessage(ctx, &request)
}

func (s *service) MessagingConsent(customer *Customer, channelNumber *MessagingChannelNumber, action MessagingConsentAction) (*hera.MessagingConsentReply, error) {
	var request hera.MessagingConsentRequest
	request.AppId = s.appID
	request.OrgId = s.orgID

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.customerNumber(customer)
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		request.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}
	request.Action = hera.MessagingConsentAction(action)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.MessagingConsent(ctx, &request)
}

// SendMessage sends a messsage to a customer
func (c *Customer) SendMessage(channelNumber *MessagingChannelNumber, body *MessageBody) (*hera.SendMessageReply, error) {
	return c.service.SendMessage(c, channelNumber, body)
}
