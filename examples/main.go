// Package main implements a the simple examples that demonstrate hot to use the elarian SDK.
package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/elarianltd/go-sdk"
	elarian "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func sendSms(client elarian.GrpcWebServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	smsRequest := elarian.SendMessageRequest{
		AppId:     "app_id",
		ProductId: "product_id",
		ChannelNumber: &elarian.MessagingChannelNumber{
			Channel: 3,
			Number:  "41012",
		},
		Body: &elarian.CustomerMessageBody{
			Entry: &elarian.CustomerMessageBody_Text{
				Text: &elarian.TextMessageBody{
					Text: &wrapperspb.StringValue{
						Value: "Test Message",
					},
				},
			},
		},
		Customer: &elarian.SendMessageRequest_CustomerNumber{
			CustomerNumber: &elarian.CustomerNumber{
				Number:   "+2547",
				Provider: 2,
			},
		},
	}
	res, err := client.SendMessage(ctx, &smsRequest)
	if err != nil {
		log.Fatalf("Could not send sms %v", err)
	}
	log.Printf("Customer id %v", res.GetCustomerId())
}

func getCustomerState(client elarian.GrpcWebServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.GetCustomerState(ctx, &elarian.GetCustomerStateRequest{
		AppId: "app_id",
		Customer: &elarian.GetCustomerStateRequest_CustomerId{
			CustomerId: "customer_id",
		},
	})
	if err != nil {
		log.Fatalf("could not get customer state %v", err)
	}
	log.Printf("customer state %v", res.GetData())
}

func addCustomerReminder(client elarian.GrpcWebServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30000*time.Second)
	defer cancel()
	exp := time.Now().Add(time.Minute + 1)

	res, err := client.AddCustomerReminder(ctx, &elarian.AddCustomerReminderRequest{
		AppId: "app_id",
		Customer: &elarian.AddCustomerReminderRequest_CustomerId{
			CustomerId: "customer_id",
		},
		Reminder: &elarian.CustomerReminder{
			ProductId:  "product_id",
			Expiration: timestamppb.New(exp),
			Key:        "12345",
			Payload: &wrapperspb.StringValue{
				Value: "i am a reminder",
			},
		},
	})
	if err != nil {
		log.Fatalf("could not set a reminder %v", err)
	}
	log.Printf("response %v", res)
}

func streamCustomerNotifications(client elarian.GrpcWebServiceClient) {
	ctx := context.Background()

	stream, err := client.StreamNotifications(ctx, &elarian.StreamNotificationRequest{
		AppId: "app_id",
	})
	if err != nil {
		log.Fatalf("could stream notification %v", err)
	}
	waitChannel := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitChannel)
				return
			}
			if err != nil {
				log.Fatalf("Failed to recieve  notifications: %v", err)
			}
			log.Printf("Got a notification %v", in.GetReminder())
		}
	}()
	<-waitChannel
}

func main() {
	client, err := elariango.Initialize("api_key", true)
	if err != nil {
		log.Fatal(err)
	}

	sendSms(client)
	getCustomerState(client)
	addCustomerReminder(client)
	streamCustomerNotifications(client)
}
