package test

import (
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

const (
	AppID      string = "zordTest"
	OrgID      string = "og-hv3yFs"
	customerID string = "el_cst_27f2e69b82a82176133aeea2cec28e9b"
	APIKey     string = "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145"
)

func GetOpts() (*elarian.Options, *elarian.ConnectionOptions) {
	opts := &elarian.Options{
		APIKey:             APIKey,
		OrgID:              OrgID,
		AppID:              AppID,
		AllowNotifications: false,
		Log:                true,
	}
	conOpts := &elarian.ConnectionOptions{
		LifeTime:  time.Hour * 60,
		Keepalive: time.Second * 6000,
		Resumable: false,
	}
	return opts, conOpts
}
