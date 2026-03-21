package stream

import "github.com/rs/zerolog/log"

// Stream is the interface for streaming bid data to Kafka.
type Stream interface {
	PutRequest(data []byte)
	PutResponse(data []byte)
	Close() error
}

type kafkaStream struct {
	producer *producer
}

type noopStream struct{}

// New creates a Stream. Returns a no-op implementation if Kafka is disabled.
func New(cfg Config) Stream {
	if cfg.Disable {
		log.Info().Msg("kafka streaming disabled")
		return &noopStream{}
	}

	log.Info().
		Strs("brokers", cfg.Brokers).
		Str("topic", cfg.Topic).
		Int("buffer_size", cfg.BufferSize).
		Msg("starting kafka stream")

	return &kafkaStream{
		producer: newProducer(cfg),
	}
}

func (s *kafkaStream) PutRequest(data []byte)  {
	// log.Info().Msg("Putting request to stream")
	s.producer.put(data)
	// log.Info().Msg("Request put to stream")
}
func (s *kafkaStream) PutResponse(data []byte) { s.producer.put(data) }
func (s *kafkaStream) Close() error            { return s.producer.close() }

func (s *noopStream) PutRequest([]byte)  {}
func (s *noopStream) PutResponse([]byte) {}
func (s *noopStream) Close() error       { return nil }
