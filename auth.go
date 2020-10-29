package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *service) GetAuthToken() (*hera.AuthTokenReply, error) {
	var request hera.AuthTokenRequest
	request.OrgId = s.orgId

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.AuthToken(ctx, &request)
}
