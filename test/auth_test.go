package test

import (
	"log"
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_GetAuthToken(t *testing.T) {
	t.Run("Should Get an Auth Token", func(t *testing.T) {
		service, err := elarian.Initialize(APIKey, OrgId)
		if err != nil {
			log.Fatal(err)
		}
		response, err := service.GetAuthToken()
		if err != nil {
			t.Errorf("Error %v", err)
		}
		if reflect.ValueOf(response.Lifetime).IsZero() {
			t.Errorf("Expected auth token lifetime")
		}
		if response.Token == "" {
			t.Errorf("Expected auth token %v", response.Token)
		}
	})
}
