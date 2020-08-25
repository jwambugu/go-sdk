// Package main implements a the simple examples that demonstrate hot to use the elarian SDK.
package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

type (
	receiver func(service elarian.Service)
)

func main() {
	service, err := elarian.Initialize("api_key")
	if err != nil {
		log.Fatal(err)
	}
	receivers := []receiver{
		sendMessage,
		getCustomerState,
		getAuthToken,
		streamCustomerNotifications,
		addCustomerReminder,
	}
	for _, receiver := range receivers {
		receiver(service)
	}
}
