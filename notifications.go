package elarian

import (
	"context"
	"io"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// USSDMenu struct
	USSDMenu struct {
		IsTerminal bool
		Text       string
	}
	// WebhookRequest struct
	WebhookRequest struct {
		AppID     string
		SessionID string
		USSDMenu  USSDMenu
	}
)

func (s *service) SendWebhookResponse(params *WebhookRequest) (*hera.WebhookResponseReply, error) {
	var request hera.WebhookResponse
	request.AppId = params.AppID
	request.SessionId = params.SessionID
	request.UssdMenu = &hera.UssdMenu{
		IsTerminal: params.USSDMenu.IsTerminal,
		Text:       params.USSDMenu.Text,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendWebhookResponse(ctx, &request)
}

func (s *service) StreamNotifications(appID string) (chan *hera.WebhookRequest, chan error) {
	var request hera.StreamNotificationRequest
	request.AppId = appID

	ctx := context.Background()
	stream, err := s.client.StreamNotifications(ctx, &request)
	streamChannel := make(chan *hera.WebhookRequest)
	errorChannel := make(chan error)
	if err != nil {
		return streamChannel, errorChannel
	}
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(streamChannel)
				close(errorChannel)
				return
			}
			if err != nil {
				errorChannel <- err
				close(streamChannel)
				close(errorChannel)
				return
			}
			streamChannel <- in
		}
	}()
	return streamChannel, nil
}
