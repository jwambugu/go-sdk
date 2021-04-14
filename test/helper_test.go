package test

import (
	"time"

	elarian "github.com/elarianltd/go-sdk"
)

const (
	customerID string = "el_cst_27f2e69b82a82176133aeea2cec28e9b"
)

func GetOpts() (*elarian.Options, *elarian.ConnectionOptions) {
	opts := &elarian.Options{
		APIKey:             "el_api_key_6b3ff181a2d5cf91f62d2133a67a25b3070d2d7305eba70288417b3ab9ebd145",
		OrgID:              "og-hv3yFs",
		AppID:              "zordTest",
		AllowNotifications: false,
		Log:                true,
	}
	conOpts := &elarian.ConnectionOptions{
		LifeTime:  time.Second * 1,
		Keepalive: time.Second * 2,
		Resumable: true,
	}
	return opts, conOpts
}
