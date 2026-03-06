package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/openrtb"
	"apps_minimalist/bidder/code/price"
	"net/url"
	"strconv"

	"gvisor.dev/gvisor/pkg/gohacks"
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

	// Request ID
	requestID := gohacks.StringFromImmutableBytes(r.Request.ID)
	pd.byteResponse = append(pd.byteResponse, `{"id":`...)
	pd.byteResponse = strconv.AppendQuote(pd.byteResponse, requestID)
	
	// Seat ID
	pd.byteResponse = append(pd.byteResponse, `,"seatbid":[{"seat":`...)
	pd.byteResponse = append(pd.byteResponse, `TESTSEAT`...)  // TMP

	// Bid - ID
	pd.byteResponse = append(pd.byteResponse, `,"bid":[{"id":"`...)  
	pd.byteResponse = pd.ksuidSequence.Get().Append(pd.byteResponse)
	// Price
	pd.byteResponse = append(pd.byteResponse, `","price":`...)
	pd.byteResponse = strconv.AppendFloat(pd.byteResponse, price.ToFloat(r.Price), 'f', -1, 64)
	// BURL
	pd.byteResponse = append(pd.byteResponse, `,"burl":"`...)  
	pd.byteResponse = append(pd.byteResponse, `https://this.is.test.burl.com/`...)
	escapedRID := url.PathEscape(requestID)
	// escapedCID := url.PathEscape(r.Campaign.HexID)
	escapedCID := "testcampaignid"
	pd.byteResponse = append(pd.byteResponse, escapedRID...)
	pd.byteResponse = append(pd.byteResponse, "/"...)
	pd.byteResponse = append(pd.byteResponse, escapedCID...)
	pd.byteResponse = append(pd.byteResponse, `/${OPENRTB_PRICE}",`...)
	// pd.byteResponse = append(pd.byteResponse, ``...)  

	pd.byteResponse = append(pd.byteResponse, `}]}]}`...)

	return pd.byteResponse
}

func buildResponse3(r *auction.Response, pd *persistentData) []byte {
	// TODO
	return nil
}
