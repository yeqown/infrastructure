package amqp

import (
	"time"

	"github.com/streadway/amqp"
)

// Produce .
func (w *Wrapper) Produce(exchange, key string, mandatory, immediate bool, dat []byte) error {
	ch, err := w.Channel(5 * time.Second)
	if err != nil {
		return err
	}

	return ch.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "text/plain",
		Body:        dat,
	})
}
