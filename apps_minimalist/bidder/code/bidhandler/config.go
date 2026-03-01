package bidhandler

import (
	"bidder/code/openrtb"
)

// TODO: add some type like for instance openrtb.Version
type Config struct {
	OpenRTBVersion openrtb.Version `envconfig:"BIDREQUEST_OPEN_RTB_VERSION" default:"2.5"`
}