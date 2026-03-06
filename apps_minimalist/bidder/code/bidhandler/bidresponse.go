package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/openrtb"
)

func buildResponse(r *auction.Response, pd *persistentData) []byte {
	switch r.Request.OpenRTBVersion {
	case openrtb.OpenRTB2_5:
		return buildResponse2(r, pd)
	case openrtb.OpenRTB3_0:
		return buildResponse3(r, pd)
	default:
		return nil
	}
}

func buildResponse2(r *auction.Response, pd *persistentData) []byte {
	pd.byteResponse = pd.byteResponse[:0]
	pd.byteResponse = append(pd.byteResponse, `{"id": "1"}`...)

	return pd.byteResponse
}

func buildResponse3(r *auction.Response, pd *persistentData) []byte {
	// TODO
	return nil
}
