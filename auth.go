package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
)

func (s *service) GenerateAuthToken() (*hera.GenerateAuthTokenReply, error) {
	req := new(hera.AppToServerCommand)
	req.Entry = new(hera.AppToServerCommand_GenerateAuthToken)
	reply := new(hera.AppToServerCommandReply)

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.GenerateAuthTokenReply{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.GenerateAuthTokenReply{}, err
	}
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetGenerateAuthToken(), err
}
