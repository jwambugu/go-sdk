package elarian

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/asaskevich/EventBus"
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	rSocketService struct {
		host string
		port int
		bus  EventBus.Bus
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
		LifeTime  time.Duration
		Keepalive time.Duration
		Resumable bool
	}

	// NotificationHandler is an interface implementation required by elarian and is contextual to receiving Notifications from Elarian
	NotificationHandler interface {
		notificationHandler()
	}
)

func (s *rSocketService) connect(options *Options, connectionOptions *ConnectionOptions) (rsocket.Client, error) {
	var metadata = new(hera.AppConnectionMetadata)
	metadata.OrgId = options.OrgID
	metadata.AppId = options.AppID

	if options.APIKey != "" {
		metadata.ApiKey = &wrapperspb.StringValue{Value: options.APIKey}
	}
	if options.AuthToken != "" {
		metadata.AuthToken = &wrapperspb.StringValue{Value: options.AuthToken}
	}

	metadata.SimplexMode = !options.IsSimulator
	metadata.SimulatorMode = options.IsSimulator
	metadata.SimplexMode = !options.AllowNotifications

	d, err := proto.Marshal(metadata)
	if err != nil {
		log.Fatalln("Error marshaling proto", err)
	}

	onConnect := func(c rsocket.Client, err error) {
		if err != nil {
			log.Fatalf("error on connection: %v", err)
		}
		if options.Log {
			log.Println("connected to elarian successfully")
		}
	}

	onClose := func(err error) {
		if err != nil {
			log.Fatalf("error closing connection: %v", err)
		}
		if options.Log {
			log.Println("elarian connection closed successfully")
		}
	}
	acceptor := func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(func(msg payload.Payload) (response mono.Mono) {
				s.bus.Publish("notification", msg)
				log.Println("PAY,OAD", msg.Data())
				req := new(hera.ServerToAppNotification)
				err := proto.Unmarshal(msg.Data(), req)
				if err != nil {
					log.Fatalln("PAYLOAD ERR", err)
				}
				log.Println("PAYLOAD", req)
				return mono.Just(msg)
			}))
	}

	tp := rsocket.TCPClient().
		SetHostAndPort(s.host, s.port).
		SetTLSConfig(&tls.Config{ServerName: s.host}).
		Build()

	client, err := rsocket.Connect().
		Resume().
		KeepAlive(time.Duration(time.Second*120), time.Duration(time.Second*10), 10).
		MetadataMimeType("application/octet-stream").
		DataMimeType("application/octet-stream").
		OnClose(onClose).
		OnConnect(onConnect).
		SetupPayload(payload.New(d, nil)).
		Acceptor(acceptor).
		Transport(tp).
		Start(context.Background())

	if err != nil {
		log.Fatalln("ELARIAN CONNECTION ERROR: ", err)
	}
	return client, err
}

// Connect establishes a connection to elarian
func Connect(options *Options, connectionOptions *ConnectionOptions) (Service, error) {
	return NewService(options, connectionOptions)
}
