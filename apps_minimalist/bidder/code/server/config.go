package server

type Config struct {
	Address string `envconfig:"SERVER_ADDRESS" required:"true"`
	BidRequestPath  string `envconfig:"SERVER_BIDREQUEST_PATH" required:"true"`
}