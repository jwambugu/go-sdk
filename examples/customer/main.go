package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

const (
	AppID      string = "zordTest"
	OrgID      string = "og-hv3yFs"
	customerID string = "el_cst_27f2e69b82a82176133aeea2cec28e9b"
	APIKey     string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
)

func main() {

	opts := new(elarian.Options)
	opts.APIKey = APIKey
	opts.OrgID = OrgID
	opts.AppID = AppID
	opts.AllowNotifications = true
	opts.Log = true

	service, err := elarian.Connect(opts, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer service.Disconnect()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Disconnecting from Elarian")
		service.Disconnect()
		os.Exit(0)
	}()

	service.On(elarian.ElarianReminderNotification, func(notf elarian.Notification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if notification, ok := notf.(*elarian.ReminderNotification); ok {
			log.Println("NOTIFICATION _KEY", notification.Reminder.Key)
		}
	})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		service.InitializeNotificationStream()
		wg.Done()
	}(wg)

	cust := service.NewCustomer(&elarian.CreateCustomer{ID: customerID})
	response, err := cust.AddReminder(&elarian.Reminder{
		Key:      "reminderKey",
		Payload:  "i am a reminder",
		RemindAt: time.Now().Add(time.Second * 5),
	})
	if err != nil {
		log.Fatalln("Error: err", err)
	}
	log.Println(response)
	wg.Wait()
}
