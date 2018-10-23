package etcd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
)

// Server ...
type Server struct {
	Name  string `json:"name"`
	Addr  string `json:"addr"`
	Alive bool   `json:"alive"`
}

// Watcher ...
// ref to http://daizuozhuo.github.io/etcd-service-discovery/
type Watcher struct {
	members        map[string]Server // members dict container
	kapi           client.KeysAPI    // KeysAPI
	prefix         string            // srvName prefix
	watchDruration time.Duration     // loop duration
	sync.RWMutex                     // RW mutext
}

// RangeMember ...
func (w *Watcher) RangeMember() []Server {
	srvs := make([]Server, 0)
	w.RLock()
	for _, srv := range w.members {
		if !srv.Alive {
			continue
		}
		srvs = append(srvs, srv)
	}
	w.RUnlock()
	return srvs
}

// AddMember ...
func (w *Watcher) AddMember(key, value string) {
	key = strings.TrimPrefix(key, w.prefix)

	w.Lock()
	if m, ok := w.members[key]; ok {
		m.Alive = true
	} else {
		w.members[key] = Server{
			Name:  key,
			Addr:  value,
			Alive: true,
		}
	}
	w.Unlock()
}

// DeleteMember ...
func (w *Watcher) DeleteMember(key string) {
	key = strings.TrimPrefix(key, w.prefix)

	w.Lock()
	delete(w.members, key)
	w.Unlock()
}

// ExpireMember ...
func (w *Watcher) ExpireMember(key string) {
	w.Lock()
	if m, ok := w.members[key]; ok {
		m.Alive = false
	}
	w.Unlock()
}

// Watch ...
func (w *Watcher) Watch() {
	etcdWatcher := w.kapi.Watcher(w.prefix, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		resp, err := etcdWatcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch: ", err)
		}

		if isDebug {
			fmt.Println(resp.Action, resp.Node.String())
		}

		switch resp.Action {
		case "expire":
			w.ExpireMember(resp.PrevNode.Key)
		case "set", "update":
			w.AddMember(resp.Node.Key, resp.Node.Value)
		case "delete":
			w.DeleteMember(resp.Node.Key)
		default:
			log.Println("Error no such action:", resp.Action)
		}
		time.Sleep(w.watchDruration)
	}
}

// NewWatcher ...
func NewWatcher(
	kapi client.KeysAPI, serverPrefix string, d time.Duration,
) *Watcher {
	return &Watcher{
		members:        make(map[string]Server),
		kapi:           kapi,
		prefix:         serverPrefix,
		watchDruration: d,
		RWMutex:        sync.RWMutex{},
	}
}
