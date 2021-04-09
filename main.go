package elarian

import (
	"context"
	"crypto/tls"
	"log"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	rSocketService struct {
		host         string
		port         int
		errChannel   chan error
		msgChannel   chan *hera.ServerToAppNotification
		replyChannel chan *hera.ServerToAppNotificationReply
	}

	// Options Elarain initialization options
	Options struct {
		OrgID              string
		AppID              string
		APIKey             string
		AuthToken          string
		IsSimulator        bool
		AllowNotifications bool
		Log                bool
	}

	// ConnectionOptions RSocket connection options
	ConnectionOptions struct {
		LifeTime   time.Duration
		Keepalive  time.Duration
		missedAcks int
		Resumable  bool
	}
)

func (s *rSocketService) connect(options *Options, connectionOptions *ConnectionOptions) (rsocket.Client, error) {
	metadata := new(hera.AppConnectionMetadata)
	metadata.OrgId = options.OrgID
	metadata.AppId = options.AppID
	metadata.SimplexMode = !options.IsSimulator
	metadata.SimulatorMode = options.IsSimulator
	metadata.SimplexMode = !options.AllowNotifications

	if options.APIKey != "" {
		metadata.ApiKey = &wrapperspb.StringValue{Value: options.APIKey}
	}
	if options.AuthToken != "" {
		metadata.AuthToken = &wrapperspb.StringValue{Value: options.AuthToken}
	}

	d, err := proto.Marshal(metadata)
	if err != nil {
		log.Fatalln("Error marshaling metadata", err)
	}

	onConnect := func(c rsocket.Client, err error) {
		if err != nil {
			log.Fatalf("Error on connection: %v \n", err)
		}
		if options.Log {
			log.Println("Connected to elarian successfully")
		}
	}

	onClose := func(err error) {
		if err != nil {
			log.Fatalf("Error closing connection: %v \n", err)
		}
		if options.Log {
			log.Println("Elarian connection closed successfully")
		}
	}

	acceptor := func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(func(msg payload.Payload) (response mono.Mono) {
				req := new(hera.ServerToAppNotification)
				if err := proto.Unmarshal(msg.Data(), req); err != nil {
					s.msgChannel <- nil
					s.errChannel <- err
					log.Fatalf("UnMarshaling Error: %v \n", err)
				}
				s.msgChannel <- req
				s.errChannel <- nil

				select {
				case <-time.After(time.Second * 15):
					reply := new(hera.ServerToAppNotificationReply)
					data, err := proto.Marshal(reply)
					if err != nil {
						s.msgChannel <- nil
						s.errChannel <- err
						log.Fatalf("Marshaling Error: %v \n", err)
					}
					return mono.Just(payload.New(data, []byte{}))
				case reply := <-s.replyChannel:
					data, err := proto.Marshal(reply)
					if err != nil {
						s.msgChannel <- nil
						s.errChannel <- err
						log.Fatalf("Marshaling Error: %v \n", err)
					}
					return mono.Just(payload.New(data, []byte{}))
				}
			}))
	}

	tp := rsocket.TCPClient().
		SetHostAndPort(s.host, s.port).
		SetTLSConfig(&tls.Config{ServerName: s.host}).
		Build()

	connectionOpts := new(ConnectionOptions)
	connectionOpts.missedAcks = 6
	connectionOpts.Resumable = true

	if reflect.ValueOf(connectionOptions).IsZero() {
		connectionOpts.Keepalive = time.Duration(time.Second * 2)
		connectionOpts.LifeTime = time.Duration(time.Second * 1)
	} else {
		connectionOpts.Keepalive = connectionOptions.Keepalive
		connectionOpts.LifeTime = connectionOptions.LifeTime
	}

	client, err := rsocket.Connect().
		KeepAlive(connectionOpts.Keepalive, connectionOpts.LifeTime, connectionOpts.missedAcks).
		MetadataMimeType("application/octet-stream").
		DataMimeType("application/octet-stream").
		OnClose(onClose).
		OnConnect(onConnect).
		SetupPayload(payload.New(d, nil)).
		Acceptor(acceptor).
		Transport(tp).
		Start(context.Background())

	if err != nil {
		log.Fatalf("Error on connection: %v \n", err)
	}
	return client, err
}

// Connect establishes a connection to elarian
func Connect(options *Options, connectionOptions *ConnectionOptions) (Service, error) {
	return NewService(options, connectionOptions)
}
