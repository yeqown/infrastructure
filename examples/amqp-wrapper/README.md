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

	applyTopology ApplyTopology
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan bool
	changeConn    chan struct{}
	chNotify      chan *amqp.Error // channel notify
	connNotify    chan *amqp.Error // conn notify

	isConnected bool // mark wrapper is connected to server
	hasConsumer bool // mark wrapper is used by a consumer
}
```

## Consumer

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

```go

func main() {
	w := amqpw.New(
		"amqp://username:password@host:port",
		amqp.Config{
			Vhost:     "vhost",
			Heartbeat: 2 * time.Second,
		}, apply)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-ticker.C:
				if err := w.Produce(ex, routingKey, false, false, []byte("hello")); err != nil {
					log.Println("could not produce: ", err)
				}
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	quit := make(chan bool)

	<-quit
}
```