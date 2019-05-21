package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/yeqown/infrastructure/framework/etcd"
)

var (
	watcher *etcd.Watcher
	handler selfHandler
)

func main() {
	endpoints := []string{
		"http://127.0.0.1:2377",
		"http://127.0.0.1:2379",
		"http://127.0.0.1:2378",
	}
	kapi, err := etcd.Connect(endpoints...)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// debug more, more log ~
	etcd.OpenDebug(true)

	watcher = etcd.NewWatcher(kapi, "prefix", time.Second)
	go watcher.Watch()

	if err := http.ListenAndServe(":9090", handler); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

type selfHandler struct{}

func (s selfHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if strings.HasPrefix(path, "/gw") {
		srvs := watcher.RangeMember()

		if len(srvs) == 0 {
			fmt.Fprintf(w, "Error-500: no server alive")
			return
		}

		// rand call
		idx := rand.Intn(len(srvs))

		target, _ := url.Parse(srvs[idx].Addr)
		proxy := httputil.NewSingleHostReverseProxy(target)
		req.URL.Path = strings.TrimPrefix(path, "/gw")
		fmt.Println(srvs, len(srvs), idx, req.URL.Path)

		proxy.ServeHTTP(w, req)
		return
	}

	fmt.Fprintf(w, "Error-404: no such path")
	return
}

type response struct {
	Message string `json:"message"`
}
