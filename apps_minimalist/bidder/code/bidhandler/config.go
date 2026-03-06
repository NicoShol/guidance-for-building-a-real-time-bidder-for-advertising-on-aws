package bidhandler

import (
	"apps_minimalist/bidder/code/openrtb"
	"time"
)

// TODO: add some type like for instance openrtb.Version
type Config struct {
	Timeout        time.Duration   `envconfig:"BIDREQUEST_TIMEOUT" required:"true"`
	TimeoutStatus  int             `envconfig:"BIDREQUEST_TIMEOUT_STATUS" required:"true"`
	OpenRTBVersion openrtb.Version `envconfig:"BIDREQUEST_OPEN_RTB_VERSION" default:"2.5"`
}