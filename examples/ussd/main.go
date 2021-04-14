package main

import (
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
	appID      string = "test_appId"
	orgID      string = "test_org"
	customerID string = "test_customernumber"
	aPIKey     string = "test_apikey"
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
		err := service.Disconnect()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	includes := func(statusArr []elarian.PaymentStatus, status elarian.PaymentStatus) bool {
		for _, s := range statusArr {
			if status == s {
				return true
			}
		}
		return false
	}

	approveLoan := func(customer *elarian.Customer, balance float64) {
		log.Printf("Processing loan for %s", customer.CustomerNumber.Number)
		metaData, err := customer.GetMetadata()
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		name, ok := metaData["name"]
		if !ok {
			metaData["name"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("")}
		}
		repaymentDate := time.Now().Add(time.Second * 6)

		res, err := service.InitiatePayment(customer, &elarian.Paymentrequest{
			Channel: elarian.PaymentChannelNumber{},
			Cash: elarian.Cash{
				Amount:       balance,
				CurrencyCode: "KES",
			},
			DebitParty: elarian.PaymentParty{
				Purse: &elarian.Purse{PurseID: ""},
			},
		})
		if err != nil {
			log.Fatalln("Error Initiating payment: ", err)
		}
		acceptableStatus := []elarian.PaymentStatus{
			elarian.PaymentStatusQueued,
			elarian.PaymentStatusPendingConfirmation,
			elarian.PaymentStatusPendingValidation,
			elarian.PaymentStatusSuccess,
		}

		if !includes(acceptableStatus, res.Status) {
			log.Fatalf("Failed to send Kes %f to %v reason: %s \n", balance, customer.CustomerNumber.Number, res.Description)
		}
		customer.UpdateMetaData(name, &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue(strconv.FormatFloat(balance, 'f', 2, 64))})
		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		customer.SendMessage(messagingChannel, elarian.TextMessage(fmt.Sprintf("Congratulations %s, Your loan of KES %f has been approved. You are expexted to pay it back by %v", name.Value, balance, repaymentDate)))
		customer.AddReminder(
			&elarian.Reminder{
				Key:      "moni",
				RemindAt: time.Now().Add(time.Second * 2),
				Payload:  "",
			},
		)
	}

	processPayment := func(customer *elarian.Customer, notification *elarian.ReceivedPaymentNotification) {
		metaData, err := customer.GetMetadata()
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		name, ok := metaData["name"]
		if !ok {
			metaData["name"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("")}
		}
		balanceMeta, ok := metaData["balance"]
		if !ok {
			metaData["balance"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("0")}
		}
		balance, err := strconv.ParseFloat(string(balanceMeta.Value.String()), 64)
		if err != nil {
			log.Fatalln("Error converting balance", err)
		}
		newBalance := balance - notification.Value.Amount
		customer.UpdateMetaData(
			name,
			&elarian.Metadata{Key: "balance", Value: elarian.StringDataValue(strconv.FormatFloat(newBalance, 'f', 2, 64))},
		)

		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		if newBalance <= 0 {
			customer.CancelReminder("moni")
			customer.SendMessage(
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Thank you for your payment %s, your loan has been fully repaid!!", name)),
			)
			customer.DeleteMetaData("name", "strike", "balance", "screen")
			return
		}
		customer.SendMessage(
			messagingChannel,
			elarian.TextMessage(fmt.Sprintf("Hey %s! \n Thank you for your payment, but you still owe me KES ${newBalance}", name)),
		)
	}

	processReminder := func(customer *elarian.Customer, notification *elarian.ReminderNotification) {
		metaData, err := customer.GetMetadata()
		if err != nil {
			log.Fatalln("Error Fetching Metadata", err)
		}
		strike, ok := metaData["strike"]
		if !ok {
			metaData["strike"] = &elarian.Metadata{Key: "strike", Value: elarian.StringDataValue("1")}
		}
		name, ok := metaData["name"]
		if !ok {
			metaData["name"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("")}
		}
		balance, ok := metaData["balance"]
		if !ok {
			metaData["balance"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("0")}
		}

		strikeValue, err := strconv.Atoi(strike.Value.String())
		if err != nil {
			log.Fatalln("Error parsing strike value")
		}

		messagingChannel := &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}
		if strikeValue == 1 {
			customer.SendMessage(
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Hey %s, this is a friendly reminder to pay back my KES %s", name, balance)),
			)
		}
		if strikeValue == 2 {
			customer.SendMessage(
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Hey %s,you still need to pay back my KES %s", name, balance)))
		}
		if strikeValue > 2 {
			customer.SendMessage(
				messagingChannel,
				elarian.TextMessage(fmt.Sprintf("Yo %s, !!! you need to pay back my KES %s", name, balance)),
			)
		}

		strikeValue++
		customer.UpdateMetaData(balance, name, &elarian.Metadata{Key: "strike", Value: elarian.StringDataValue(strconv.Itoa(strikeValue))})
		customer.AddReminder(&elarian.Reminder{Key: "moni", RemindAt: time.Now().Add(time.Minute * 1), Payload: ""})
	}

	processUssd := func(customer *elarian.Customer, notification *elarian.UssdSessionNotification, appdata *elarian.Appdata, cb elarian.NotificationCallBack) {
		fmt.Printf("Processing USSD from %v \n", customer.CustomerNumber.Number)

		// Create an abitrary  struct to hold our appdata
		appData := struct {
			Screen string `json:"screen,omitempty"`
		}{Screen: "home"}

		if appdata.Value.String() != "" {
			err := json.Unmarshal([]byte(appdata.Value.String()), &appData)
			if err != nil {
				log.Fatalln("Error unmarshalling app data", err)
			}
		}

		// get a customer's metadata
		customerData, err := customer.GetMetadata()
		if err != nil {
			log.Fatalln("Error getting customer data: ", err)
		}
		nameVal, ok := customerData["name"]
		if !ok {
			customerData["name"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("")}
		}
		balanceVal, ok := customerData["balance"]
		if !ok {
			customerData["balance"] = &elarian.Metadata{Key: "balance", Value: elarian.StringDataValue("0")}
		}

		name := nameVal.Value
		balance, err := strconv.ParseFloat(balanceVal.Value.String(), 64)
		if err != nil {
			log.Fatalln("Error converting balance", err)
		}

		menu := &elarian.UssdMenu{Text: "", IsTerminal: false}
		nextScreen := appData.Screen
		if appData.Screen == "home" && notification.Input == "" {
			if notification.Input == "1" {
				nextScreen = "request-name"
			}
			if notification.Input == "2" {
				nextScreen = "quit"
			}
			if name.String() != "" {
				nextScreen = "info"
			}
		}

		switch nextScreen {
		case "quit":
			menu.Text = "Happy Coding"
			menu.IsTerminal = true
			nextScreen = "home"
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})

		case "info":
			menu.Text = fmt.Sprintf("Hey %s, ", name)
			if balance > 0 {
				menu.Text += fmt.Sprintf("you still owe me KES %f !", balance)
			} else {
				menu.Text += "you have repaid your loan, good for you !"
			}
			menu.IsTerminal = false
			nextScreen = "home"
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})

		case "request-name":
			menu.Text = "Alright, what is your name?"
			nextScreen = "request-amount"
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})
		case "approve-amount":
			balance, err := strconv.ParseFloat(notification.Input, 64)
			if err != nil {
				log.Fatalln("Error converting balance", err)
			}
			menu.Text = fmt.Sprintf("Awesome! %s we are reviewing your application and will be in touch shortly \n Have a lovely day!", name)
			menu.IsTerminal = true
			nextScreen = "home"
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})
			approveLoan(customer, balance)
		case "home":
			menu.Text = "Welcome to MoniMoni!\n1. Apply for loan\n2. Quit"
			menu.IsTerminal = false
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})
		default:
			menu.Text = "Welcome to MoniMoni!\n1. Apply for loan\n2. Quit"
			menu.IsTerminal = false
			appData.Screen = nextScreen
			val, _ := json.Marshal(appData)
			cb(menu, &elarian.Appdata{Value: elarian.StringDataValue(string(val))})
		}
		customer.UpdateMetaData(nameVal,
			&elarian.Metadata{Key: "balance", Value: elarian.StringDataValue(strconv.FormatFloat(balance, 'f', 2, 64))})
	}

	service.On(elarian.ElarianReminderNotification, func(notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if reminderNotification, ok := notf.(*elarian.ReminderNotification); ok {
			processReminder(customer, reminderNotification)
		}
	})

	service.On(elarian.ElarianReceivedUssdSessionNotification, func(notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
		if messageNotification, ok := notf.(*elarian.UssdSessionNotification); ok {
			processUssd(customer, messageNotification, appData, cb)
		}
	})
	service.On(elarian.ElarianReceivedPaymentNotification, func(notf elarian.IsNotification, appData *elarian.Appdata, customer *elarian.Customer, cb elarian.NotificationCallBack) {
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

	cust := service.NewCustomer(&elarian.CreateCustomer{ID: customerID})
	response, err := cust.AddReminder(&elarian.Reminder{Key: "KEY",
		Payload:  "i am a reminder",
		RemindAt: time.Now().Add(time.Second * 5),
	})
	if err != nil {
		log.Fatalln("Error: err", err)
	}
	log.Println(response)
	wg.Wait()
	defer service.Disconnect()
}
