package main

import (
	"log"
	"sync"

	elarian "github.com/elarianltd/go-sdk"
)

func streamCustomerNotifications(service elarian.Service) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(2)
	streamChan, errorChan := service.StreamNotifications("app_id")
	go func() {
		defer wg.Done()
		err := <-errorChan
		if err != nil {
			log.Fatalf("notification stream error %v", err)
		}

	}()
	go func() {
		for {
			res := <-streamChan
			log.Printf("response %v", res.GetReminder())
		}
	}()
}
