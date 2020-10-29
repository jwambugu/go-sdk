package main

import (
	"log"

	elarian "github.com/elarianltd/go-sdk"
)

func getAuthToken(service elarian.Service) {
	response, err := service.GetAuthToken()
	if err != nil {
		log.Fatalf("could not get an auth token %v", err)
	}
	log.Printf("Auth token %v", response.Token)
}
