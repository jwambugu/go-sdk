package elarian

import (
	"context"
	"fmt"
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

	MessagingConsentStatus int32

	MessageDeliveryStatus int32

	MessagingSessionStatus int32

	MediaType int32

	// MessagingChannelNumber struct
	MessagingChannelNumber struct {
		Number  string           `json:"number"`
		Channel MessagingChannel `json:"channel"`
	}

	// Media defines the necessary attributes required to send a file as a message
	Media struct {
		URL  string    `json:"url"`
		Type MediaType `json:"type"`
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

	// PaymentStatusNotification struct
	MessageStatusNotification struct {
		CustomerId string                `json:"customerId,omitempty"`
		Status     MessageDeliveryStatus `json:"status,omitempty"`
		MessageId  string                `json:"messageId,omitempty"`
	}

	MessageSessionStatusNotification struct {
		CustomerId     string                  `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		Expiration     int64                   `json:"expiration,omitempty"`
		Status         MessagingSessionStatus  `json:"status,omitempty"`
	}

	MessagingConsentStatusNotification struct {
		CustomerId     string                  `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		Status         MessagingConsentStatus  `json:"status,omitempty"`
	}

	RecievedMessageNotification struct {
		CustomerId     string                  `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber         `json:"customerNumber,omitempty"`
		ChannelNumber  *MessagingChannelNumber `json:"channelNumber,omitempty"`
		MessageId      string                  `json:"messageId,omitempty"`
		Text           string                  `json:"text"`
		Media          []*Media                `json:"media"`
		Location       *Location               `json:"location"`
		Template       *Template               `json:"template"`
	}
)

const (
	MEDIA_TYPE_UNSPECIFIED MediaType = iota
	MEDIA_TYPE_IMAGE
	MEDIA_TYPE_AUDIO
	MEDIA_TYPE_VIDEO
	MEDIA_TYPE_DOCUMENT
	MEDIA_TYPE_VOICE
	MEDUA_TYPE_STICKER
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

const (
	MESSAGING_SESSION_STATUSN_UNSPECIFIED MessagingSessionStatus = 0
	MESSAGING_SESSION_STATUSN_ACTIVE      MessagingSessionStatus = 100
	MESSAGING_SESSION_STATUSN_EXPIRED     MessagingSessionStatus = 200
)

const (
	MESSAGE_DELIVERY_STATUS_UNSEPCIFIED                 MessageDeliveryStatus = 0
	MESSAGE_DELIVERY_STATUS_SENT                        MessageDeliveryStatus = 101
	MESSAGE_DELIVERY_STATUS_DELIVERED                   MessageDeliveryStatus = 300
	MESSAGE_DELIVERY_STATUS_READ                        MessageDeliveryStatus = 301
	MESSAGE_DELIVERY_STATUS_RECEIVED                    MessageDeliveryStatus = 302
	MESSAGE_DELIVERY_STATUS_FAILED                      MessageDeliveryStatus = 400
	MESSAGE_DELIVERY_STATUS_NO_CONSENT                  MessageDeliveryStatus = 401
	MESSAGE_DELIVERY_STATUS_NO_CAPABILITY               MessageDeliveryStatus = 402
	MESSAGE_DELIVERY_STATUS_EXPIRED                     MessageDeliveryStatus = 403
	MESSAGE_DELIVERY_STATUS_ONLY_TEMPLATE_ALLOWED       MessageDeliveryStatus = 404
	MESSAGE_DELIVERY_STATUS_INVALID_CHANNEL_NUMBER      MessageDeliveryStatus = 405
	MESSAGE_DELIVERY_STATUS_NOT_SUPPORTED               MessageDeliveryStatus = 406
	MESSAGE_DELIVERY_STATUS_INVALID_REPLY_TO_MESSAGE_ID MessageDeliveryStatus = 407
	MESSAGE_DELIVERY_STATUS_INVALID_CUSTOMER_ID         MessageDeliveryStatus = 408
	MESSAGE_DELIVERY_STATUS_DUPLICATE_REQUEST           MessageDeliveryStatus = 409
	MESSAGE_DELIVERY_STATUS_TAG_NOT_FOUND               MessageDeliveryStatus = 410
	MESSAGE_DELIVERY_STATUS_CUSTOMER_NUMBER_NOT_FOUND   MessageDeliveryStatus = 411
	MESSAGE_DELIVERY_STATUS_DECOMMISSIONED_CUSTOMERID   MessageDeliveryStatus = 412
	MESSAGE_DELIVERY_STATUS_INVALID_REQUEST             MessageDeliveryStatus = 413
	MESSAGE_DELIVERY_STATUS_APPLICATION_ERROR           MessageDeliveryStatus = 501
)

const (
	MESSAGING_CONSENT_STATUS_UNSPECIFIED                MessagingConsentStatus = 0
	MESSAGING_CONSENT_STATUS_OPT_IN_REQUEST_SENT        MessagingConsentStatus = 101
	MESSAGING_CONSENT_STATUS_OPT_IN_COMPLETED           MessagingConsentStatus = 300
	MESSAGING_CONSENT_STATUS_OPT_OUT_COMPLETED          MessagingConsentStatus = 301
	MESSAGING_CONSENT_STATUS_INVALID_CHANNEL_NUMBER     MessagingConsentStatus = 401
	MESSAGING_CONSENT_STATUS_DECOMMISSIONED_CUSTOMER_ID MessagingConsentStatus = 402
	MESSAGING_CONSENT_STATUS_APPLICATION_ERROR          MessagingConsentStatus = 501
)

func (s *service) SendMessage(
	customer *Customer,
	channelNumber *MessagingChannelNumber,
	body *MessageBody,
) (*hera.SendMessageReply, error) {
	var request hera.SendMessageRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

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
			Entry: s.messageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(&body.Media),
		}
	}
	fmt.Printf("Request %v \n", &request)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.SendMessage(ctx, &request)
}

func (s *service) SendMessageByTag(
	tag *Tag,
	channelNumber *MessagingChannelNumber,
	body *MessageBody,
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
			Entry: s.messageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(&body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
			Entry: s.messageBodyAsText(body.Text),
		}
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsTemplate(&body.Template),
		}
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsLocation(&body.Location),
		}
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		request.Body = &hera.CustomerMessageBody{
			Entry: s.messageBodyAsMedia(&body.Media),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
func (c *Customer) SendMessage(
	channelNumber *MessagingChannelNumber,
	body *MessageBody,
) (*hera.SendMessageReply, error) {
	return c.service.SendMessage(c, channelNumber, body)
}
