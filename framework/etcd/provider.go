package etcd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/client"
)

var (
	errEmptyKeysAPI = errors.New("empty client.KeysAPI")
)

// ServerProvider ....
type ServerProvider interface {
	Name() string
	Addr() string
	// Heartbeat loop to set key and vlaue with ttl...
	Heartbeat(context.Context, client.KeysAPI, *ProvideOptions)
	Quit(client.KeysAPI, *ProvideOptions) error
}

type defaultServerProvider struct {
	name string
	addr string
}

// Name ...
func (d defaultServerProvider) Name() string {
	return d.name
}

// Addr ...
func (d defaultServerProvider) Addr() string {
	return d.addr
}

// ref to: http://ralphbupt.github.io/2017/05/04/etcd-%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0/
func (d defaultServerProvider) Heartbeat(
	ctx context.Context, kapi client.KeysAPI, opt *ProvideOptions,
) {
	for {
		select {
		case <-ctx.Done():
			if isDebug {
				fmt.Println("context done")
			}
			return
		default:
			if err := Provide(kapi, d, opt); err != nil {
				fmt.Println(err)
			}
		}
		if isDebug {
			fmt.Println("provider heartbeat doing")
		}
		time.Sleep(10 * time.Second)
	}
}

func (d defaultServerProvider) Quit(kapi client.KeysAPI, opt *ProvideOptions) error {
	return Delete(kapi, d, opt)
}

// ProvideOptions ...
type ProvideOptions struct {
	NamePrefix string             // ServerName Prefix if "" means no prefix
	SetOpts    *client.SetOptions // etcd client SetOptions
}

// Provide ...
func Provide(kapi client.KeysAPI, provider ServerProvider, opt *ProvideOptions) error {
	if kapi == nil {
		return errEmptyKeysAPI
	}

	// do put the info into with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// should use Create or Set
	// https://godoc.org/go.etcd.io/etcd/client#KeysAPI
	var (
		key     = strings.TrimPrefix(provider.Name(), "/")
		value   = provider.Addr()
		setOpts = &client.SetOptions{TTL: time.Second * 12}
	)

	if opt != nil {
		if len(opt.NamePrefix) != 0 {
			key = fmt.Sprintf("/%s/%s", opt.NamePrefix, key)
		}
		if opt.SetOpts != nil {
			setOpts = opt.SetOpts
		}
	}

	_, err := kapi.Set(ctx, key, value, setOpts)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Delete ...
func Delete(kapi client.KeysAPI, provider ServerProvider, opt *ProvideOptions) error {
	if isDebug {
		fmt.Println("called delete func by provider")
	}

	if kapi == nil {
		return errEmptyKeysAPI
	}

	var (
		key = strings.TrimPrefix(provider.Name(), "/")
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if opt != nil {
		if len(opt.NamePrefix) != 0 {
			key = fmt.Sprintf("/%s/%s", opt.NamePrefix, key)
		}
	}

	if _, err := kapi.Delete(ctx, key, nil); err != nil {
		return err
	}

	return nil
}

// NewProvider ...
func NewProvider(name, addr string) ServerProvider {
	return defaultServerProvider{name, addr}
}
