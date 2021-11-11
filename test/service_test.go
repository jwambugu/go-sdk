package test

import (
	elarian "github.com/elarianltd/go-sdk"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_NewService(t *testing.T) {
	t.Run("Should connect to elarian successfully", func(t *testing.T) {
		service, err := elarian.Connect(GetOpts())
		if err != nil {
			log.Fatal(err)
		}

		defer service.Disconnect()

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}
