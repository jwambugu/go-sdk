package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// VoiceChannel type
	VoiceChannel int32

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
)

const (
	VOICE_CHANNEL_UNSPECIFIED VoiceChannel = iota
	VOICE_CHANNEL_TELCO
)

func (s *service) transformVoiceCallActions(actions []interface{}) []*hera.VoiceCallAction {
	var voiceActions = []*hera.VoiceCallAction{}

	for _, voiceAction := range actions {
		if action, ok := voiceAction.(VoiceCallActionDequeue); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Dequeue{
					Dequeue: &hera.DequeueCallAction{
						Record:    action.Record,
						QueueName: wrapperspb.String(action.QueueName),
						ChannelNumber: &hera.VoiceChannelNumber{
							Channel: hera.VoiceChannel(action.Channel.Channel),
							Number:  action.Channel.Number,
						},
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionEnqueue); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Enqueue{
					Enqueue: &hera.EnqueueCallAction{
						HoldMusic: wrapperspb.String(action.HoldMusic),
						QueueName: wrapperspb.String(action.QueueName),
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionDail); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Dial{
					Dial: &hera.DialCallAction{
						Record:          action.Record,
						Sequential:      action.Sequential,
						MaxDuration:     wrapperspb.Int32(action.MaxDuration),
						CallerId:        wrapperspb.String(action.CallerId),
						RingbackTone:    wrapperspb.String(action.RingBackTone),
						CustomerNumbers: s.setCustomerNumbers(action.CustomerNumbers),
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionGetDigits); ok && !reflect.ValueOf(action).IsZero() {
			var getDigits *hera.GetDigitsCallAction
			getDigits.FinishOnKey = wrapperspb.String(action.FinishOnKey)
			getDigits.NumDigits = wrapperspb.Int32(action.NumDigits)
			getDigits.Timeout = durationpb.New(action.Timeout)

			if !reflect.ValueOf(action.Prompt.Play).IsZero() {
				getDigits.Prompt = &hera.GetDigitsCallAction_Play{
					Play: &hera.PlayCallAction{
						Url: action.Prompt.Play.URL,
					},
				}
			}

			if !reflect.ValueOf(action.Prompt.Say).IsZero() {
				getDigits.Prompt = &hera.GetDigitsCallAction_Say{
					Say: &hera.SayCallAction{
						PlayBeep: action.Prompt.Say.PlayBeep,
						Text:     action.Prompt.Say.Text,
						Voice:    hera.TextToSpeechVoice(action.Prompt.Say.Voice),
					},
				}
			}

			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_GetDigits{
					GetDigits: getDigits,
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionGetRecording); ok && !reflect.ValueOf(action).IsZero() {
			var getRecording *hera.GetRecordingCallAction
			getRecording.FinishOnKey = wrapperspb.String(action.FinishOnKey)
			getRecording.MaxLength = durationpb.New(action.MaxLength)
			getRecording.PlayBeep = action.PlayBeep
			getRecording.TrimSilence = action.TrimSilence
			getRecording.Timeout = durationpb.New(action.Timeout)
			getRecording.Prompt = &hera.GetRecordingCallAction_Say{}

			if !reflect.ValueOf(action.Prompt.Play).IsZero() {
				getRecording.Prompt = &hera.GetRecordingCallAction_Play{
					Play: &hera.PlayCallAction{
						Url: action.Prompt.Play.URL,
					},
				}
			}

			if !reflect.ValueOf(action.Prompt.Say).IsZero() {
				getRecording.Prompt = &hera.GetRecordingCallAction_Say{
					Say: &hera.SayCallAction{
						PlayBeep: action.Prompt.Say.PlayBeep,
						Text:     action.Prompt.Say.Text,
						Voice:    hera.TextToSpeechVoice(action.Prompt.Say.Voice),
					},
				}
			}

			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_GetRecording{
					GetRecording: getRecording,
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionPlay); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Play{
					Play: &hera.PlayCallAction{
						Url: action.URL,
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionSay); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Say{
					Say: &hera.SayCallAction{
						PlayBeep: action.PlayBeep,
						Text:     action.Text,
						Voice:    hera.TextToSpeechVoice(action.Voice), //FIXME
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionRedirect); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Redirect{
					Redirect: &hera.RedirectCallAction{
						Url: action.URL,
					},
				},
			})
			continue
		}

		if _, ok := voiceAction.(VoiceCallActionReject); ok {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Reject{
					Reject: &hera.RejectCallAction{},
				},
			})
			continue
		}

		if _, ok := voiceAction.(VoiceCallActionRecordSession); ok {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_RecordSession{
					RecordSession: &hera.RecordSessionCallAction{},
				},
			})
			continue
		}
	}
	return voiceActions
}

func (s *service) MakeVoiceCall(
	customer *Customer,
	channel *VoiceChannelNumber,
) (*hera.MakeVoiceCallReply, error) {
	var request hera.MakeVoiceCallRequest
	request.AppId = s.appId
	request.OrgId = s.orgId

	if !reflect.ValueOf(customer.CustomerNumber).IsZero() {
		request.CustomerNumber = s.setCustomerNumber(customer)
	}
	if !reflect.ValueOf(channel).IsZero() {
		request.ChannelNumber = &hera.VoiceChannelNumber{
			Channel: hera.VoiceChannel(channel.Channel),
			Number:  channel.Number,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendWebhookResponse(ctx, &request)
}
