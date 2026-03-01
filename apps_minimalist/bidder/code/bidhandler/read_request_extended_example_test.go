package bidhandler

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"apps_minimalist/bidder/code/auction"
)

// TestExtendedParserWithFixtures demonstrates the extended parser using
// the OpenRTB 2.5 fixtures we created.
func TestExtendedParserWithFixtures(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
	}{
		{
			name:        "Banner Site Request",
			fixturePath: "../../../fixtures/banner-site.json",
		},
		{
			name:        "Video Site Request",
			fixturePath: "../../../fixtures/video-site.json",
		},
		{
			name:        "Audio Site Request",
			fixturePath: "../../../fixtures/audio-site.json",
		},
		{
			name:        "Banner App Request",
			fixturePath: "../../../fixtures/banner-app.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read fixture file
			data, err := os.ReadFile(tt.fixturePath)
			if err != nil {
				t.Fatalf("Failed to read fixture: %v", err)
			}

			// Initialize persistent data with parser
			pd := newPersistenData()

			// Parse the request
			req, err := parseBidRequestExtended_2_5(data, pd)
			if err != nil {
				t.Fatalf("Failed to parse request: %v", err)
			}

			// Print parsed data to show what we extracted
			fmt.Printf("\n=== %s ===\n", tt.name)
			printExtendedRequest(req)
		})
	}
}

// printExtendedRequest pretty-prints the extracted data to show
// the difference between minimal and extended parsing.
func printExtendedRequest(req *auction.ExtendedRequest) {
	fmt.Printf("Request ID: %s\n", req.ID)
	fmt.Printf("OpenRTB Version: %v\n", req.OpenRTBVersion)

	// Print Device Information
	fmt.Println("\n--- DEVICE ---")
	fmt.Printf("  IFA: %s\n", req.Device.IFA)
	fmt.Printf("  UA: %s\n", req.Device.UA)
	fmt.Printf("  IP: %s\n", req.Device.IP)
	fmt.Printf("  Device Type: %d\n", req.Device.DeviceType)
	fmt.Printf("  Make: %s\n", req.Device.Make)
	fmt.Printf("  Model: %s\n", req.Device.Model)
	fmt.Printf("  OS: %s v%s\n", req.Device.OS, req.Device.OSV)
	fmt.Printf("  Screen: %dx%d (PPI: %d, Ratio: %.2f)\n",
		req.Device.W, req.Device.H, req.Device.PPI, req.Device.PxRatio)
	fmt.Printf("  Language: %s\n", req.Device.Language)
	fmt.Printf("  Carrier: %s (MCC-MNC: %s)\n", req.Device.Carrier, req.Device.MCCMNC)
	fmt.Printf("  Connection: %d, DNT: %d, Lmt: %d\n",
		req.Device.ConnectionType, req.Device.DNT, req.Device.Lmt)

	// Print Device Geo
	if req.Device.Geo.Lat != 0 || req.Device.Geo.Lon != 0 {
		fmt.Println("\n  Device Geo:")
		fmt.Printf("    Location: %.4f, %.4f (Type: %d)\n",
			req.Device.Geo.Lat, req.Device.Geo.Lon, req.Device.Geo.Type)
		fmt.Printf("    City: %s, Region: %s, Country: %s\n",
			req.Device.Geo.City, req.Device.Geo.Region, req.Device.Geo.Country)
		fmt.Printf("    Zip: %s, Metro: %s\n",
			req.Device.Geo.Zip, req.Device.Geo.Metro)
	}

	// Print User Information
	fmt.Println("\n--- USER ---")
	fmt.Printf("  ID: %s\n", req.User.ID)
	fmt.Printf("  Buyer UID: %s\n", req.User.BuyerUID)
	fmt.Printf("  YOB: %d\n", req.User.YOB)
	fmt.Printf("  Gender: %s\n", req.User.Gender)
	fmt.Printf("  Keywords: %s\n", req.User.Keywords)

	// Print User Geo (if different from device geo)
	if req.User.Geo.Lat != 0 || req.User.Geo.Lon != 0 {
		fmt.Println("\n  User Home Geo:")
		fmt.Printf("    Location: %.4f, %.4f\n", req.User.Geo.Lat, req.User.Geo.Lon)
		fmt.Printf("    City: %s, Region: %s\n", req.User.Geo.City, req.User.Geo.Region)
	}

	// Print User Data Segments
	if len(req.User.Data) > 0 {
		fmt.Println("\n  User Audience Segments:")
		for _, data := range req.User.Data {
			fmt.Printf("    Provider: %s (ID: %s)\n", data.Name, data.ID)
			for _, seg := range data.Segments {
				fmt.Printf("      - %s (ID: %s, Value: %s)\n",
					seg.Name, seg.ID, seg.Value)
			}
		}
	}

	// Print Impression Items
	fmt.Printf("\n--- IMPRESSIONS (%d) ---\n", len(req.Items))
	for i, item := range req.Items {
		fmt.Printf("\n  [%d] ID: %s\n", i+1, item.ID)
		fmt.Printf("      Bid Floor: %.2f %s\n", item.BidFloor, item.BidFloorCur)
		fmt.Printf("      Tag ID: %s\n", item.TagID)
		fmt.Printf("      Secure: %d, Exp: %d\n", item.Secure, item.Exp)

		if item.BannerWidth > 0 {
			fmt.Printf("      Banner: %dx%d\n", item.BannerWidth, item.BannerHeight)
			if len(item.BannerMimes) > 0 {
				fmt.Printf("      Mimes: ")
				for j, mime := range item.BannerMimes {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", mime)
				}
				fmt.Println()
			}
		}
	}

	fmt.Println("\n" + string(make([]byte, 60, 60)[:60])) // separator
}

// BenchmarkMinimalVsExtendedParsing compares performance of minimal vs extended parsing.
func BenchmarkMinimalVsExtendedParsing(b *testing.B) {
	// Read a realistic fixture
	data, err := os.ReadFile("../../../fixtures/banner-app.json")
	if err != nil {
		b.Fatalf("Failed to read fixture: %v", err)
	}

	b.Run("Extended Parsing", func(b *testing.B) {
		pd := newPersistenData()

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := parseBidRequestExtended_2_5(data, pd)
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
		}
	})

	b.Run("Standard JSON Unmarshal", func(b *testing.B) {
		var v map[string]interface{}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(data, &v)
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
		}
	})
}
