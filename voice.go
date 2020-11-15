package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// VoiceChannel type
	VoiceChannel int32

	CustomerEventDirection int32

	VoiceCallStatus int32

	VoiceCallHangupCause int32

	// VoiceChannelNumber struct
	VoiceChannelNumber struct {
		Channel VoiceChannel `json:"channel,omitempty"`
		Number  string       `json:"number,omitempty"`
	}

	// VoiceCallActionDequeue struct
	VoiceCallActionDequeue struct {
		Channel   VoiceChannelNumber `json:"channel,omitempty"`
		QueueName string             `json:"queueName,omitempty"`
		Record    bool               `json:"record,omitempty"`
	}

	// VoiceCallActionEnqueue struct
	VoiceCallActionEnqueue struct {
		HoldMusic string `json:"holdMusic,omitempty"`
		QueueName string `json:"queueName,omitempty"`
	}

	// VoiceCallActionDail struct
	VoiceCallActionDail struct {
		CallerId        string            `json:"callerId,omitempty"`
		MaxDuration     int32             `json:"maxDuration,omitempty"`
		Record          bool              `json:"record,omitempty"`
		RingBackTone    string            `json:"ringBackTone,omitempty"`
		Sequential      bool              `json:"sequential,omitempty"`
		CustomerNumbers []*CustomerNumber `json:"customerNumbers,omitempty"`
	}

	// VoiceCallActionGetDigits struct
	VoiceCallActionGetDigits struct {
		FinishOnKey string        `json:"finishOnKey,omitempty"`
		NumDigits   int32         `json:"numDigits,omitempty"`
		Timeout     time.Duration `json:"timeout,omitempty"`
		Prompt      Prompt        `json:"prompt,omitempty"`
	}

	// VoiceCallActionGetRecording struct
	VoiceCallActionGetRecording struct {
		FinishOnKey string        `json:"finishOnKey,omitempty"`
		MaxLength   time.Duration `json:"maxLength,omitempty"`
		PlayBeep    bool          `json:"playBeep,omitempty"`
		Timeout     time.Duration `json:"timeout,omitempty"`
		TrimSilence bool          `json:"trimSilence,omitempty"`
		Prompt      Prompt        `json:"prompt,omitempty"`
	}

	// VoiceCallActionPlay struct
	VoiceCallActionPlay struct {
		URL string `json:"url,omitempty"`
	}

	// VoiceCallActionSay struct
	VoiceCallActionSay struct {
		PlayBeep bool   `json:"playBeep,omitempty"`
		Text     string `json:"text,omitempty"`
		Voice    int    `json:"voice,omitempty"`
	}

	// VoiceCallActionRedirect struct
	VoiceCallActionRedirect struct {
		URL string `json:"url,omitempty"`
	}

	// VoiceCallActionReject struct
	VoiceCallActionReject struct{}

	// VoiceCallActionRecordSession struct
	VoiceCallActionRecordSession struct{}

	// Prompt struct
	Prompt struct {
		Play VoiceCallActionPlay `json:"voiceCallActionPlay,omitempty"`
		Say  VoiceCallActionSay  `json:"voiceCallActionSay,omitempty"`
	}

	VoiceCallDailInput struct {
		DestinationNumber string        `json:"destinationNumber,omitempty"`
		Duration          time.Duration `json:"duration,omitempty"`
		StartedAt         time.Time     `json:"startedAt,omitempty"`
	}

	VoiceCallQueueInput struct {
		DequeuedAt          time.Time     `json:"dequeuedAt,omitempty"`
		DequeuedToNumber    string        `json:"dequeuedToNumber,omitempty"`
		DequeuedToSessionId string        `json:"dequeuedToSessionId,omitempty"`
		EnqueuedAt          time.Time     `json:"enqueuedAt,omitempty"`
		QueueDuration       time.Duration `json:"queueDuration,omitempty"`
	}

	VoiceCallHopInput struct {
		DtmfDigits   string               `json:"dtmfDigits,omitempty"`
		RecordingUrl string               `json:"recordingUrl,omitempty"`
		StartedAt    time.Time            `json:"startedAt,omitempty"`
		Status       VoiceCallStatus      `json:"status,omitempty"`
		HangupCase   VoiceCallHangupCause `json:"hangupCase,omitempty"`
		DailData     *VoiceCallDailInput  `json:"dailData,omitempty"`
		QueueData    *VoiceCallQueueInput `json:"queueData,omitempty"`
	}

	VoiceCallNotification struct {
		SessionId     string                 `json:"sessionId,omitempty"`
		Cost          *Cash                  `json:"cost,omitempty"`
		Duration      time.Duration          `json:"duration,omitempty"`
		Input         *VoiceCallHopInput     `json:"input,omitempty"`
		Direction     CustomerEventDirection `json:"direction,omitempty"`
		ChannelNumber *VoiceChannelNumber    `json:"channelNumber,omitempty"`
	}
)

const (
	VOICE_CHANNEL_UNSPECIFIED VoiceChannel = iota
	VOICE_CHANNEL_TELCO
)

const (
	CUSTOMER_EVENT_DIRECTION_UNSPECIFIED CustomerEventDirection = iota
	CUSTOMER_EVENT_DIRECTION_INBOUND                            = 1
	CUSTOMER_EVENT_DIRECTION_OUTBOUND                           = 2
)

const (
	VOICE_CALL_STATUS_UNSPECIFIED               VoiceCallStatus = 0
	VOICE_CALL_STATUS_QUEUED                    VoiceCallStatus = 100
	VOICE_CALL_STATUS_ANSWERED                  VoiceCallStatus = 101
	VOICE_CALL_STATUS_RINGING                   VoiceCallStatus = 102
	VOICE_CALL_STATUS_ACTIVE                    VoiceCallStatus = 200
	VOICE_CALL_STATUS_DIALING                   VoiceCallStatus = 201
	VOICE_CALL_STATUS_DIAL_COMPLETED            VoiceCallStatus = 202
	VOICE_CALL_STATUS_BRIDGED                   VoiceCallStatus = 203
	VOICE_CALL_STATUS_ENQUEUED                  VoiceCallStatus = 204
	VOICE_CALL_STATUS_DEQUEUED                  VoiceCallStatus = 205
	VOICE_CALL_STATUS_TRANSFERRED               VoiceCallStatus = 206
	VOICE_CALL_STATUS_TRANSFER_COMPLETED        VoiceCallStatus = 207
	VOICE_CALL_STATUS_COMPLETED                 VoiceCallStatus = 300
	VOICE_CALL_STATUS_INSUFFICIENT_CREDIT       VoiceCallStatus = 400
	VOICE_CALL_STATUS_NOT_ANSWERED              VoiceCallStatus = 401
	VOICE_CALL_STATUS_INVALID_PHONE_NUMBER      VoiceCallStatus = 402
	VOICE_CALL_STATUS_DESTINATION_NOT_SUPPORTED VoiceCallStatus = 403
	VOICE_CALL_STATUS_DECOMMISSIONED_CUSTOMERID VoiceCallStatus = 404
	VOICE_CALL_STATUS_EXPIRED                   VoiceCallStatus = 405
	VOICE_CALL_STATUS_INVALID_CHANNEL_NUMBER    VoiceCallStatus = 406
	VOICE_CALL_STATUS_APPLICATION_ERROR         VoiceCallStatus = 501
)

const (
	VOICE_CALL_HANGUP_CAUSE_UNSPECIFIED              VoiceCallHangupCause = 0
	VOICE_CALL_HANGUP_CAUSE_UNALLOCATED_NUMBER       VoiceCallHangupCause = 1
	VOICE_CALL_HANGUP_CAUSE_USER_BUSY                VoiceCallHangupCause = 17
	VOICE_CALL_HANGUP_CAUSE_NORMAL_CLEARING          VoiceCallHangupCause = 16
	VOICE_CALL_HANGUP_CAUSE_NO_USER_RESPONSE         VoiceCallHangupCause = 18
	VOICE_CALL_HANGUP_CAUSE_NO_ANSWER                VoiceCallHangupCause = 19
	VOICE_CALL_HANGUP_CAUSE_SUBSCRIBER_ABSENT        VoiceCallHangupCause = 20
	VOICE_CALL_HANGUP_CAUSE_CALL_REJECTED            VoiceCallHangupCause = 21
	VOICE_CALL_HANGUP_CAUSE_NORMAL_UNSPECIFIED       VoiceCallHangupCause = 31
	VOICE_CALL_HANGUP_CAUSE_NORMAL_TEMPORARY_FAILURE VoiceCallHangupCause = 41
	VOICE_CALL_HANGUP_CAUSE_SERVICE_UNAVAILABLE      VoiceCallHangupCause = 63
	VOICE_CALL_HANGUP_CAUSE_RECOVERY_ON_TIMER_EXPIRE VoiceCallHangupCause = 102
	VOICE_CALL_HANGUP_CAUSE_ORIGINATOR_CANCEL        VoiceCallHangupCause = 487
	VOICE_CALL_HANGUP_CAUSE_LOSE_RACE                VoiceCallHangupCause = 502
	VOICE_CALL_HANGUP_CAUSE_USER_NOT_REGISTERED      VoiceCallHangupCause = 606
)

func (s *service) MakeVoiceCall(
	customer *Customer,
	channel *VoiceChannelNumber,
) (*hera.MakeVoiceCallReply, error) {
	var request hera.MakeVoiceCallRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.customerNumber(customer)
	}
	if !reflect.ValueOf(channel).IsZero() {
		request.ChannelNumber = &hera.VoiceChannelNumber{
			Channel: hera.VoiceChannel(channel.Channel),
			Number:  channel.Number,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.MakeVoiceCall(ctx, &request)
}

func (s *service) ReplyToVoiceCall(
	sessionId string,
	actions []interface{},
) (*hera.WebhookResponseReply, error) {
	var request hera.WebhookResponse
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.SessionId = sessionId
	request.VoiceCallActions = s.transformVoiceCallActions(actions)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.SendWebhookResponse(ctx, &request)
}
