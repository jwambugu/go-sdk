package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
						CustomerNumbers: s.customerNumbers(action.CustomerNumbers),
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

func (s *service) voiceCallNotification(notf *hera.VoiceCallNotification) *VoiceCallNotification {
	return &VoiceCallNotification{
		SessionId: notf.SessionId,
		ChannelNumber: &VoiceChannelNumber{
			Channel: VoiceChannel(notf.ChannelNumber.Channel),
			Number:  notf.ChannelNumber.Number,
		},
		Cost: &Cash{
			CurrencyCode: notf.Cost.CurrencyCode,
			Amount:       notf.Cost.Amount,
		},
		Duration:  notf.Duration.AsDuration(),
		Direction: CustomerEventDirection(notf.Direction),
		Input: &VoiceCallHopInput{
			DtmfDigits:   notf.Input.DtmfDigits.Value,
			RecordingUrl: notf.Input.RecordingUrl.Value,
			Status:       VoiceCallStatus(notf.Input.Status),
			StartedAt:    notf.Input.StartedAt.AsTime(),
			HangupCase:   VoiceCallHangupCause(notf.Input.HangupCause),
			DailData: &VoiceCallDailInput{
				DestinationNumber: notf.Input.DialData.DestinationNumber,
				Duration:          notf.Input.DialData.Duration.AsDuration(),
				StartedAt:         notf.Input.DialData.StartedAt.AsTime(),
			},
			QueueData: &VoiceCallQueueInput{
				DequeuedAt:          notf.Input.QueueData.DequeuedAt.AsTime(),
				DequeuedToNumber:    notf.Input.QueueData.DequeuedToNumber.Value,
				EnqueuedAt:          notf.Input.QueueData.EnqueuedAt.AsTime(),
				QueueDuration:       notf.Input.QueueData.QueueDuration.AsDuration(),
				DequeuedToSessionId: notf.Input.QueueData.DequeuedToSessionId.Value,
			},
		},
	}
}
