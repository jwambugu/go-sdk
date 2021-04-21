package elarian

import (
	"context"
	"crypto/tls"
	"fmt"
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
	service struct {
		host                         string
		port                         int
		errorChannel                 chan<- error
		replyChannel                 <-chan *hera.ServerToAppNotificationReply
		notificationChannel          chan<- *hera.ServerToAppNotification
		simulatorNotificationChannel chan<- *hera.ServerToSimulatorNotification
	}

	// Options Elarain initialization options
	Options struct {
		OrgID              string `json:"orgId,omitempty"`
		AppID              string `json:"appId,omitempty"`
		APIKey             string `json:"apiKey,omitempty"`
		AuthToken          string `json:"authToken,omitempty"`
		IsSimulator        bool   `json:"isSimulator,omitempty"`
		AllowNotifications bool   `json:"allowNotifications,omitempty"`
		Log                bool   `json:"log,omitempty"`
	}

	// ConnectionOptions RSocket connection options
	ConnectionOptions struct {
		LifeTime   time.Duration `json:"lifeTime,omitempty"`
		Keepalive  time.Duration `json:"keepAlive,omitempty"`
		MissedAcks int           `json:"missedAcks,omitempty"`
		Resumable  bool          `json:"resumable,omitempty"`
	}
)

func (s *service) connect(options *Options, connectionOptions *ConnectionOptions) (rsocket.Client, error) {
	metadata := &hera.AppConnectionMetadata{
		OrgId:         options.OrgID,
		AppId:         options.AppID,
		SimulatorMode: options.IsSimulator,
		SimplexMode:   !options.AllowNotifications,
		ApiKey:        wrapperspb.String(options.APIKey),
		AuthToken:     wrapperspb.String(options.AuthToken),
	}

	data, err := proto.Marshal(metadata)
	if err != nil {
		log.Println("Error marshaling metadata", err)
	}

	onConnect := func(c rsocket.Client, err error) {
		if err != nil {
			log.Fatalf("Error on connection: %v \n", err)
			return
		}
		if options.Log {
			log.Println("Connected to elarian successfully")
		}
	}

	onClose := func(err error) {
		if err != nil {
			log.Printf("Error closing connection: %v \n", err)
			return
		}
		close(s.errorChannel)
		close(s.notificationChannel)
		close(s.simulatorNotificationChannel)
		if options.Log {
			log.Println("Elarian connection closed successfully")
		}
	}

	notificationHandler := func(req *hera.ServerToAppNotification) mono.Mono {
		s.notificationChannel <- req
		select {
		case <-time.After(time.Second * 15):
			reply := new(hera.ServerToAppNotificationReply)
			data, _ := proto.Marshal(reply)
			return mono.Just(payload.New(data, []byte{}))
		case reply := <-s.replyChannel:
			data, err := proto.Marshal(reply)
			if err != nil {
				s.errorChannel <- fmt.Errorf("Marshling error: %w ", err)
			}
			return mono.Just(payload.New(data, []byte{}))
		}
	}

	simulatorNotificationHandler := func(req *hera.ServerToSimulatorNotification) mono.Mono {
		s.simulatorNotificationChannel <- req
		reply := new(hera.ServerToSimulatorNotificationReply)
		data, _ := proto.Marshal(reply)
		return mono.Just(payload.New(data, []byte{}))
	}

	acceptor := func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(func(msg payload.Payload) (response mono.Mono) {
				req := new(hera.ServerToAppNotification)
				if err := proto.Unmarshal(msg.Data(), req); err == nil {
					return notificationHandler(req)
				}
				simReq := new(hera.ServerToSimulatorNotification)
				if err := proto.Unmarshal(msg.Data(), simReq); err == nil {
					return simulatorNotificationHandler(simReq)
				}
				reply := new(hera.ServerToSimulatorNotificationReply)
				data, _ := proto.Marshal(reply)
				return mono.Just(payload.New(data, []byte{}))
			}))
	}

	tp := rsocket.TCPClient().
		SetHostAndPort(s.host, s.port).
		SetTLSConfig(&tls.Config{ServerName: s.host}).
		Build()

	connectionOpts := &ConnectionOptions{
		MissedAcks: 6,
		Resumable:  true,
	}

	if reflect.ValueOf(connectionOptions).IsZero() {
		connectionOpts.Keepalive = time.Duration(time.Second * 2)
		connectionOpts.LifeTime = time.Duration(time.Second * 1)
	} else {
		connectionOpts.Keepalive = connectionOptions.Keepalive
		connectionOpts.LifeTime = connectionOptions.LifeTime
	}

	client, err := rsocket.Connect().
		KeepAlive(connectionOpts.Keepalive, connectionOpts.LifeTime, connectionOpts.MissedAcks).
		MetadataMimeType("application/octet-stream").
		DataMimeType("application/octet-stream").
		OnClose(onClose).
		OnConnect(onConnect).
		SetupPayload(payload.New(data, nil)).
		Acceptor(acceptor).
		Transport(tp).
		Start(context.Background())

	if err != nil {
		log.Fatalf("Error on connection: %v \n", err)
	}
	return client, err
}

// Connect establishes a connection to elarian
func Connect(options *Options, connectionOptions *ConnectionOptions) (Elarian, error) {
	return NewService(options, connectionOptions)
}
