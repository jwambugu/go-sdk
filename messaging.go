package elarian

import (
	"context"
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
		Text         string        `json:"text,omitempty"`
		Media        *Media        `json:"media,omitempty"`
		Location     *Location     `json:"location,omitempty"`
		Template     *Template     `json:"template,omitempty"`
		Ussd         *UssdMenu     `json:"ussd,omitempty"`
		Email        *Email        `json:"email,omitempty"`
		VoiceActions []VoiceAction `json:"voiceActions,omitempty"`
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

func (s *service) SendMessage(customer *Customer, channelNumber *MessagingChannelNumber, body *OutBoundMessageBody) (*hera.SendMessageReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_SendMessage)
	command.SendMessage = &hera.SendMessageCommand{}
	req.Entry = command

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		command.SendMessage.CustomerNumber = s.customerNumber(customer)
	}

	if !reflect.ValueOf(channelNumber).IsZero() {
		command.SendMessage.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}

	message := new(hera.OutboundMessage)
	command.SendMessage.Message = message

	if body.Text != "" {
		message.Body = s.textMessage(body.Text)
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		message.Body = s.templateMesage(body.Template)
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		message.Body = s.locationMessage(body.Location)
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		message.Body = s.mediaMessage(body.Media)
	}
	if !reflect.ValueOf(body.Ussd).IsZero() {
		message.Body = s.ussdMessage(body.Ussd)
	}
	if !reflect.ValueOf(body.VoiceActions).IsZero() {
		message.Body = s.voiceMessage(body.VoiceActions)
	}
	if !reflect.ValueOf(body.Email).IsZero() {
		message.Body = s.email(body.Email)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.SendMessageReply{}, err
	}

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.SendMessageReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetSendMessage(), err
}

func (s *service) SendMessageByTag(tag *Tag, channelNumber *MessagingChannelNumber, body *OutBoundMessageBody) (*hera.TagCommandReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_SendMessageTag)
	req.Entry = command

	if !reflect.ValueOf(tag).IsZero() {
		command.SendMessageTag.Tag = &hera.IndexMapping{
			Key:   tag.Key,
			Value: wrapperspb.String(tag.Value),
		}
	}
	if !reflect.ValueOf(channelNumber).IsZero() {
		command.SendMessageTag.ChannelNumber = &hera.MessagingChannelNumber{
			Channel: hera.MessagingChannel(channelNumber.Channel),
			Number:  channelNumber.Number,
		}
	}

	var message = new(hera.OutboundMessage)
	command.SendMessageTag.Message = message

	if body.Text != "" {
		message.Body = s.textMessage(body.Text)
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		message.Body = s.templateMesage(body.Template)
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		message.Body = s.locationMessage(body.Location)
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		message.Body = s.mediaMessage(body.Media)
	}
	if !reflect.ValueOf(body.Ussd).IsZero() {
		message.Body = s.ussdMessage(body.Ussd)
	}
	if !reflect.ValueOf(body.VoiceActions).IsZero() {
		message.Body = s.voiceMessage(body.VoiceActions)
	}
	if !reflect.ValueOf(body.Email).IsZero() {
		message.Body = s.email(body.Email)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.TagCommandReply{}, err
	}
	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetTagCommand(), err
}

func (s *service) ReplyToMessage(customer *Customer, messageID string, body *OutBoundMessageBody) (*hera.SendMessageReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_ReplyToMessage)
	command.ReplyToMessage = &hera.ReplyToMessageCommand{}
	req.Entry = command

	if customer.ID != "" {
		command.ReplyToMessage.CustomerId = customer.ID
	}
	command.ReplyToMessage.MessageId = messageID

	var message = new(hera.OutboundMessage)
	command.ReplyToMessage.Message = message

	if body.Text != "" {
		message.Body = s.textMessage(body.Text)
	}
	if !reflect.ValueOf(body.Template).IsZero() {
		message.Body = s.templateMesage(body.Template)
	}
	if !reflect.ValueOf(body.Location).IsZero() {
		message.Body = s.locationMessage(body.Location)
	}
	if !reflect.ValueOf(body.Media).IsZero() {
		message.Body = s.mediaMessage(body.Media)
	}
	if !reflect.ValueOf(body.Ussd).IsZero() {
		message.Body = s.ussdMessage(body.Ussd)
	}
	if !reflect.ValueOf(body.VoiceActions).IsZero() {
		message.Body = s.voiceMessage(body.VoiceActions)
	}
	if !reflect.ValueOf(body.Email).IsZero() {
		message.Body = s.email(body.Email)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.SendMessageReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.SendMessageReply{}, err
	}
	reply := &hera.AppToServerCommandReply{}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetSendMessage(), err
}

func (s *service) ReceiveMessage(customerNumber string, channel *MessagingChannelNumber, parts []*InBoundMessageBody) (*hera.SimulatorToServerCommandReply, error) {
	req := new(hera.SimulatorToServerCommand)
	command := new(hera.SimulatorToServerCommand_ReceiveMessage)
	req.Entry = command
	command.ReceiveMessage.CustomerNumber = customerNumber
	command.ReceiveMessage.ChannelNumber = &hera.MessagingChannelNumber{
		Channel: hera.MessagingChannel(channel.Channel),
		Number:  channel.Number,
	}
	command.ReceiveMessage.Parts = []*hera.InboundMessageBody{}
	for _, part := range parts {
		if !reflect.ValueOf(part.Text).IsZero() {
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Text{Text: part.Text},
			})
			continue
		}
		if !reflect.ValueOf(part.Ussd).IsZero() {
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
				Entry: &hera.InboundMessageBody_Ussd{
					Ussd: wrapperspb.String(part.Ussd.Input),
				},
			})
			continue
		}
		if !reflect.ValueOf(part.Email).IsZero() {
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
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
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
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
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
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
			command.ReceiveMessage.Parts = append(command.ReceiveMessage.Parts, &hera.InboundMessageBody{
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.SimulatorToServerCommandReply{}, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.SimulatorToServerCommandReply)
	if err != nil {
		return reply, err
	}
	err = proto.Unmarshal(payload.Data(), reply)
	return reply, err
}
