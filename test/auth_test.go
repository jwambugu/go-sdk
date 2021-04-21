package test

import (
	"context"
	"log"
	"testing"
	"time"

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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
		defer cancel()
		response, err := service.GenerateAuthToken(ctx)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		assert.NotNil(t, response)
		assert.NotEqual(t, response.Token, "")
	})
}
