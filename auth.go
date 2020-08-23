package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) GetAuthToken(appID string) (*hera.AuthTokenReply, error) {
	var request hera.AuthTokenRequest
	request.AppId = appID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AuthToken(ctx, &request)
}
