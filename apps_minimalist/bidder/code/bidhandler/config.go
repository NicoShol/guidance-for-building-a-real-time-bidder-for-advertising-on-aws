package bidhandler

// TODO: add some type like for instance openrtb.Version
type Config struct {
	OpenRTBVersion float32 `envconfig:"BIDREQUEST_OPEN_RTB_VERSION" default:"3.0"`
}