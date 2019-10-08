package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	amqpw "github.com/yeqown/infrastructure/framework/amqp"
)

const (
	ex         = "t-ex"
	queue      = "t-queue"
	routingKey = "t-routing"
)

func apply(ch *amqp.Channel) (err error) {
	if err = ch.ExchangeDeclare(ex, "direct", true, true, false, true, nil); err != nil {
		return err
	}

	if _, err = ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	if err = ch.QueueBind(queue, routingKey, ex, false, nil); err != nil {
		return err
	}
	return nil
}

func consume(d amqp.Delivery) error {
	fmt.Printf("data: %s\n", d.Body)
	d.Ack(true)
	return nil
}

func main() {
	w := amqpw.New(
		"amqp://username:password@host:port",
		amqp.Config{
			Vhost:     "vhost",
			Heartbeat: 2 * time.Second,
		}, apply)

	go w.Consume(queue, "", false, false, false, false, nil, consume)

	quit := make(chan bool)

	<-quit
}
