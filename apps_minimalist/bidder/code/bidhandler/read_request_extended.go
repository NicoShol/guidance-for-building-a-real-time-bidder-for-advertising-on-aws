package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/openrtb"
	"strconv"

	"emperror.dev/errors"
)

// parseBidRequestExtended parses a complete OpenRTB 2.5 bid request including
// all device and user fields. This is a comprehensive version that extracts
// everything from the specification.
func parseBidRequestExtended_2_5(byteRequest []byte, pd *persistentData) (*auction.ExtendedRequest, error) {
	v, err := pd.parser.ParseBytes(byteRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error while parsing request")
	}

	// -----------------------------------------------------------------
	// 1. PARSE TOP-LEVEL FIELDS
	// -----------------------------------------------------------------
	ID := v.GetStringBytes("id")
	if len(ID) == 0 {
		return nil, errors.New("empty request ID")
	}

	// -----------------------------------------------------------------
	// 2. PARSE DEVICE OBJECT
	// -----------------------------------------------------------------
	deviceObj := v.Get("device")
	device := auction.Device{}

	if deviceObj != nil {
		device.IFA = deviceObj.GetStringBytes("ifa")
		device.UA = deviceObj.GetStringBytes("ua")
		device.IP = deviceObj.GetStringBytes("ip")
		device.IPv6 = deviceObj.GetStringBytes("ipv6")
		device.DeviceType = deviceObj.GetInt("devicetype")
		device.Make = deviceObj.GetStringBytes("make")
		device.Model = deviceObj.GetStringBytes("model")
		device.OS = deviceObj.GetStringBytes("os")
		device.OSV = deviceObj.GetStringBytes("osv")
		device.HWV = deviceObj.GetStringBytes("hwv")
		device.W = deviceObj.GetInt("w")
		device.H = deviceObj.GetInt("h")
		device.PPI = deviceObj.GetInt("ppi")
		device.PxRatio = deviceObj.GetFloat64("pxratio")
		device.JS = deviceObj.GetInt("js")
		device.Language = deviceObj.GetStringBytes("language")
		device.Carrier = deviceObj.GetStringBytes("carrier")
		device.MCCMNC = deviceObj.GetStringBytes("mccmnc")
		device.ConnectionType = deviceObj.GetInt("connectiontype")
		device.DNT = deviceObj.GetInt("dnt")
		device.Lmt = deviceObj.GetInt("lmt")

		// Parse nested Geo object
		geoObj := deviceObj.Get("geo")
		if geoObj != nil {
			device.Geo = auction.Geo{
				Lat:       geoObj.GetFloat64("lat"),
				Lon:       geoObj.GetFloat64("lon"),
				Type:      geoObj.GetInt("type"),
				Country:   geoObj.GetStringBytes("country"),
				Region:    geoObj.GetStringBytes("region"),
				Metro:     geoObj.GetStringBytes("metro"),
				City:      geoObj.GetStringBytes("city"),
				Zip:       geoObj.GetStringBytes("zip"),
				UTCOffset: geoObj.GetInt("utcoffset"),
			}
		}
	}

	// -----------------------------------------------------------------
	// 3. PARSE USER OBJECT
	// -----------------------------------------------------------------
	userObj := v.Get("user")
	user := auction.User{}

	if userObj != nil {
		user.ID = userObj.GetStringBytes("id")
		user.BuyerUID = userObj.GetStringBytes("buyeruid")
		user.YOB = userObj.GetInt("yob")
		user.Gender = userObj.GetStringBytes("gender")
		user.Keywords = userObj.GetStringBytes("keywords")

		// Parse user Geo object (home location, different from device geo)
		userGeoObj := userObj.Get("geo")
		if userGeoObj != nil {
			user.Geo = auction.Geo{
				Lat:       userGeoObj.GetFloat64("lat"),
				Lon:       userGeoObj.GetFloat64("lon"),
				Type:      userGeoObj.GetInt("type"),
				Country:   userGeoObj.GetStringBytes("country"),
				Region:    userGeoObj.GetStringBytes("region"),
				Metro:     userGeoObj.GetStringBytes("metro"),
				City:      userGeoObj.GetStringBytes("city"),
				Zip:       userGeoObj.GetStringBytes("zip"),
				UTCOffset: userGeoObj.GetInt("utcoffset"),
			}
		}

		// Parse user data segments (third-party audience data)
		dataArray := userObj.GetArray("data")

		// Use temporary slices for building up data
		var userData []auction.UserData

		for _, dataObj := range dataArray {
			data := auction.UserData{
				ID:   dataObj.GetStringBytes("id"),
				Name: dataObj.GetStringBytes("name"),
			}

			// Parse segments within each data provider
			segmentArray := dataObj.GetArray("segment")
			var segments []auction.UserSegment

			for _, segObj := range segmentArray {
				segment := auction.UserSegment{
					ID:    segObj.GetStringBytes("id"),
					Name:  segObj.GetStringBytes("name"),
					Value: segObj.GetStringBytes("value"),
				}
				segments = append(segments, segment)
			}

			data.Segments = segments
			userData = append(userData, data)
		}

		user.Data = userData
	}

	// -----------------------------------------------------------------
	// 4. PARSE IMPRESSION OBJECTS
	// -----------------------------------------------------------------
	impArray := v.GetArray("imp")
	if len(impArray) == 0 {
		return nil, errors.New("bidrequest must contain at least one imp object")
	}

	// Reset items slice, keep capacity
	pd.auctionRequest.Items = pd.auctionRequest.Items[:0]

	for _, impObj := range impArray {
		item := auction.ExtendedItem{
			ID:          impObj.GetStringBytes("id"),
			BidFloor:    impObj.GetFloat64("bidfloor"),
			BidFloorCur: impObj.GetStringBytes("bidfloorcur"),
			Secure:      impObj.GetInt("secure"),
			TagID:       impObj.GetStringBytes("tagid"),
			Exp:         impObj.GetInt("exp"),
		}

		// Parse banner object if present
		bannerObj := impObj.Get("banner")
		if bannerObj != nil {
			item.BannerWidth = bannerObj.GetInt("w")
			item.BannerHeight = bannerObj.GetInt("h")

			// Parse banner mimes array
			mimesArray := bannerObj.GetArray("mimes")
			var mimes [][]byte

			for _, mimeVal := range mimesArray {
				mimes = append(mimes, mimeVal.GetStringBytes())
			}

			item.BannerMimes = mimes
		}

		pd.auctionRequest.Items = append(pd.auctionRequest.Items, item)
	}

	// -----------------------------------------------------------------
	// 5. BUILD EXTENDED REQUEST
	// -----------------------------------------------------------------
	pd.auctionRequest.ID = ID
	pd.auctionRequest.Device = device
	pd.auctionRequest.User = user
	pd.auctionRequest.OpenRTBVersion = openrtb.OpenRTB2_5

	return &pd.auctionRequest, nil
}

// parseDeviceIDExtended extracts and validates device IFA from byte slice.
func parseDeviceIDExtended(ifa []byte) ([]byte, error) {
	if len(ifa) == 0 {
		return nil, errors.New("empty device IFA")
	}

	// Remove dashes from UUID format if present
	// e.g., "9c1ce5bd-3013-5c90-2598-cd744e1c96d6" -> "9c1ce5bd30135c902598cd744e1c96d6"
	result := make([]byte, 0, 32)
	for _, b := range ifa {
		if b != '-' {
			result = append(result, b)
		}
	}

	if len(result) != 32 {
		return nil, errors.New("invalid IFA length, expected 32 hex chars: " + strconv.Itoa(len(result)))
	}

	return result, nil
}
