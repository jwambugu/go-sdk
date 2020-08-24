package elarian

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
	productionHost = "api.elarian.dev:443"
)

type (
	// Credentials options required to make a connection to the elarian server.
	Credentials struct {
		APIKey string
	}
	rpcservice struct {
		credentials *Credentials
		tlsConfig   *tls.Config
		keepAlive   keepalive.ClientParameters
	}
)

// GetRequestMetadata returns a map[string]string as another format of the Credentials above
func (c *Credentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"api-key": c.APIKey,
	}, nil
}

// RequireTransportSecurity returns a boolean as whether to enable transport layer security
func (c *Credentials) RequireTransportSecurity() bool {
	return true
}

func (s *rpcservice) connect() (hera.GrpcWebServiceClient, error) {
	var host string
	host = productionHost
	var opts []grpc.DialOption

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

// Initialize creates a secure grpc connection with the elarian server and returns a Service
func Initialize(apikey string) (Service, error) {
	var rpc rpcservice
	rpc.tlsConfig = &tls.Config{
		InsecureSkipVerify: false,
	}
	rpc.keepAlive = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
	rpc.credentials = &Credentials{
		APIKey: apikey,
	}

	client, err := rpc.connect()
	if err != nil {
		return &service{}, err
	}
	return NewService(&client), nil
}