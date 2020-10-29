package ussd

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func main() {
	service, err := elarian.Initialize("", "")
	if err != nil {
		log.Fatal(err)
	}
	err = service.InitializeNotificationStream("")
	if err != nil {
		log.Fatal(err)
	}

	service.AddNotificationSubscriber(
		elarian.ElarianUSSDSessionEvent,
		func(data interface{}, customer *elarian.Customer) {

		})
}
