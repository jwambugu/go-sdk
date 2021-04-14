package elarian

import (
	"context"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/golang/protobuf/proto"
	"github.com/rsocket/rsocket-go/payload"
)

type (
	// GenerateAuthTokenReply struct
	GenerateAuthTokenReply struct {
		LifeTime time.Duration `json:"lifeTime,omitempty"`
		Token    string        `json:"token,omitempty"`
	}
)

func (s *elarian) GenerateAuthToken() (*GenerateAuthTokenReply, error) {
	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_GenerateAuthToken{},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := &hera.AppToServerCommandReply{}
	if err := proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &GenerateAuthTokenReply{
		LifeTime: reply.GetGenerateAuthToken().Lifetime.AsDuration(),
		Token:    reply.GetGenerateAuthToken().Token,
	}, nil
}
