package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	elarian "github.com/elarianltd/go-sdk"
	"github.com/rsocket/rsocket-go/payload"
)

const (
	AppID      string = "zordTest"
	OrgID      string = "og-hv3yFs"
	customerID string = "el_cst_27f2e69b82a82176133aeea2cec28e9b"
	APIKey     string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
)

func getOpts() (*elarian.Options, *elarian.ConnectionOptions) {
	opts := &elarian.Options{
		APIKey:             APIKey,
		OrgID:              OrgID,
		AppID:              AppID,
		AllowNotifications: false,
		Log:                true,
	}
	conOpts := &elarian.ConnectionOptions{
		LifeTime:  time.Hour * 60,
		Keepalive: time.Second * 6000,
		Resumable: false,
	}
	return opts, conOpts
}

func addReminder(service elarian.Service) {
	cust := service.NewCustomer(&elarian.CreateCustomer{
		ID: customerID,
	})
	response, err := cust.AddReminder(&elarian.Reminder{Key: "KEY",
		Payload:  "i am a reminder",
		RemindAt: time.Now().Add(time.Second * 3),
		Interval: time.Duration(time.Second * 60),
	})
	if err != nil {
		log.Fatalln("Error: err")
	}
	log.Println(response)
}

func main() {
	service, err := elarian.Connect(getOpts())
	service.On("notification", func(msg payload.Payload) {
		fmt.Println("MESSAGE RECEIVED", msg)
	})
	if err != nil {
		log.Fatalln(err)
	}
	addReminder(service)
	defer service.Disconnect()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(time.Second * 120)
		wg.Done()
	}(wg)
	wg.Wait()
}
