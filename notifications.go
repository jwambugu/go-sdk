package elariango

import (
	"context"
	"io"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (e *elarian) SendWebhookResponse() (*hera.WebhookResponseReply, error) {
	var request hera.WebhookResponse
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return e.client.SendWebhookResponse(ctx, &request)
}

func (e *elarian) StreamNotifications(appID string) (chan *hera.WebhookRequest, chan error) {
	var request hera.StreamNotificationRequest
	request.AppId = appID

	ctx := context.Background()
	stream, err := e.client.StreamNotifications(ctx, &request)
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
			}
			streamChannel <- in
		}
	}()
	return streamChannel, nil
}
