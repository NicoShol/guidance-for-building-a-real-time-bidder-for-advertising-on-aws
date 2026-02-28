package server

type Config struct {
	Address string `envconfig:"SERVER_ADDRESS" required:"true"`
}