package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// UssdChannel type
	UssdChannel int32

	// UssdMenu struct
	UssdMenu struct {
		IsTerminal bool   `json:"isTerminal,omitempty"`
		Text       string `json:"text,omitempty"`
	}

	// UssdChannelNumber struct
	UssdChannelNumber struct {
		Channel UssdChannel `json:"channel,omitempty"`
		Number  string      `json:"number,omitempty"`
	}

	// Ussdrequest struct
	UssdOptions struct {
		SessionId string    `json:"sessionId,omitempty"`
		UssdMenu  *UssdMenu `json:"UssdMenu,omitempty"`
	}

	// UssdSessionNotification struct
	UssdSessionNotification struct {
		SessionId      string            `json:"sessionId,omitempty"`
		CustomerId     string            `json:"customerId,omitempty"`
		Input          string            `json:"input,omitempty"`
		CustomerNumber *CustomerNumber   `json:"customerNumber,omitempty"`
		ChannelNumber  UssdChannelNumber `json:"channelNumber,omitempty"`
	}
)

const (
	USSD_CHANNEL_UNSPECIFIED UssdChannel = iota
	USSD_CHANNEL_TELCO
)

func (s *service) ReplyToUssdSession(
	sessionId string,
	ussdMenu *UssdMenu,
) (*hera.WebhookResponseReply, error) {
	var request hera.WebhookResponse
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.SessionId = sessionId
	request.UssdMenu = &hera.UssdMenu{
		IsTerminal: ussdMenu.IsTerminal,
		Text:       ussdMenu.Text,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.SendWebhookResponse(ctx, &request)
}
