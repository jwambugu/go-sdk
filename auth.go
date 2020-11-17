package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) GetAuthToken() (*hera.AuthTokenReply, error) {
	var request hera.AuthTokenRequest
	request.OrgId = s.orgID
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.AuthToken(ctx, &request)
}
