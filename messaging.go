package elarian

import (
	"context"
	"errors"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (

	// MessageReaction Enum
	MessageReaction int32

	// MessagingChannel is an enum
	MessagingChannel int32

	// MessagingConsentUpdate Enum
	MessagingConsentUpdate int32

	// MessagingConsentUpdateStatus int
	MessagingConsentUpdateStatus int32

	// MessageDeliveryStatus int
	MessageDeliveryStatus int32

	// MessagingSessionStatus int
	MessagingSessionStatus int32

	// MessagingSessionEndReason enum
	MessagingSessionEndReason int32

	// MediaType int
	MediaType int32

	// IsOutBoundMessageBody interface
	IsOutBoundMessageBody interface {
		isOutBoundMessageBody()
	}

	// TextMessage implements the IsOutBoundMessageBody interface
	TextMessage string

	// URLMessage implements the IsOutBoundMessageBody interface
	URLMessage string

	// VoiceCallActions implements the IsOutBoundMessageBody interface
	VoiceCallActions []VoiceAction

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
		Label     string  `json:"label,omitempty"`
		Address   string  `json:"address,omitempty"`
	}

	// Template This refers to a predefined template for your message, the name of the template is used as the identifier and the params should be added in their logical order
	Template struct {
		ID     string            `json:"name,omitempty"`
		Params map[string]string `json:"params,omitempty"`
	}

	// Email defines the email fields you can use as a message
	Email struct {
		Subject     string   `json:"subject,omitempty"`
		Body        string   `json:"body,omitempty"`
		HTML        string   `json:"html,omitempty"`
		CcList      []string `json:"ccList,omitempty"`
		BccList     []string `json:"bccList,omitempty"`
		Attachments []string `json:"attachments,omitempty"`
	}

	// OutBoundMessageBody defines how the message body should look like Note all the options are optional and the construction of this struct depends on your needs.
	OutBoundMessageBody struct {
		Entry IsOutBoundMessageBody `json:"entry,omitempty"`
	}

	// PromptMessageReplyAction enum
	PromptMessageReplyAction int32

	// PromptMessageMenuItemBody struct
	PromptMessageMenuItemBody struct {
		Text  string `json:"text,omitempty"`
		Media *Media `json:"media,omitempty"`
	}

	// OutboundMessageReplyPrompt struct
	OutboundMessageReplyPrompt struct {
		Action PromptMessageReplyAction     `json:"action,omitempty"`
		Menu   []*PromptMessageMenuItemBody `json:"menu,omitempty"`
	}

	// OutBoundMessage struct
	OutBoundMessage struct {
		Body        *OutBoundMessageBody        `json:"body,omitempty"`
		Labels      []string                    `json:"labels,omitempty"`
		ProviderTag string                      `json:"providerTag,omitempty"`
		ReplyToken  string                      `json:"replyToken,omitempty"`
		ReplyPrompt *OutboundMessageReplyPrompt `json:"replyPrompt,omitempty"`
	}

	// InBoundMessageBody defines how the message body should look like Note all the options are optional and the construction of this struct depends on your needs.
	InBoundMessageBody struct {
		Text     string                   `json:"text,omitempty"`
		Media    *Media                   `json:"media,omitempty"`
		Location *Location                `json:"location,omitempty"`
		Template *Template                `json:"template,omitempty"`
		Ussd     *UssdSessionNotification `json:"ussd,omitempty"`
		Email    *Email                   `json:"email,omitempty"`
		Voice    *Voice                   `json:"voice,omitempty"`
	}

	// SendMessageReply struct
	SendMessageReply struct {
		CustomerID  string                `json:"customerId,omitempty"`
		Description string                `json:"description,omitempty"`
		MessageID   string                `json:"messageId,omitempty"`
		SessionID   string                `json:"sessionId,omitempty"`
		Status      MessageDeliveryStatus `json:"status,omitempty"`
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
	MediaTypeContact
)

// MessagingChannel constants
const (
	MessagingChannelUnspecified MessagingChannel = iota
	MessagingChannelSms
	MessagingChannelVoice
	MessagingChannelUssd
	MessagingChannelFBMessanger
	MessagingChannelTelegram
	MessagingChannelWhatsapp
	MessagingChannelEmail
)

// MessagingSessionStatus constants
const (
	MessagingSessionStatusnUnspecified MessagingSessionStatus = 0
	MessagingSessionStatusnActive      MessagingSessionStatus = 100
	MessagingSessionStatusnExpired     MessagingSessionStatus = 200
)

// MessagingSessionEndReason constants
const (
	MessagingSessionEndReasonUnspecified    MessagingSessionEndReason = 0
	MessagingSessionEndReasonNormalClearing MessagingSessionEndReason = 100
	MessagingSessionEndReasonInactivity     MessagingSessionEndReason = 200
	MessagingSessionEndReasonFailure        MessagingSessionEndReason = 300
)

// MessageDeliveryStatus constants
const (
	MessageDeliveryStatusUnspecified              MessageDeliveryStatus = 0
	MessageDeliveryStatusQueued                   MessageDeliveryStatus = 100
	MessageDeliveryStatusSent                     MessageDeliveryStatus = 101
	MessageDeliveryStatusDelivered                MessageDeliveryStatus = 300
	MessageDeliveryStatusRead                     MessageDeliveryStatus = 301
	MessageDeliveryStatusReceived                 MessageDeliveryStatus = 302
	MessageDeliveryStatusSessionInitiated         MessageDeliveryStatus = 303
	MessageDeliveryStatusFailed                   MessageDeliveryStatus = 400
	MessageDeliveryStatusNoConsent                MessageDeliveryStatus = 401
	MessageDeliveryStatusNoCapability             MessageDeliveryStatus = 402
	MessageDeliveryStatusExpired                  MessageDeliveryStatus = 403
	MessageDeliveryStatusNoSessionInProgress      MessageDeliveryStatus = 404
	MessageDeliveryStatusOtherSessionInProgress   MessageDeliveryStatus = 405
	MessageDeliveryStatusInvalidReplyToken        MessageDeliveryStatus = 406
	MessageDeliveryStatusInvalidChannelNumber     MessageDeliveryStatus = 407
	MessageDeliveryStatusNotSupported             MessageDeliveryStatus = 408
	MessageDeliveryStatusInvalidReplyToMessageID  MessageDeliveryStatus = 409
	MessageDeliveryStatusInvalidCustomerID        MessageDeliveryStatus = 410
	MessageDeliveryStatusDuplicateRequest         MessageDeliveryStatus = 411
	MessageDeliveryStatusTagNotFound              MessageDeliveryStatus = 412
	MessageDeliveryStatusCustomerNumberNotFound   MessageDeliveryStatus = 413
	MessageDeliveryStatusDecommissionedCustomerid MessageDeliveryStatus = 414
	MessageDeliveryStatusRejected                 MessageDeliveryStatus = 415
	MessageDeliveryStatusInvalidRequest           MessageDeliveryStatus = 416
	MessageDeliveryStatusApplicationError         MessageDeliveryStatus = 501
)

// MessagingConsentStatus constants
const (
	MessagingConsentStatusUnspecified              MessagingConsentUpdateStatus = 0
	MessagingConsentStatusQueued                   MessagingConsentUpdateStatus = 100
	MessagingConsentStatusCompleted                MessagingConsentUpdateStatus = 300
	MessagingConsentStatusInvalidChannelNumber     MessagingConsentUpdateStatus = 401
	MessagingConsentStatusDecommissionedCustomerID MessagingConsentUpdateStatus = 402
	MessagingConsentStatusApplicationError         MessagingConsentUpdateStatus = 501
)

// MesageReaction constants
const (
	MessageReactionUnspecified  MessageReaction = 0
	MessageReactionClicked      MessageReaction = 100
	MessageReactionUnsubscribed MessageReaction = 200
	MessageReactionComplained   MessageReaction = 201
)

// PromptMessageReply constants
const (
	PromptMessageReplyActionUnspecified PromptMessageReplyAction = iota
	PromptMessageReplyActionText
	PromptMessageReplyActionPhoneNumber
	PromptMessageReplyActionEmail
	PromptMessageReplyActionLocation
	PromptMessageReplyActionURL
)

func (TextMessage) isOutBoundMessageBody()      {}
func (*Media) isOutBoundMessageBody()           {}
func (*Location) isOutBoundMessageBody()        {}
func (*Template) isOutBoundMessageBody()        {}
func (*UssdMenu) isOutBoundMessageBody()        {}
func (*Email) isOutBoundMessageBody()           {}
func (URLMessage) isOutBoundMessageBody()       {}
func (VoiceCallActions) isOutBoundMessageBody() {}

func (s *elarian) SendMessage(ctx context.Context, number *CustomerNumber, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	if number == nil || reflect.ValueOf(number).IsZero() {
		return nil, errors.New("customerNumber required")
	}

	if channelNumber == nil || reflect.ValueOf(channelNumber).IsZero() {
		return nil, errors.New("channelNumber required")
	}

	message := &hera.OutboundMessage{}
	if entry, ok := body.(TextMessage); ok {
		message.Body = s.heraOutBoundTextMessage(string(entry))
	}
	if entry, ok := body.(*Template); ok {
		message.Body = s.heraOutBoundTemplateMesage(entry)
	}
	if entry, ok := body.(*Location); ok {
		message.Body = s.heraOutBoundLocationMessage(entry)
	}
	if entry, ok := body.(*Media); ok {
		message.Body = s.heraOutBoundMediaMessage(entry)
	}
	if entry, ok := body.(*UssdMenu); ok {
		message.Body = s.heraOutBoundUssdMessage(entry)
	}
	if entry, ok := body.(*Email); ok {
		message.Body = s.heraOutBoundEmail(entry)
	}
	if entry, ok := body.(VoiceCallActions); ok {
		message.Body = s.heraOutBoundVoiceMessage(entry)
	}

	command := &hera.SendMessageCommand{
		CustomerNumber: s.heraCustomerNumber(number),
		ChannelNumber: &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		},
		Message: message,
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_SendMessage{SendMessage: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &SendMessageReply{
		CustomerID:  reply.GetSendMessage().CustomerId.Value,
		Description: reply.GetSendMessage().Description,
		MessageID:   reply.GetSendMessage().MessageId.Value,
		Status:      MessageDeliveryStatus(reply.GetSendMessage().Status),
	}, nil
}

func (s *elarian) SendMessageByTag(ctx context.Context, tag *Tag, channelNumber *MessagingChannelNumber, body IsOutBoundMessageBody) (*TagCommandReply, error) {
	if tag == nil || reflect.ValueOf(tag).IsZero() {
		return nil, errors.New("tag is required")
	}

	if channelNumber == nil || reflect.ValueOf(channelNumber).IsZero() {
		return nil, errors.New("channelNumber is required")
	}

	var message = &hera.OutboundMessage{}
	if entry, ok := body.(TextMessage); ok {
		message.Body = s.heraOutBoundTextMessage(string(entry))
	}
	if entry, ok := body.(*Template); ok {
		message.Body = s.heraOutBoundTemplateMesage(entry)
	}
	if entry, ok := body.(*Location); ok {
		message.Body = s.heraOutBoundLocationMessage(entry)
	}
	if entry, ok := body.(*Media); ok {
		message.Body = s.heraOutBoundMediaMessage(entry)
	}
	if entry, ok := body.(*UssdMenu); ok {
		message.Body = s.heraOutBoundUssdMessage(entry)
	}
	if entry, ok := body.(*Email); ok {
		message.Body = s.heraOutBoundEmail(entry)
	}
	if entry, ok := body.(VoiceCallActions); ok {
		message.Body = s.heraOutBoundVoiceMessage(entry)
	}

	command := &hera.SendMessageTagCommand{
		Tag: &hera.IndexMapping{
			Key:   tag.Key,
			Value: wrapperspb.String(tag.Value),
		},
		ChannelNumber: &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		},
		Message: message,
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_SendMessageTag{SendMessageTag: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.AppToServerCommandReply)
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &TagCommandReply{
		Status:      reply.GetTagCommand().Status,
		Description: reply.GetTagCommand().Description,
		WorkID:      reply.GetTagCommand().WorkId.Value,
	}, err
}

func (s *elarian) ReplyToMessage(ctx context.Context, customerID, messageID string, body IsOutBoundMessageBody) (*SendMessageReply, error) {
	var message = &hera.OutboundMessage{}
	if entry, ok := body.(TextMessage); ok {
		message.Body = s.heraOutBoundTextMessage(string(entry))
	}
	if entry, ok := body.(*Template); ok {
		message.Body = s.heraOutBoundTemplateMesage(entry)
	}
	if entry, ok := body.(*Location); ok {
		message.Body = s.heraOutBoundLocationMessage(entry)
	}
	if entry, ok := body.(*Media); ok {
		message.Body = s.heraOutBoundMediaMessage(entry)
	}
	if entry, ok := body.(*UssdMenu); ok {
		message.Body = s.heraOutBoundUssdMessage(entry)
	}
	if entry, ok := body.(*Email); ok {
		message.Body = s.heraOutBoundEmail(entry)
	}
	if entry, ok := body.(VoiceCallActions); ok {
		message.Body = s.heraOutBoundVoiceMessage(entry)
	}

	command := &hera.ReplyToMessageCommand{
		CustomerId: customerID,
		MessageId:  messageID,
		Message:    message,
	}
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_ReplyToMessage{ReplyToMessage: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &SendMessageReply{
		CustomerID:  reply.GetSendMessage().CustomerId.Value,
		Description: reply.GetSendMessage().Description,
		MessageID:   reply.GetSendMessage().MessageId.Value,
		Status:      MessageDeliveryStatus(reply.GetSendMessage().Status),
	}, err
}

func (s *elarian) ReceiveMessage(ctx context.Context, customerNumber string, channel *MessagingChannelNumber, sessionID string, parts []*InBoundMessageBody) (*SimulatorToServerCommandReply, error) {
	messageparts := []*hera.InboundMessageBody{}
	for _, part := range parts {
		if !reflect.ValueOf(part.Text).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Text{Text: part.Text},
			})
			continue
		}
		if !reflect.ValueOf(part.Ussd).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Ussd{
					Ussd: wrapperspb.String(part.Ussd.Input),
				},
			})
			continue
		}
		if !reflect.ValueOf(part.Email).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Email{
					Email: &hera.EmailMessageBody{
						BodyPlain:   part.Email.Body,
						BodyHtml:    part.Email.HTML,
						CcList:      part.Email.CcList,
						BccList:     part.Email.BccList,
						Attachments: part.Email.Attachments,
						Subject:     part.Email.Subject,
					},
				},
			})
			continue
		}
		if !reflect.ValueOf(part.Location).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Location{
					Location: &hera.LocationMessageBody{
						Latitude:  part.Location.Latitude,
						Longitude: part.Location.Longitude,
						Label:     wrapperspb.String(part.Location.Label),
						Address:   wrapperspb.String(part.Location.Address),
					},
				},
			})
			continue
		}

		if !reflect.ValueOf(part.Media).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Media{
					Media: &hera.MediaMessageBody{
						Url:   part.Media.URL,
						Media: hera.MediaType(part.Media.Type),
					},
				},
			})
			continue
		}

		if !reflect.ValueOf(part.Voice).IsZero() {
			messageparts = append(messageparts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Voice{Voice: &hera.VoiceCallInputMessageBody{
					Direction:    hera.CustomerEventDirection(part.Voice.Direction),
					Status:       hera.VoiceCallStatus(part.Voice.Status),
					StartedAt:    timestamppb.New(part.Voice.StartedAt),
					HangupCause:  hera.VoiceCallHangupCause(part.Voice.HangupCase),
					DtmfDigits:   wrapperspb.String(part.Voice.DtmfDigits),
					RecordingUrl: wrapperspb.String(part.Voice.RecordingURL),
					DialData: &hera.VoiceCallDialInput{
						DestinationNumber: part.Voice.DailData.DestinationNumber,
						StartedAt:         timestamppb.New(part.Voice.DailData.StartedAt),
						Duration:          durationpb.New(part.Voice.DailData.Duration),
					},
					QueueData: &hera.VoiceCallQueueInput{
						EnqueuedAt:          timestamppb.New(part.Voice.QueueData.EnqueuedAt),
						DequeuedAt:          timestamppb.New(part.Voice.QueueData.DequeuedAt),
						DequeuedToNumber:    wrapperspb.String(part.Voice.QueueData.DequeuedToNumber),
						DequeuedToSessionId: wrapperspb.String(part.Voice.QueueData.DequeuedToSessionID),
						QueueDuration:       durationpb.New(part.Voice.QueueData.QueueDuration),
					},
				}},
			})
		}
	}

	command := &hera.ReceiveMessageSimulatorCommand{
		CustomerNumber: customerNumber,
		ChannelNumber: &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channel.Channel),
			Number:  channel.Number,
		},
		Parts:     messageparts,
		SessionId: wrapperspb.String(sessionID),
	}
	req := &hera.SimulatorToServerCommand{
		Entry: &hera.SimulatorToServerCommand_ReceiveMessage{ReceiveMessage: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.SimulatorToServerCommandReply)
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &SimulatorToServerCommandReply{
		Status:      reply.Status,
		Message:     s.OutboundMessage(reply.Message),
		Description: reply.Description,
	}, nil
}
