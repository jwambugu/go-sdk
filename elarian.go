package elariango

import (
	"context"
	"crypto/tls"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	productionHost = "api.elarian.com:443"
	sandBoxHost    = "api.elarian.dev:443"
)

type (
	elarianCredentials struct {
		APIKey    string
		AuthToken string
	}
	service struct {
		options struct {
			SandBox bool
		}
		credentials *elarianCredentials
		tlsConfig   *tls.Config
		keepAlive   keepalive.ClientParameters
	}
)

func (c *elarianCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"api-key":    c.APIKey,
		"auth-token": c.AuthToken,
	}, nil
}

func (c *elarianCredentials) RequireTransportSecurity() bool {
	return true
}

func (s *service) connect() (hera.GrpcWebServiceClient, error) {
	var host string
	var opts []grpc.DialOption

	if s.options.SandBox {
		host = sandBoxHost
	} else {
		host = productionHost
	}

	creds := credentials.NewTLS(s.tlsConfig)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(s.credentials))
	opts = append(opts, grpc.WithKeepaliveParams(s.keepAlive))
	opts = append(opts, grpc.WithBlock())

	channel, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, err
	}
	client := hera.NewGrpcWebServiceClient(channel)
	return client, nil
}

// Initialize creates a secure grpc connection with the elarian api and returns an elarian grpc web service client
func Initialize(apikey string, sandbox bool) (Elarian, error) {
	var elarianService service

	elarianService.tlsConfig = &tls.Config{
		InsecureSkipVerify: false,
	}
	elarianService.keepAlive = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
	elarianService.credentials = &elarianCredentials{
		APIKey:    apikey,
		AuthToken: "",
	}
	elarianService.options.SandBox = sandbox

	client, err := elarianService.connect()
	if err != nil {
		return &elarian{}, err
	}
	return NewService(&client), nil
}
