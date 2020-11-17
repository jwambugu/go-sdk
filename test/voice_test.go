package test

import (
	"reflect"
	"testing"

	elarian "github.com/elarianltd/go-sdk"
)

func Test_Voice(t *testing.T) {
	var opts *elarian.Options
	opts.APIKey = APIKey
	opts.AppID = AppID
	opts.OrgID = OrgID

	service, err := elarian.Initialize(opts)
	if err != nil {
		t.Fatal(err)
	}
	var customer *elarian.Customer
	customer.CustomerNumber = &elarian.CustomerNumber{
		Number:   "+254712876967",
		Provider: elarian.CustomerNumberProviderTelco,
	}

	t.Run("It should make a voice call", func(t *testing.T) {
		var channelNumber *elarian.VoiceChannelNumber
		channelNumber.Number = "+245712876967"
		channelNumber.Channel = elarian.VoiceChannelTelco

		res, err := service.MakeVoiceCall(customer, channelNumber)
		if err != nil {
			t.Error(err)
		}
		if reflect.ValueOf(res.SessionId).IsZero() {
			t.Error("Expected a session id but got none")
		}
		t.Logf("Response %v", res)
	})
}
