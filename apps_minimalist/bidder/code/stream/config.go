package stream

import "time"

type Config struct {
	Disable      bool          `envconfig:"KAFKA_DISABLE" default:"true"`
	Brokers      []string      `envconfig:"KAFKA_BROKERS"`
	Topic        string        `envconfig:"KAFKA_TOPIC" default:"bidder-events"`
	BatchSize    int           `envconfig:"KAFKA_BATCH_SIZE" default:"1000"`
	BatchTimeout time.Duration `envconfig:"KAFKA_BATCH_TIMEOUT" default:"500ms"`
	BufferSize   int           `envconfig:"KAFKA_BUFFER_SIZE" default:"100000"`
}
