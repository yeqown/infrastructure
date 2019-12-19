package amqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// ConsumeHandler .
type ConsumeHandler func(d amqp.Delivery) error

// Consume .
func (w *Wrapper) Consume(
	queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table,
	handler ConsumeHandler,
) {
	var (
		ch       *amqp.Channel
		delivery <-chan amqp.Delivery
		err      error
	)

	w.hasConsumer = true

	for {
		select {
		case <-w.changeConn:
			log.Println("evt 'changeConn' triggered.")
			if ch, err = w.Channel(5 * time.Second); err != nil {
				log.Println("could not get channel for now with error: ", err)
				break
			}
			if delivery, err = ch.Consume(
				queue,     // queue
				consumer,  // consumer
				autoAck,   // auto-ack
				exclusive, // exclusive
				noLocal,   // no-local
				noWait,    // no-wait
				args,      // args
			); err != nil {
				log.Println("could not start consuming with error: ", err)
				break
			}
			log.Println("initial comsumer finished")
		default:
			if !w.isConnected || delivery == nil {
				// true: wrapper has not connected or consumer has not initialized
				// must to wait `changeConn` evt
				time.Sleep(1 * time.Second)
				break
			}
			// delivery will be closed, then this `range` will be finished
			for d := range delivery {
				if err := handler(d); err != nil {
					log.Printf("could not consume message: %v with error: %v", d, err)
				}
			}
		}
	}
}
