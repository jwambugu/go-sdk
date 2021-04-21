package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
	includes := func(statusArr []elarian.PaymentStatus, status elarian.PaymentStatus) bool {
		for _, s := range statusArr {
			if status == s {
				return true
			}
		}
		return false
	}

	approveLoan := func(service elarian.Elarian, customer *elarian.Customer, balance float64) {
		log.Printf("Processing loan for %s", customer.CustomerNumber.Number)
		metaData, err := customer.GetMetadata(context.Background())
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		name, ok := metaData["name"]
		if !ok {
			name = &elarian.Metadata{Key: "balance", Value: ""}
			metaData["name"] = name
		}
		repaymentDate := time.Now().Add(time.Second * 6)

		res, err := service.InitiatePayment(
			context.Background(),
			&elarian.PaymentCounterParty{
				DebitParty: &elarian.Purse{PurseID: "prs-PZSvFO"},
				CreditParty: &elarian.CustomerPaymentParty{
					CustomerNumber: customer.CustomerNumber,
					ChannelNumber: &elarian.PaymentChannelNumber{
						Channel: elarian.PaymentChannelCellular,
						Number:  "525900",
					},
				},
			},
			&elarian.Cash{
				CurrencyCode: "KES",
				Amount:       balance,
			},
		)
		if err != nil {
			log.Fatalln("Error Initiating payment: ", err)
		}
		acceptableStatus := []elarian.PaymentStatus{
			elarian.PaymentStatusQueued,
			elarian.PaymentStatusPendingConfirmation,
			elarian.PaymentStatusPendingValidation,
			elarian.PaymentStatusSuccess,
		}

		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		if !includes(acceptableStatus, res.Status) {
			customer.SendMessage(context.Background(), messagingChannel, elarian.TextMessage("Processing your loan failed please try again"))
			log.Printf("Failed to send Kes %f to %v reason: %s \n", balance, customer.CustomerNumber.Number, res.Description)
			return
		}

		customer.UpdateMetaData(context.Background(), name, &elarian.Metadata{Key: "balance", Value: strconv.FormatFloat(balance, 'f', 2, 64)})
		customer.SendMessage(context.Background(), messagingChannel, elarian.TextMessage(fmt.Sprintf("Congratulations %s, Your loan of KES %f has been approved. You are expexted to pay it back by %v", name.Value, balance, repaymentDate)))
		customer.AddReminder(
			context.Background(),
			&elarian.Reminder{
				Key:      "moni",
				RemindAt: time.Now().Add(time.Second * 2),
				Payload:  "",
			},
		)
	}

	processPayment := func(customer *elarian.Customer, notification *elarian.ReceivedPaymentNotification) {
		metaData, err := customer.GetMetadata(context.Background())
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		name, ok := metaData["name"]
		if !ok {
			name = &elarian.Metadata{Key: "balance", Value: ""}
			metaData["name"] = name
		}
		balanceMeta, ok := metaData["balance"]
		if !ok {
			balanceMeta = &elarian.Metadata{Key: "balance", Value: "0"}
			metaData["balance"] = balanceMeta
		}
		balance, err := strconv.ParseFloat(string(balanceMeta.Value), 64)
		if err != nil {
			log.Fatalln("Error converting balance", err)
		}
		newBalance := balance - notification.Value.Amount
		customer.UpdateMetaData(
			context.Background(),
			name,
			&elarian.Metadata{Key: "balance", Value: strconv.FormatFloat(newBalance, 'f', 2, 64)},
		)

		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		if newBalance <= 0 {
			customer.CancelReminder(context.Background(), "moni")
			customer.SendMessage(
				context.Background(),
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Thank you for your payment %s, your loan has been fully repaid!!", name.Value)),
			)
			customer.DeleteMetaData(context.Background(), "name", "strike", "balance", "screen")
			return
		}
		customer.SendMessage(
			context.Background(),
			messagingChannel,
			elarian.TextMessage(fmt.Sprintf("Hey %s! \n Thank you for your payment, but you still owe me KES ${newBalance}", name)),
		)
	}

	processReminder := func(customer *elarian.Customer, notification *elarian.ReminderNotification) {
		metaData, err := customer.GetMetadata(context.Background())
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		strike, ok := metaData["strike"]
		if !ok {
			strike = &elarian.Metadata{Key: "strike", Value: "1"}
			metaData["strike"] = strike
		}
		name, ok := metaData["name"]
		if !ok {
			name = &elarian.Metadata{Key: "balance", Value: ""}
			metaData["name"] = name
		}
		balance, ok := metaData["balance"]
		if !ok {
			balance = &elarian.Metadata{Key: "balance", Value: "0"}
			metaData["balance"] = balance
		}
		strikeValue, err := strconv.Atoi(strike.Value)
		if err != nil {
			log.Fatalln("Error parsing strike value")
		}

		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		if strikeValue == 1 {
			customer.SendMessage(
				context.Background(),
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Hey %s, this is a friendly reminder to pay back my KES %s", name, balance)),
			)
		}
		if strikeValue == 2 {
			customer.SendMessage(
				context.Background(),
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Hey %s,you still need to pay back my KES %s", name, balance)))
		}
		if strikeValue > 2 {
			customer.SendMessage(
				context.Background(),
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Yo %s, !!! you need to pay back my KES %s", name, balance)),
			)
		}
		strikeValue++
		customer.UpdateMetaData(context.Background(), balance, name, &elarian.Metadata{Key: "strike", Value: strconv.Itoa(strikeValue)})
		customer.AddReminder(context.Background(), &elarian.Reminder{Key: "moni", RemindAt: time.Now().Add(time.Minute * 1), Payload: ""})
	}

	processUssd := func(service elarian.Elarian, customer *elarian.Customer, notification *elarian.UssdSessionNotification, appdata *elarian.Appdata, cb elarian.NotificationCallBack) {
		fmt.Printf("Processing USSD from %v %s \n", customer.CustomerNumber.Number, notification.SessionID)

		// Create an abitrary  struct to hold our appdata
		appData := struct {
			State     string `json:"state"`
			SessionID string `json:"sessionId"`
		}{}

		if appdata.Value != "" {
			err := json.Unmarshal([]byte(appdata.Value), &appData)
			if err != nil {
				log.Fatalln("Error unmarshalling app data", err)
			}
		}

		if notification.SessionID != appData.SessionID {
			appData.State = "initial"
			appData.SessionID = notification.SessionID
		}

		// get a customer's metadata
		customerData, err := customer.GetMetadata(context.Background())
		if err != nil {
			log.Fatalln("Error getting customer data: ", err)
		}

		nameVal, ok := customerData["name"]
		if !ok {
			nameVal = &elarian.Metadata{Key: "name", Value: ""}
			customerData["name"] = nameVal
		}
		balanceVal, ok := customerData["balance"]
		if !ok {
			balanceVal = &elarian.Metadata{Key: "balance", Value: "0"}
			customerData["balance"] = balanceVal
		}

		name := nameVal.Value
		balance, err := strconv.ParseFloat(balanceVal.Value, 64)
		if err != nil {
			log.Fatalln("Error converting balance", err)
		}

		menu := &elarian.UssdMenu{Text: "", IsTerminal: false}
		state := appData.State
		switch state {
		case "initial":
			appData.State = "home"
			menu.Text = "Welcome to MoniMoni!\n1. Apply for loan\n2. Quit"
			val, _ := json.Marshal(&appData)
			cb(menu, &elarian.Appdata{Value: string(val)})
			return
		case "home":
			if notification.Input == "1" {
				if name == "" {
					appData.State = "request-name"
					menu.Text = "Alright, what is your name?"
					val, _ := json.Marshal(&appData)
					cb(menu, &elarian.Appdata{Value: string(val)})
					return
				}
				if balance > 0 {
					appData.State = "request-amount"
					menu.Text += fmt.Sprintf("Hey %s, you still owe KES %f !", name, balance)
					val, _ := json.Marshal(&appData)
					cb(menu, &elarian.Appdata{Value: string(val)})
					return
				}
				menu.Text = fmt.Sprintf("Okay %s, how much do you need?", name)
				appData.State = "request-amount"
				val, _ := json.Marshal(&appData)
				cb(menu, &elarian.Appdata{Value: string(val)})
				return
			}
			if notification.Input == "2" {
				menu.Text = "Have a great day"
				menu.IsTerminal = true
				cb(menu, nil)
				return
			}
		case "request-name":
			name = notification.Input
			menu.Text = fmt.Sprintf("Okay %s, how much do you need?", name)
			appData.State = "request-amount"
			val, _ := json.Marshal(&appData)
			cb(menu, &elarian.Appdata{Value: string(val)})
			customer.UpdateMetaData(context.Background(), &elarian.Metadata{Key: "name", Value: name})
			return
		case "request-amount":
			appData.State = "approve-amount"
			balance, err := strconv.ParseFloat(notification.Input, 64)
			if err != nil {
				menu.Text = "Incorrect amount, please try again"
				menu.IsTerminal = true
				val, _ := json.Marshal(&appData)
				cb(menu, &elarian.Appdata{Value: string(val)})
				log.Println("Error converting balance", err)
				return
			}
			menu.Text = fmt.Sprintf("Awesome! %s we are reviewing your application and will be in touch shortly \n Have a lovely day!", name)
			menu.IsTerminal = true
			val, _ := json.Marshal(&appData)
			cb(menu, &elarian.Appdata{Value: string(val)})
			approveLoan(service, customer, balance)
			return
		case "approve-amount":
			appData.State = "approve-amount"
			menu.Text = fmt.Sprintf("Hi! %s your loan was approved. \n Please contact our support team for further inquiries!", name)
			menu.IsTerminal = true
			val, _ := json.Marshal(&appData)
			cb(menu, &elarian.Appdata{Value: string(val)})
			return
		default:
			appData.State = "home"
			menu.Text = "Welcome to MoniMoni!\n1. Apply for loan\n2. Quit"
			val, _ := json.Marshal(&appData)
			cb(menu, &elarian.Appdata{Value: string(val)})
			return
		}
	}

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
		err := service.Disconnect()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	service.On(elarian.ElarianReminderNotification, func(service elarian.Elarian, notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if reminderNotification, ok := notf.(*elarian.ReminderNotification); ok {
			processReminder(customer, reminderNotification)
		}
	})

	service.On(elarian.ElarianReceivedUssdSessionNotification, func(service elarian.Elarian, notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if ussdNotification, ok := notf.(*elarian.UssdSessionNotification); ok {
			processUssd(service, customer, ussdNotification, appData, cb)
		}
	})

	service.On(elarian.ElarianReceivedPaymentNotification, func(service elarian.Elarian, notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if paymentNotification, ok := notf.(*elarian.ReceivedPaymentNotification); ok {
			processPayment(customer, paymentNotification)
		}
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		err := <-service.InitializeNotificationStream()
		if err != nil {
			log.Println("Notification err", err)
		}
		wg.Done()
	}(wg)
	wg.Wait()
}
