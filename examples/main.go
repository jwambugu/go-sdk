// Package main implements a the simple examples that demonstrate hot to use the elarian SDK.
package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func main() {
	service, err := elarian.Initialize("api_key", true)
	if err != nil {
		log.Fatal(err)
	}
	sendMessage(service)
	getCustomerState(service)
	addCustomerReminder(service)
}
