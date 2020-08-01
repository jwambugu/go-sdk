package main

import (
	"context"
	"log"
	"time"

	"github.com/elarian/elariango"
	elarian "github.com/elarian/elariango/com_elarian_hera_proto"
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
				Number:   "+254712345678",
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
		AppId: "app-id",
		Customer: &elarian.GetCustomerStateRequest_CustomerId{
			CustomerId: "el_cst_35ff1eb3r448652dv55556fvff",
		},
	})
	if err != nil {
		log.Fatalf("could not get customer state %v", err)
	}
	log.Printf("customer state %v", res.GetData())
}

// AddCustomerReminder func
func AddCustomerReminder(client elarian.GrpcWebServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.AddCustomerReminder(ctx, &elarian.AddCustomerReminderRequest{
		AppId: "some_app",
		Customer: &elarian.AddCustomerReminderRequest_CustomerId{
			CustomerId: "",
		},
		Reminder: &elarian.CustomerReminder{
			ProductId: "some_product_id",
			Key:       "some_key",
			Payload: &wrapperspb.StringValue{
				Value: "i am some payload",
			},
		},
	})
	if err != nil {
		log.Fatalf("could not set a reminder %v", err)
	}
	log.Printf("response %v", res.GetDescription())
}

func main() {
	client, err := elariango.Initialize("api_key", true)
	if err != nil {
		log.Fatal(err)
	}
	sendSms(client)
	getCustomerState(client)
}
