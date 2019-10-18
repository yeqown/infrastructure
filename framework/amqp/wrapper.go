package amqp

import (
	"errors"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	reconnectDelay     = 5 * time.Second // 连接断开后多久重连
	reconnectDetectDur = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to the AMQP server")
	errAlreadyClosed = errors.New("already closed: not connected to the AMQP server")
)

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

// handleReconnct
func (w *Wrapper) handleReconnect() {
	for {
		if !w.isConnected {
			log.Println("Attempting to connect")
			var (
				connected = false
				err       error
			)

			for cnt := 0; !connected; cnt++ {
				if connected, err = w.connect(); err != nil {
					log.Printf("Failed to connect: %s.\n", err)
				}
				if !connected {
					log.Printf("Retrying... %d\n", cnt)
				}
				time.Sleep(reconnectDelay)
			}
		}

		select {
		case <-w.done:
			println("evt `w.done` triggered")
			return
		case err := <-w.chNotify:
			log.Printf("channel close notify: %v", err)
			w.isConnected = false
		case err := <-w.connNotify:
			log.Printf("conn close notify: %v", err)
			w.isConnected = false
		}
		time.Sleep(reconnectDetectDur)
	}
}

// Connect .
func (w *Wrapper) connect() (bool, error) {
	conn, err := amqp.DialConfig(w.Addr, w.Config)
	if err != nil {
		return false, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return false, err
	}

	if err := w.applyTopology(ch); err != nil {
		return false, err
	}
	w.isConnected = true
	w.changeConnection(conn, ch)
	return true, nil
}

// 监听Rabbit channel的状态
func (w *Wrapper) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	w.connection = connection
	w.connNotify = make(chan *amqp.Error, 1)
	w.connection.NotifyClose(w.chNotify)

	w.channel = channel
	w.chNotify = make(chan *amqp.Error, 1)
	w.channel.NotifyClose(w.chNotify)

	// TOFIX: only producer will be blocked here
	if w.hasConsumer {
		// true: cause only consumer will be  notify for now.
		w.changeConn <- struct{}{}
	}
}

// Conn .
// func (w *Wrapper) Conn() *amqp.Connection {
// 	return w.connection
// }

// Channel . it will blocked
func (w *Wrapper) Channel(timeout time.Duration) (*amqp.Channel, error) {
	timer := time.NewTimer(timeout)
	for !w.isConnected {
		select {
		case <-timer.C:
			return nil, errNotConnected
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return w.channel, nil
}

// Close .
func (w *Wrapper) Close() error {
	if !w.isConnected {
		return errAlreadyClosed
	}
	err := w.channel.Close()
	if err != nil {
		return err
	}
	err = w.connection.Close()
	if err != nil {
		return err
	}
	close(w.done)
	w.isConnected = false
	return nil
}

// New .
// addr schema://username:pwd@host:port
func New(addr string, cfg amqp.Config, f ApplyTopology) *Wrapper {
	// ctx, cancel := context.WithCancel(context.Background())
	w := &Wrapper{
		Addr:          addr,
		applyTopology: f,
		Config:        cfg,
		changeConn:    make(chan struct{}, 1),
	}

	go w.handleReconnect()

	return w
}
