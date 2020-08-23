package test

import (
	"log"
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_GetAuthToken(t *testing.T) {
	service, err := elarian.Initialize(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	response, err := service.GetAuthToken(AppID)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if reflect.ValueOf(response.Lifetime).IsZero() {
		t.Errorf("Expected auth token lifetime")
	}
	if response.Token == "" {
		t.Errorf("Expected auth token %v", response.Token)
	}
}
