package stream

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer    *kafka.Writer
	inputChan chan []byte
	wg        sync.WaitGroup
}

func newProducer(cfg Config) *producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		//Async:        false,  // DEBUG, set to true for better performance but more complex error handling (so debugging --> false for now)
		Async:        true,  // DEBUG, set to true for better performance but more complex error handling (so debugging --> false for now)
		Compression:  kafka.Snappy,
		ErrorLogger:  kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Error().Msgf(msg, args...)
		}),
	}

	p := &producer{
		writer:    w,
		inputChan: make(chan []byte, cfg.BufferSize),
	}

	p.wg.Add(1)
	go p.drain()

	return p
}

// put enqueues a message. Returns immediately as long as the buffer isn't full.
// The data is copied because fasthttp reuses request body buffers across requests.
func (p *producer) put(data []byte) {
	copied := make([]byte, len(data))
	copy(copied, data)

	select {
	case p.inputChan <- copied:
	default:
		log.Warn().Msg("kafka producer buffer full, dropping message")
	}
}

// drain reads from inputChan and writes to Kafka.
// kafka.Writer handles batching internally (BatchSize / BatchTimeout).
func (p *producer) drain() {
	defer p.wg.Done()
	for data := range p.inputChan {
		if err := p.writer.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
			log.Error().Err(err).Msg("kafka write failed")
		}
	}
}

// close flushes remaining messages and shuts down the writer.
func (p *producer) close() error {
	close(p.inputChan)
	p.wg.Wait()
	return p.writer.Close()
}
