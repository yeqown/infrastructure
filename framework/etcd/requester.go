package etcd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
)

var (
	errNoneNamePrefix = errors.New("none length of request name prefix")
	isDebug           = false
)

// RequesterOptions ...
type RequesterOptions struct {
	NamePrefix string
	GetOpts    *client.GetOptions
}

// Server ...
type Server struct {
	Name  string `json:"name"`
	Addr  string `json:"addr"`
	Alive bool   `json:"alive"`
}

// RequestWithPrefix ...
func RequestWithPrefix(kapi client.KeysAPI, reqOpts *RequesterOptions) ([]Server, error) {
	if kapi == nil {
		return nil, errEmptyKeysAPI
	}

	var (
		// key     = strings.TrimPrefix(requester.Name(), "/")
		key     string
		getOpts *client.GetOptions
	)

	if reqOpts != nil {
		if len(reqOpts.NamePrefix) == 0 {
			return nil, errNoneNamePrefix
		}
		// make sure the format of key: "/prefix"
		key = fmt.Sprintf("/%s", strings.TrimPrefix(reqOpts.NamePrefix, "/"))

		if reqOpts.GetOpts != nil {
			getOpts = reqOpts.GetOpts
		}
	}

	response, err := kapi.Get(context.Background(), key, getOpts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	servers := make([]Server, 0)
	for _, node := range response.Node.Nodes {
		servers = append(servers, Server{
			Name:  strings.TrimPrefix(node.Key, key),
			Addr:  node.Value,
			Alive: true,
		})
	}

	return servers, nil
}

// Watcher ...
// ref to http://daizuozhuo.github.io/etcd-service-discovery/
type Watcher struct {
	members map[string]Server
	kapi    client.KeysAPI
	prefix  string
	sync.RWMutex
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
		response, err := etcdWatcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch: ", err)
		}

		if isDebug {
			fmt.Println(response.Action, response.Node.String())
		}

		switch response.Action {
		case "expire":
			w.ExpireMember(response.PrevNode.Key)
		case "set", "update":
			// TODO: set or updatte s.members with lock
			w.AddMember(response.Node.Key, response.Node.Value)
		case "delete":
			// TODO: delete s.members with lock
			w.DeleteMember(response.Node.Key)
		default:
			log.Println("Error no such action:", response.Action)
		}
		time.Sleep(5 * time.Second)
	}
}

// NewWatcher ...
func NewWatcher(kapi client.KeysAPI, serverPrefix string) *Watcher {
	return &Watcher{
		members: make(map[string]Server),
		kapi:    kapi,
		prefix:  serverPrefix,
		RWMutex: sync.RWMutex{},
	}
}
