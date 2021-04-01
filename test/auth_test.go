package test

import (
	"log"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateAuthToken(t *testing.T) {
	service, err := elarian.Connect(GetOpts())
	if err != nil {
		log.Fatal(err)
	}
	defer service.Disconnect()
	t.Run("Should Generate an Auth Token", func(t *testing.T) {
		response, err := service.GenerateAuthToken()
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.NotEqual(t, response.Token, "")
	})
}
