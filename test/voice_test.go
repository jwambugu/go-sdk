package test

import (
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Voice(t *testing.T) {
	service, err := elarian.Initialize(APIKey, OrgId)
	if err != nil {
		t.Fatal(err)
	}
	var customer elarian.Customer
	customer.CustomerNumber = elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CUSTOMER_NUMBER_PROVIDER_TELCO,
	}

	t.Run("It should make a voice call", func(t *testing.T) {
		var request elarian.VoiceCallRequest
		request.AppId = AppId

		request.Channel = elarian.VoiceChannelNumber{
			Channel: elarian.VOICE_CHANNEL_TELCO,
			Number:  "+245712876967",
		}

		res, err := service.MakeVoiceCall(&customer, &request)
		if err != nil {
			t.Error(err)
		}
		if reflect.ValueOf(res.SessionId).IsZero() {
			t.Error("Expected a session id but got none")
		}
		t.Logf("Response %v", res)
	})
}
