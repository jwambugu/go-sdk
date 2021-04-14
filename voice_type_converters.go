package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *elarian) heraVoiceCallActions(actions []VoiceAction) []*hera.VoiceCallAction {
	var voiceActions = []*hera.VoiceCallAction{}

	for _, voiceAction := range actions {
		if action, ok := voiceAction.(VoiceCallActionDequeue); ok && !reflect.ValueOf(action).IsZero() {
			voiceActions = append(voiceActions, &hera.VoiceCallAction{
				Entry: &hera.VoiceCallAction_Dequeue{
					Dequeue: &hera.DequeueCallAction{
						Record:    action.Record,
						QueueName: wrapperspb.String(action.QueueName),
						ChannelNumber: &hera.MessagingChannelNumber{
							Channel: hera.MessagingChannel_MESSAGING_CHANNEL_VOICE,
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
						CallerId:        wrapperspb.String(action.CallerID),
						RingbackTone:    wrapperspb.String(action.RingBackTone),
						CustomerNumbers: s.heraCustomerNumbers(action.CustomerNumbers),
					},
				},
			})
			continue
		}

		if action, ok := voiceAction.(VoiceCallActionGetDigits); ok && !reflect.ValueOf(action).IsZero() {
			getDigits := new(hera.GetDigitsCallAction)
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
						Voice:    hera.TextToSpeechVoice(action.Prompt.Say.TextToSpeechVoice),
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
			getRecording := new(hera.GetRecordingCallAction)
			getRecording.FinishOnKey = wrapperspb.String(action.FinishOnKey)
			getRecording.MaxLength = durationpb.New(action.MaxLength)
			getRecording.PlayBeep = action.PlayBeep
			getRecording.TrimSilence = action.TrimSilence
			getRecording.Timeout = durationpb.New(action.Timeout)

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
						Voice:    hera.TextToSpeechVoice(action.Prompt.Say.TextToSpeechVoice),
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
						Voice:    hera.TextToSpeechVoice(action.TextToSpeechVoice),
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

func (s *elarian) voiceCallActions(actions []*hera.VoiceCallAction) VoiceCallActions {
	var voiceActions = []VoiceAction{}
	for _, voiceCallAction := range actions {
		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Dequeue); ok {
			voiceActions = append(voiceActions, &VoiceCallActionDequeue{
				Channel: VoiceChannelNumber{
					Channel: VoiceChannel(entry.Dequeue.ChannelNumber.Channel),
					Number:  entry.Dequeue.ChannelNumber.Number,
				},
				QueueName: entry.Dequeue.QueueName.Value,
				Record:    entry.Dequeue.Record,
			})
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Enqueue); ok {
			voiceActions = append(voiceActions, VoiceCallActionEnqueue{
				HoldMusic: entry.Enqueue.HoldMusic.Value,
				QueueName: entry.Enqueue.QueueName.Value,
			})
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Dial); ok {
			action := &VoiceCallActionDail{
				Record:          entry.Dial.Record,
				Sequential:      entry.Dial.Sequential,
				MaxDuration:     entry.Dial.MaxDuration.Value,
				CallerID:        entry.Dial.CallerId.Value,
				RingBackTone:    entry.Dial.RingbackTone.Value,
				CustomerNumbers: []*CustomerNumber{},
			}
			for _, number := range entry.Dial.CustomerNumbers {
				action.CustomerNumbers = append(action.CustomerNumbers, &CustomerNumber{
					Number:    number.Number,
					Provider:  NumberProvider(number.Provider),
					Partition: number.Partition.Value,
				})
			}
			voiceActions = append(voiceActions, action)
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_GetDigits); ok {
			getDigits := new(VoiceCallActionGetDigits)
			getDigits.FinishOnKey = entry.GetDigits.FinishOnKey.Value
			getDigits.NumDigits = entry.GetDigits.NumDigits.Value
			getDigits.Timeout = entry.GetDigits.Timeout.AsDuration()

			if prompt, ok := entry.GetDigits.Prompt.(*hera.GetDigitsCallAction_Play); ok {
				getDigits.Prompt = &Prompt{
					Play: &VoiceCallActionPlay{
						URL: prompt.Play.Url,
					},
				}
			}
			if prompt, ok := entry.GetDigits.Prompt.(*hera.GetDigitsCallAction_Say); ok {
				getDigits.Prompt = &Prompt{
					Say: &VoiceCallActionSay{
						PlayBeep:          prompt.Say.PlayBeep,
						Text:              prompt.Say.Text,
						TextToSpeechVoice: TextToSpeechVoice(prompt.Say.Voice),
					},
				}
			}
			voiceActions = append(voiceActions, getDigits)
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_GetRecording); ok {
			recording := new(VoiceCallActionGetRecording)
			recording.FinishOnKey = entry.GetRecording.FinishOnKey.Value
			recording.MaxLength = entry.GetRecording.MaxLength.AsDuration()
			recording.PlayBeep = entry.GetRecording.PlayBeep
			recording.TrimSilence = entry.GetRecording.TrimSilence
			recording.Timeout = entry.GetRecording.Timeout.AsDuration()

			if entry, ok := entry.GetRecording.Prompt.(*hera.GetRecordingCallAction_Say); ok {
				recording.Prompt = &Prompt{
					Say: &VoiceCallActionSay{
						PlayBeep:          entry.Say.PlayBeep,
						Text:              entry.Say.Text,
						TextToSpeechVoice: TextToSpeechVoice(entry.Say.Voice),
					},
				}
			}
			if entry, ok := entry.GetRecording.Prompt.(*hera.GetRecordingCallAction_Play); ok {
				recording.Prompt = &Prompt{
					Play: &VoiceCallActionPlay{
						URL: entry.Play.Url,
					},
				}
			}

			voiceActions = append(voiceActions, recording)
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Play); ok {
			voiceActions = append(voiceActions, &VoiceCallActionPlay{
				URL: entry.Play.Url,
			})
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Say); ok {
			voiceActions = append(voiceActions, &VoiceCallActionSay{
				PlayBeep:          entry.Say.PlayBeep,
				Text:              entry.Say.Text,
				TextToSpeechVoice: TextToSpeechVoice(entry.Say.Voice),
			})
			continue
		}

		if entry, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Redirect); ok {
			voiceActions = append(voiceActions, &VoiceCallActionRedirect{
				URL: entry.Redirect.Url,
			})
			continue
		}

		if _, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_Reject); ok {
			voiceActions = append(voiceActions, &VoiceCallActionReject{})
			continue
		}

		if _, ok := voiceCallAction.Entry.(*hera.VoiceCallAction_RecordSession); ok {
			voiceActions = append(voiceActions, &VoiceCallActionRecordSession{})
			continue
		}

	}
	return voiceActions
}

func (s *elarian) voiceCallNotification(notf *hera.InboundMessageBody_Voice) *Voice {
	return &Voice{
		Direction:    CustomerEventDirection(notf.Voice.Direction),
		Status:       VoiceCallStatus(notf.Voice.Status),
		StartedAt:    notf.Voice.StartedAt.AsTime(),
		HangupCase:   VoiceCallHangupCause(notf.Voice.HangupCause),
		DtmfDigits:   notf.Voice.DtmfDigits.Value,
		RecordingURL: notf.Voice.RecordingUrl.Value,
		DailData: &VoiceCallDailInput{
			DestinationNumber: notf.Voice.DialData.DestinationNumber,
			StartedAt:         notf.Voice.DialData.StartedAt.AsTime(),
			Duration:          notf.Voice.DialData.Duration.AsDuration(),
		},
		QueueData: &VoiceCallQueueInput{
			EnqueuedAt:          notf.Voice.QueueData.EnqueuedAt.AsTime(),
			DequeuedAt:          notf.Voice.QueueData.DequeuedAt.AsTime(),
			DequeuedToNumber:    notf.Voice.QueueData.DequeuedToNumber.Value,
			DequeuedToSessionID: notf.Voice.QueueData.DequeuedToSessionId.Value,
			QueueDuration:       notf.Voice.QueueData.QueueDuration.AsDuration(),
		},
	}
}
