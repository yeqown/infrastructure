package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/yeqown/server-common/framework/etcd"
)

var (
	port = flag.Int("port", 3456, "setting port to listen")
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
	http.HandleFunc("/hello", hdlHello)

	fmt.Println("server listen on: ", *port)

	endpoints := []string{
		"http://127.0.0.1:2377",
		"http://127.0.0.1:2379",
		"http://127.0.0.1:2378",
	}
	kapi, err := etcd.Connect(endpoints...)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(2)
	}

	provider := etcd.NewProvider(
		fmt.Sprintf("srv_%d", *port),              // name
		fmt.Sprintf("http://127.0.0.1:%d", *port), // addr
	)
	go provider.Heartbeat(context.Background(), kapi, &etcd.ProvideOptions{
		NamePrefix: "prefix",
		SetOpts:    nil,
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Errorf(err.Error())
		os.Exit(2)
	}
}
