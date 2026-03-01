package auction

import (
	"apps_minimalist/bidder/code/openrtb"
)

type Item struct {
	ID []byte
}

type Request struct {
	ID				[]byte
	Item			[]Item
	OpenRTBVersion	openrtb.Version
}