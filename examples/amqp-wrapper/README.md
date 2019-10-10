# amqp-warapper

Wrapping `connection` and `channel` of `amqp` to handle reconnection of `Consumer` and `Producer`.

## Wrapper

```go
// ApplyTopology to apply resource from MQ server
// eg. QueueDeclare, ExchangeDeclare
type ApplyTopology func(ch *amqp.Channel) error

// Wrapper .
type Wrapper struct {
	Addr   string
	Config amqp.Config

	applyTopology ApplyTopology         // request Topology resource function
	connection    *amqp.Connection      // connectio to amqp.Server
	channel       *amqp.Channel         // channel
	done          chan bool
	changeConn    chan struct{}         // the signal of changing the connection and channel
	notifyClose   chan *amqp.Error      // channel be closed or any error
	// notifyConfirm chan amqp.Confirmation
	isConnected bool                    // mark the wrapper has connected to server
}
```

## Consumer

```go
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
			log.Println("delivery was cleared")
		}
	}
}
```

```go
w := amqpw.New(
    "amqp://user:password@host:port",
    amqp.Config{
        Vhost:     "vhost",
        Heartbeat: 2 * time.Second,
    }, apply)

go w.Consume(queue, "", false, false, false, false, nil, consume)
```

## Producer
