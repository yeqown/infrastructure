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

	for {
		select {
		case <-w.changeConn:
			log.Println("evt 'changeConn' triggered.")
			if ch, err = w.Channel(10 * time.Second); err != nil {
				panic(err)
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
				panic(err)
			}
			log.Println("init comsumer done")
		default:
			if !w.isConnected || delivery == nil {
				time.Sleep(1 * time.Second)
				break
			}
			// 如果异常失去链接，delivery 会被关闭，默认情况下会一直从 delivery 通道中获取数据
			for d := range delivery {
				if err := handler(d); err != nil {
					log.Printf("consume msg error: %v", err)
				}
			}
		}
	}
}
