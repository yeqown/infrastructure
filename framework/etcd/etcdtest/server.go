package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yeqown/server-common/framework/etcd"
)

var (
	port = flag.Int("port", 3456, "setting port to listen")
	sig  = make(chan os.Signal)
)

type response struct {
	Message string `json:"message"`
}

func hdlHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling this handler")
	resp := response{
		Message: fmt.Sprintf(
			"Hello! this is server listen on: %d", *port),
	}
	bs, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(bs))
}

func main() {
	flag.Parse()
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP)

	http.HandleFunc("/hello", hdlHello)

	endpoints := []string{
		"http://127.0.0.1:2377",
		"http://127.0.0.1:2379",
		"http://127.0.0.1:2378",
	}
	etcd.OpenDebug(true)
	kapi, err := etcd.Connect(endpoints...)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(2)
	}

	provider := etcd.NewProvider(
		fmt.Sprintf("srv_%d", *port),              // name
		fmt.Sprintf("http://127.0.0.1:%d", *port), // addr
	)

	ctx, cancel := context.WithCancel(context.Background())
	go provider.Heartbeat(ctx, kapi, &etcd.ProvideOptions{
		NamePrefix: "prefix",
		SetOpts:    nil,
	})

	fmt.Println("server listen on: ", *port)
	go http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	select {
	case <-sig:
		fmt.Println("cancel called")
		cancel()
		provider.Quit(kapi, &etcd.ProvideOptions{
			NamePrefix: "prefix",
			SetOpts:    nil,
		})
	}
}
