package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

const (
	appID  string = "zordTest"
	orgID  string = "og-hv3yFs"
	aPIKey string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
)

func main() {
	var (
		custNumber *elarian.CustomerNumber
		channel    *elarian.MessagingChannelNumber
		opts       *elarian.Options
	)

	opts = &elarian.Options{
		APIKey:             aPIKey,
		OrgID:              orgID,
		AppID:              appID,
		AllowNotifications: true,
		Log:                true,
	}
	service, err := elarian.Connect(opts, nil)
	if err != nil {
		log.Fatalf("Error Initializing Elarian: %v \n", err)
	}
	defer service.Disconnect()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Disconnecting From Elarian")
		service.Disconnect()
		os.Exit(0)
	}()

	custNumber = &elarian.CustomerNumber{Number: "+254708752502", Provider: elarian.CustomerNumberProviderCellular}
	channel = &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()
	response, err := service.SendMessage(ctx, custNumber, channel, elarian.TextMessage("Hello world from the go sdk"))
	if err != nil {
		log.Fatalf("Message not send %v \n", err.Error())
	}
	log.Printf("Status %d Description %s \n customerID %s \n", response.Status, response.Description, response.CustomerID)
}
