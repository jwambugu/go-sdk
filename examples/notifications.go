package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func streamCustomerNotifications(service elarian.Service) {
	streamChan, errorChan := service.StreamNotifications("app_id")
	err := <-errorChan
	if err != nil {
		log.Fatalf("notification stream error %v", err)
	}
	for {
		res := <-streamChan
		log.Printf("response %v", res.GetReminder())
	}
}
