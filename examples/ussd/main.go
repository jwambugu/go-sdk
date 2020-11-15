package main

import (
	"encoding/json"
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

type (
	AppData struct {
		Name  string `json:"name,omitempty"`
		State string `json:"state,omitempty"`
	}
)

func errorHandler(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ussdApp(srv elarian.Service, cust *elarian.Customer, data interface{}) {
	sessionData, ok := data.(elarian.UssdSessionNotification)
	if !ok {
		log.Fatalln("No session data")
	}
	meta, err := cust.LeaseMetaData("awesomeNameSurvey")
	errorHandler(err)

	var menu *elarian.UssdMenu
	menu.IsTerminal = true
	menu.Text = ""

	metadata := &AppData{}
	metadata.State = "newbie"

	if len(meta.Value.GetBytesVal()) > 0 {
		err = json.Unmarshal(meta.Value.GetBytesVal(), metadata)
		errorHandler(err)
	}

	switch metadata.State {
	case "veteran":
		if metadata.Name != "" {
			menu.Text = "Welcome back " + metadata.Name + "! What's your new name?"
			menu.IsTerminal = false
			return
		}
		metadata.Name = sessionData.Input
		menu.Text = "Thank you for trying Elarian" + metadata.Name
		menu.IsTerminal = true

		body := &elarian.MessageBody{
			Text: "Hey" + metadata.Name + "Thank you for trying our Elarian",
		}
		channelNumber := &elarian.MessagingChannelNumber{
			Number:  "Elarian",
			Channel: elarian.MESSAGING_CHANNEL_SMS,
		}
		_, err := cust.SendMessage(channelNumber, body)
		errorHandler(err)
	case "newbie":
	default:
		menu.Text = "Hey there welcome to Elarian! What's your name?"
		menu.IsTerminal = false
		metadata.State = "veteran"
	}
	_, err = cust.UpdateMetaData(
		map[string]string{
			"name":  metadata.Name,
			"state": metadata.State,
		},
	)
	errorHandler(err)
	_, err = srv.ReplyToUssdSession(
		sessionData.SessionId,
		menu)
	errorHandler(err)
}

func main() {
	service, err := elarian.Initialize(&elarian.Options{})
	errorHandler(err)
	err = service.AddNotificationSubscriber(
		elarian.ELARIAN_USSD_SESSION_NOTIFICATION,
		ussdApp,
	)
	errorHandler(err)
}
