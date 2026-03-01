# OpenRTB 2.5 Bid Request Fixtures

This directory contains example OpenRTB 2.5 bid request fixtures based on the official IAB OpenRTB 2.5 specification and industry examples.

## Fixtures

### 1. Banner (Display) - Website
**File:** `banner-site.json`

Standard display banner ad request for a website. Features:
- 300x250 medium rectangle banner
- Site context with publisher and content metadata
- Desktop device (Macbook)
- User demographic data
- Brand safety controls (blocked categories and advertisers)

**Use case:** Traditional web display advertising

### 2. Video - Website
**File:** `video-site.json`

In-stream video ad request for pre-roll placement. Features:
- 640x480 video player
- VAST 2.0/3.0 protocol support
- 5-30 second duration constraint
- Companion banner ads (300x250, 728x90)
- Video content metadata (car show episode)
- User interest segments (auto intenders/enthusiasts)

**Use case:** Pre-roll video advertising on content sites

### 3. Audio - Podcast
**File:** `audio-site.json`

Digital audio ad request for podcast content. Features:
- DAAST-compliant audio protocols
- 15-30 second audio spot
- Bitrate range: 128-320 kbps
- Companion display ads
- Podcast episode metadata (Tech Talk episode 42)
- Mobile device context (iPhone)
- Volume normalization settings

**Use case:** Podcast and streaming audio advertising

### 4. Banner (Display) - Mobile App
**File:** `banner-app.json`

Mobile app banner ad request. Features:
- 320x50 mobile banner (with alternative 300x250)
- App object (vs site) with bundle ID and store URL
- Android mobile device (Google Pixel 5)
- Device advertising identifier (IFA/IDFA)
- App category: Gaming (IAB9-30)
- Multiple creative format support

**Use case:** In-app mobile advertising

## OpenRTB 2.5 Key Concepts

### Impression Types
OpenRTB 2.5 supports four main impression types via subordinate objects of the `Imp` object:
- **Banner** - Display ads (images, expandable units, in-banner video)
- **Video** - In-stream video (typically VAST)
- **Audio** - Audio ads (typically DAAST)
- **Native** - Native ads (not included in these fixtures)

### Site vs App
Each bid request must include **either** a `site` or `app` object (never both):
- **Site** - Browser-based web content
- **App** - Non-browser mobile/tablet applications

### Required Fields

According to OpenRTB 2.5 spec, minimal required fields are:
- `id` - Unique auction ID
- `imp` - Array of at least one impression object
  - `id` - Unique impression ID within the request
  - One of: `banner`, `video`, `audio`, or `native`

### Recommended Fields

Highly recommended for better bid optimization:
- `site` or `app` - Publisher context
- `device` - Device and geo information
- `user` - User demographics and segments
- `bcat` - Blocked advertiser categories (brand safety)
- `badv` - Blocked advertiser domains (competitive separation)

## Field Explanations

### Common Fields

- **`at`** - Auction type (1 = first price, 2 = second price)
- **`tmax`** - Max timeout in milliseconds (default: 120ms)
- **`cur`** - Currency (ISO-4217 codes)
- **`bcat`** - Blocked categories (IAB taxonomy)
- **`badv`** - Blocked advertiser domains

### Impression Object

- **`bidfloor`** - Minimum bid floor (CPM)
- **`bidfloorcur`** - Currency for bid floor
- **`secure`** - HTTPS required (1 = yes)
- **`pmp`** - Private marketplace deals

### Banner Object

- **`w`, `h`** - Width and height in pixels
- **`pos`** - Ad position (0 = unknown, 1 = above fold, etc.)
- **`battr`** - Blocked creative attributes
- **`topframe`** - Is ad in top frame? (1 = yes)
- **`expdir`** - Expandable directions (1 = left, 2 = right, 3 = up, 4 = down)
- **`api`** - Supported API frameworks (3 = MRAID-1, 5 = MRAID-2)

### Video Object

- **`mimes`** - Supported MIME types (e.g., "video/mp4")
- **`minduration`, `maxduration`** - Min/max duration in seconds
- **`protocols`** - Supported video protocols (2 = VAST 2.0, 3 = VAST 3.0)
- **`w`, `h`** - Player width/height
- **`startdelay`** - Pre-roll (0), mid-roll (>0), post-roll (-1)
- **`linearity`** - Linear (1) or non-linear (2)
- **`playbackmethod`** - Auto-play sound on (1), auto-play muted (2), click-to-play (3)
- **`companionad`** - Array of companion banner ads

### Audio Object

- **`mimes`** - Supported MIME types (e.g., "audio/mp4")
- **`minduration`, `maxduration`** - Min/max duration in seconds
- **`protocols`** - Supported audio protocols (9 = DAAST 1.0, 10 = DAAST 1.0 wrapper)
- **`minbitrate`, `maxbitrate`** - Bitrate range in kbps
- **`feed`** - Type of audio feed (1 = music, 2 = broadcast, 3 = podcast)
- **`nvol`** - Volume normalization mode
- **`companionad`** - Array of companion banner ads

### Device Object

- **`ua`** - User agent string
- **`geo`** - Geographic location data
- **`dnt`** - Do not track flag (0 = tracking OK, 1 = do not track)
- **`lmt`** - Limit ad tracking (iOS)
- **`ip`** - IPv4 address
- **`devicetype`** - Type (1 = mobile, 2 = PC, 4 = phone, 5 = tablet)
- **`ifa`** - Device advertising identifier (IDFA/AAID)

### User Object

- **`id`** - Exchange-specific user ID
- **`buyeruid`** - Buyer-specific user ID (from cookie sync)
- **`yob`** - Year of birth
- **`gender`** - Gender (M/F/O)
- **`data`** - Additional user segments from data providers

## IAB Content Taxonomy Examples

Sample category codes used in these fixtures:

- **IAB1** - Arts & Entertainment
- **IAB2** - Automotive
- **IAB3** - Business
- **IAB7-39** - Weapons (blocked category)
- **IAB8-5** - Drugs/Alcohol (blocked category)
- **IAB8-18** - Gambling (blocked category)
- **IAB9-30** - Video & Computer Games
- **IAB19-6** - Podcasts
- **IAB25-3** - Politics (blocked category)
- **IAB26** - Illegal Content (blocked category)

Full taxonomy: https://iabtechlab.com/standards/content-taxonomy/

## Testing These Fixtures

You can use these fixtures to test your RTB bidder:

```bash
# Test banner request
curl -X POST http://localhost:8080/bidrequest \
  -H "Content-Type: application/json" \
  -H "x-openrtb-version: 2.5" \
  -d @fixtures/banner-site.json

# Test video request
curl -X POST http://localhost:8080/bidrequest \
  -H "Content-Type: application/json" \
  -H "x-openrtb-version: 2.5" \
  -d @fixtures/video-site.json

# Test audio request
curl -X POST http://localhost:8080/bidrequest \
  -H "Content-Type: application/json" \
  -H "x-openrtb-version: 2.5" \
  -d @fixtures/audio-site.json

# Test mobile app request
curl -X POST http://localhost:8080/bidrequest \
  -H "Content-Type: application/json" \
  -H "x-openrtb-version: 2.5" \
  -d @fixtures/banner-app.json
```

## Sources

These fixtures are based on:

1. **IAB OpenRTB 2.5 Specification** - [Official PDF](https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-API-Specification-Version-2-5-FINAL.pdf)
2. **OpenRTB GitHub Examples** - [Video Request Example](https://github.com/openrtb/examples/blob/master/spotxchange/example-video-request-single_impr.md)
3. **Bidswitch Protocol Documentation** - [Audio Object Example](https://protocol.bidswitch.com/standards-v57/audio.html)
4. **Google OpenRTB Integration Guide** - [OpenRTB Guide](https://developers.google.com/authorized-buyers/rtb/openrtb-guide)
5. **Smaato OpenRTB 2.5 Specifications** - [Developer Docs](https://developers.smaato.com/marketers/openrtb-2-5-specifications/)

## Validation

To validate these fixtures against the OpenRTB 2.5 schema, you can use tools like:

- **OpenRTB Validator** - https://github.com/InteractiveAdvertisingBureau/openrtb
- **JSON Schema validation** using the official OpenRTB JSON schemas

## Notes

- All device IPs are examples (192.168.x.x)
- IFA/device IDs are randomly generated UUIDs
- User IDs are anonymized/hashed examples
- Currency is USD for all examples
- All examples use first-price auction (`at: 1`) except video which uses second-price (`at: 2`)
- GDPR is set to 0 (not applicable) for simplicity - adjust for EU traffic

## Further Reading

- **OpenRTB 2.5 Full Spec**: https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-API-Specification-Version-2-5-FINAL.pdf
- **OpenRTB GitHub**: https://github.com/InteractiveAdvertisingBureau/openrtb2.x
- **IAB Tech Lab**: https://iabtechlab.com/standards/openrtb/
- **VAST Specification**: https://iabtechlab.com/standards/vast/
- **DAAST Specification**: https://iabtechlab.com/standards/daast/
