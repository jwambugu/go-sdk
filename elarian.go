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

	// Options required to start get an elarian instance
	Options struct {
		ApiKey string `json:"apiKey,omitempty"`
		OrgId  string `json:"orgId,omitempty"`
		AppId  string `json:"appId,omitempty"`
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
	var opts []grpc.DialOption

	creds := credentials.NewTLS(s.tlsConfig)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(s.credentials))
	opts = append(opts, grpc.WithKeepaliveParams(s.keepAlive))
	opts = append(opts, grpc.WithBlock())

	channel, err := grpc.Dial("api.elarian.dev:443", opts...)
	if err != nil {
		return nil, err
	}
	client := hera.NewGrpcWebServiceClient(channel)
	return client, nil
}

// Initialize creates a secure grpc connection with the elarian server and returns a Service
func Initialize(opts *Options) (Service, error) {
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
		APIKey: opts.ApiKey,
	}

	client, err := rpc.connect()
	if err != nil {
		return &service{}, err
	}
	return NewService(&client, opts.OrgId, opts.AppId), nil
}
