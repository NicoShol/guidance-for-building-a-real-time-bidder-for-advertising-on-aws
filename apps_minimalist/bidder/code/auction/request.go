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

// ExtendedRequest contains comprehensive parsing of OpenRTB 2.5 bid request
// including full Device and User objects.
type ExtendedRequest struct {
	ID             []byte
	Items          []ExtendedItem
	Device         Device
	User           User
	OpenRTBVersion openrtb.Version
}

// ExtendedItem represents an impression object
type ExtendedItem struct {
	ID           []byte
	BidFloor     float64
	BidFloorCur  []byte
	Secure       int
	TagID        []byte
	Exp          int
	BannerWidth  int
	BannerHeight int
	BannerMimes  [][]byte
}

// Device represents the device object with all OpenRTB 2.5 fields
type Device struct {
	IFA            []byte  // Device advertising identifier (IDFA/AAID)
	UA             []byte  // User agent string
	IP             []byte  // IPv4 address
	IPv6           []byte  // IPv6 address
	DeviceType     int     // Device type (1=mobile, 2=PC, 4=phone, 5=tablet)
	Make           []byte  // Device make (e.g., "Apple")
	Model          []byte  // Device model (e.g., "iPhone")
	OS             []byte  // Operating system (e.g., "iOS")
	OSV            []byte  // OS version (e.g., "14.6")
	HWV            []byte  // Hardware version
	W              int     // Physical width in pixels
	H              int     // Physical height in pixels
	PPI            int     // Pixels per inch
	PxRatio        float64 // Pixel ratio
	JS             int     // JavaScript support (0=no, 1=yes)
	Language       []byte  // Browser language
	Carrier        []byte  // Carrier or ISP
	MCCMNC         []byte  // Mobile carrier code (MCC-MNC)
	ConnectionType int     // Connection type (0=unknown, 1=ethernet, 2=wifi, 3=cellular, etc.)
	DNT            int     // Do Not Track (0=tracking OK, 1=do not track)
	Lmt            int     // Limit Ad Tracking (0=tracking OK, 1=limit tracking)
	Geo            Geo     // Geographic location
}

// Geo represents geographic location data
type Geo struct {
	Lat     float64 // Latitude
	Lon     float64 // Longitude
	Type    int     // Location type (1=GPS, 2=IP, 3=user)
	Country []byte  // Country code (ISO-3166-1 Alpha-3)
	Region  []byte  // Region code (ISO-3166-2)
	Metro   []byte  // Google metro code
	City    []byte  // City name
	Zip     []byte  // Zip/postal code
	UTCOffset int   // UTC offset in minutes
}

// User represents the user object with all OpenRTB 2.5 fields
type User struct {
	ID        []byte        // Exchange-specific user ID
	BuyerUID  []byte        // Buyer-specific user ID
	YOB       int           // Year of birth
	Gender    []byte        // Gender ("M", "F", "O")
	Keywords  []byte        // Comma-separated keywords
	Geo       Geo           // User home geo (different from device geo)
	Data      []UserData    // Additional user data from providers
}

// UserData represents third-party audience segment data
type UserData struct {
	ID       []byte        // Data provider ID
	Name     []byte        // Data provider name
	Segments []UserSegment // User segments
}

// UserSegment represents a single audience segment
type UserSegment struct {
	ID    []byte // Segment ID
	Name  []byte // Segment name
	Value []byte // Segment value
}
