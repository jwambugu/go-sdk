package elariango

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	voiceChannelNumber struct {
		Channel hera.VoiceChannel `json:"channel"`
		Number  string            `json:"number"`
	}

	// VoiceCallRequest struct defines the parameters required to make a voice call request.
	VoiceCallRequest struct {
		AppID     string             `json:"appId,omitempty"`
		ProductID string             `json:"productId,omitempty"`
		Channel   voiceChannelNumber `json:"channel,omitempty"`
	}
)

func (e *elarian) MakeVoiceCall(customer *Customer, params *VoiceCallRequest) (*hera.MakeVoiceCallReply, error) {
	var request hera.MakeVoiceCallRequest

	if customer.ID != "" {
		request.Customer = &hera.MakeVoiceCallRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.MakeVoiceCallRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: customer.PhoneNumber.Provider,
			},
		}
	}
	request.ChannelNumber = &hera.VoiceChannelNumber{
		Channel: params.Channel.Channel,
		Number:  params.Channel.Number,
	}
	request.AppId = params.AppID
	request.ProductId = params.ProductID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return e.client.MakeVoiceCall(ctx, &request)
}
