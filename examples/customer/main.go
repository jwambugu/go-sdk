package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

const (
	appID      string = "zordTest"
	orgID      string = "og-hv3yFs"
	customerID string = "el_cst_27f2e69b82a82176133aeea2cec28e9b"
	aPIKey     string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
)

func main() {
	opts := &elarian.Options{
		APIKey:             aPIKey,
		OrgID:              orgID,
		AppID:              appID,
		AllowNotifications: true,
		Log:                true,
	}

	service, err := elarian.Connect(opts, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer service.Disconnect()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		_, err := service.CancelCustomerReminder(context.Background(), elarian.CustomerID(customerID), "reminderKey")
		if err != nil {
			log.Fatalln(err)
		}
		err = service.Disconnect()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		err := <-service.InitializeNotificationStream()
		if err != nil {
			log.Println("Notification error: ", err)
		}
		wg.Done()
	}(wg)

	service.On(elarian.ElarianReminderNotification, func(service elarian.Elarian, notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if notification, ok := notf.(*elarian.ReminderNotification); ok {
			log.Println("NOTIFICATION_KEY", notification.Reminder.Key)
			cb(nil, nil)
		}
	})

	cust := service.NewCustomer(&elarian.CreateCustomer{ID: customerID})
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()
	response, err := cust.AddReminder(ctx, &elarian.Reminder{
		Key:      "reminderKey",
		Payload:  "i am a reminder",
		RemindAt: time.Now().Add(time.Second * 5),
		Interval: time.Second * 60,
	})
	if err != nil {
		log.Fatalln("Error: err", err)
	}
	log.Println(response)
	wg.Wait()
}
