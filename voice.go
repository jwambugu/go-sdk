package elarian

import (
	"time"
)

type (
	// VoiceChannel type
	VoiceChannel int32

	// CustomerEventDirection int
	CustomerEventDirection int32

	// VoiceCallStatus int
	VoiceCallStatus int32

	// VoiceCallHangupCause int
	VoiceCallHangupCause int32

	// TextToSpeechVoice int
	TextToSpeechVoice int32

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
		CallerID        string            `json:"callerId,omitempty"`
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
		PlayBeep          bool              `json:"playBeep,omitempty"`
		Text              string            `json:"text,omitempty"`
		TextToSpeechVoice TextToSpeechVoice `json:"textToSpeechVoice,omitempty"`
	}

	// VoiceCallActionRedirect struct
	VoiceCallActionRedirect struct {
		URL string `json:"url,omitempty"`
	}

	// VoiceCallActionReject struct
	VoiceCallActionReject struct{}

	// VoiceCallActionRecordSession struct
	VoiceCallActionRecordSession struct{}

	// VoiceAction interface
	VoiceAction interface {
		voice()
	}

	// Voice struct
	Voice struct {
		Direction    CustomerEventDirection `json:"direction,omitempty"`
		Status       VoiceCallStatus        `json:"status,omitempty"`
		StartedAt    time.Time              `json:"startedAt,omitempty"`
		HangupCase   VoiceCallHangupCause   `json:"hangupCase,omitempty"`
		DtmfDigits   string                 `json:"dtmfDigits,omitempty"`
		RecordingURL string                 `json:"recordingUrl,omitempty"`
		DailData     *VoiceCallDailInput    `json:"dailData,omitempty"`
		QueueData    *VoiceCallQueueInput   `json:"queueData,omitempty"`
	}

	// Prompt struct
	Prompt struct {
		Play VoiceCallActionPlay `json:"voiceCallActionPlay,omitempty"`
		Say  VoiceCallActionSay  `json:"voiceCallActionSay,omitempty"`
	}

	// VoiceCallDailInput struct
	VoiceCallDailInput struct {
		DestinationNumber string        `json:"destinationNumber,omitempty"`
		Duration          time.Duration `json:"duration,omitempty"`
		StartedAt         time.Time     `json:"startedAt,omitempty"`
	}

	// VoiceCallQueueInput struct
	VoiceCallQueueInput struct {
		DequeuedAt          time.Time     `json:"dequeuedAt,omitempty"`
		DequeuedToNumber    string        `json:"dequeuedToNumber,omitempty"`
		DequeuedToSessionID string        `json:"dequeuedToSessionId,omitempty"`
		EnqueuedAt          time.Time     `json:"enqueuedAt,omitempty"`
		QueueDuration       time.Duration `json:"queueDuration,omitempty"`
	}

	// VoiceCallNotification struct
	VoiceCallNotification struct {
		Direction    CustomerEventDirection `json:"direction,omitempty"`
		Status       VoiceCallStatus        `json:"status,omitempty"`
		StartedAt    time.Time              `json:"startedAt,omitempty"`
		HangupCase   VoiceCallHangupCause   `json:"hangupCase,omitempty"`
		DtmfDigits   string                 `json:"dtmfDigits,omitempty"`
		RecordingURL string                 `json:"recordingUrl,omitempty"`
		DailData     *VoiceCallDailInput    `json:"dailData,omitempty"`
		QueueData    *VoiceCallQueueInput   `json:"queueData,omitempty"`
	}
)

// VoiceChannel constants
const (
	VoiceChannelUnspecified VoiceChannel = iota
	VoiceChannelTelco
)

// TextToSpeechVoice
const (
	TextToSpeechVoiceUnspecified TextToSpeechVoice = 0
	TextToSpeechVoiceMale        TextToSpeechVoice = 1
	TextToSpeechVoiceFemale      TextToSpeechVoice = 2
)

// CustomerEventDirect constants
const (
	CustomerEventDirectionUnspecified CustomerEventDirection = iota
	CustomerEventDirectionInbound                            = 1
	CustomerEventDirectionOutbound                           = 2
)

// VoiceCallStatus constants
const (
	VoiceCallStatusUnspecified              VoiceCallStatus = 0
	VoiceCallStatusQueued                   VoiceCallStatus = 100
	VoiceCallStatusAnswered                 VoiceCallStatus = 101
	VoiceCallStatusRinging                  VoiceCallStatus = 102
	VoiceCallStatusActive                   VoiceCallStatus = 200
	VoiceCallStatusDialing                  VoiceCallStatus = 201
	VoiceCallStatusDialCompleted            VoiceCallStatus = 202
	VoiceCallStatusBridged                  VoiceCallStatus = 203
	VoiceCallStatusEnqueued                 VoiceCallStatus = 204
	VoiceCallStatusDequeued                 VoiceCallStatus = 205
	VoiceCallStatusTransferred              VoiceCallStatus = 206
	VoiceCallStatusTransferCompleted        VoiceCallStatus = 207
	VoiceCallStatusCompleted                VoiceCallStatus = 300
	VoiceCallStatusInsufficientCredit       VoiceCallStatus = 400
	VoiceCallStatusNotAnswered              VoiceCallStatus = 401
	VoiceCallStatusInvalidPhoneNumber       VoiceCallStatus = 402
	VoiceCallStatusDestinationNotSupported  VoiceCallStatus = 403
	VoiceCallStatusDecommissionedCustomerid VoiceCallStatus = 404
	VoiceCallStatusExpired                  VoiceCallStatus = 405
	VoiceCallStatusInvalidChannelNumber     VoiceCallStatus = 406
	VoiceCallStatusApplicationError         VoiceCallStatus = 501
)

// VoiceCallHangupCause constants
const (
	VoiceCallHangupCauseUnspecified            VoiceCallHangupCause = 0
	VoiceCallHangupCauseUnallocatedNumber      VoiceCallHangupCause = 1
	VoiceCallHangupCauseUserBusy               VoiceCallHangupCause = 17
	VoiceCallHangupCauseNormalClearing         VoiceCallHangupCause = 16
	VoiceCallHangupCauseNoUserResponse         VoiceCallHangupCause = 18
	VoiceCallHangupCauseNoAnswer               VoiceCallHangupCause = 19
	VoiceCallHangupCauseSubscriberAbsent       VoiceCallHangupCause = 20
	VoiceCallHangupCauseCallRejected           VoiceCallHangupCause = 21
	VoiceCallHangupCauseNormalUnspecified      VoiceCallHangupCause = 31
	VoiceCallHangupCauseNormalTemporaryFailure VoiceCallHangupCause = 41
	VoiceCallHangupCauseServiceUnavailable     VoiceCallHangupCause = 63
	VoiceCallHangupCauseRecoveryOnTimerExpire  VoiceCallHangupCause = 102
	VoiceCallHangupCauseOriginatorCancel       VoiceCallHangupCause = 487
	VoiceCallHangupCauseLoseRace               VoiceCallHangupCause = 502
	VoiceCallHangupCauseUserNotRegistered      VoiceCallHangupCause = 606
)

func (va VoiceCallActionDequeue) voice()       {}
func (va VoiceCallActionEnqueue) voice()       {}
func (va VoiceCallActionDail) voice()          {}
func (va VoiceCallActionGetDigits) voice()     {}
func (va VoiceCallActionGetRecording) voice()  {}
func (va VoiceCallActionPlay) voice()          {}
func (va VoiceCallActionSay) voice()           {}
func (va VoiceCallActionRedirect) voice()      {}
func (va VoiceCallActionReject) voice()        {}
func (va VoiceCallActionRecordSession) voice() {}
