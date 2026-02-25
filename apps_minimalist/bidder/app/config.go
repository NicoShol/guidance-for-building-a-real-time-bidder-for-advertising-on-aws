package app

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" required:"true"`
}