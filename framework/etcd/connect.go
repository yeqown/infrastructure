// Package etcd includes etcd clients ops related
package etcd

import (
	"errors"
	"fmt"
	"time"

	"go.etcd.io/etcd/client"
)

var (
	isDebug           = false
	errEmptyKeysAPI   = errors.New("empty client.KeysAPI")
	errNoneNamePrefix = errors.New("none length of request name prefix")
	errNilProvideOpt  = errors.New("nil provide options")
)

// Connect to etcd client
// addr format like: http://host:port
func Connect(addrs ...string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: addrs,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 3 * time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client.NewKeysAPI(c), nil
}

// OpenDebug set debug mode on
func OpenDebug(open bool) {
	isDebug = open
}
