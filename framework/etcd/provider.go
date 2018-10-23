package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/client"
)

// ServerProvider ....
type ServerProvider interface {
	Name() string
	Addr() string
	KeysAPI() client.KeysAPI
	// Provide support set key-value to etcd
	Provide(*ProvideOptions) error
	// Quit while should be called while Server quit
	Quit(*ProvideOptions) error
}

type defaultServerProvider struct {
	name string
	addr string
	kapi client.KeysAPI
}

// Name ...
func (d defaultServerProvider) Name() string {
	return d.name
}

// Addr ...
func (d defaultServerProvider) Addr() string {
	return d.addr
}

// KeysAPI ...
func (d defaultServerProvider) KeysAPI() client.KeysAPI {
	return d.kapi
}

// Provide ...
func (d defaultServerProvider) Provide(opt *ProvideOptions) error {
	kapi := d.KeysAPI()
	if kapi == nil {
		return errEmptyKeysAPI
	}
	// do put the info into with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// should use Create or Set
	// https://godoc.org/go.etcd.io/etcd/client#KeysAPI
	var (
		err     error
		key     = strings.TrimPrefix(d.Name(), "/")
		value   = d.Addr()
		setOpts = &client.SetOptions{
			TTL: opt.TTLDuration,
		}
	)

	key, value, err = handleWithProvideOption(key, value, opt)
	if err != nil {
		return err
	}
	if opt.SetOpts != nil {
		setOpts = opt.SetOpts
	}

	if _, err = kapi.Set(ctx, key, value, setOpts); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Quit ...
func (d defaultServerProvider) Quit(opt *ProvideOptions) error {
	if isDebug {
		fmt.Println("called delete func by provider")
	}

	var (
		key  = strings.TrimPrefix(d.Name(), "/")
		kapi = d.KeysAPI()
	)

	if kapi == nil {
		return errEmptyKeysAPI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key, _, _ = handleWithProvideOption(key, "", opt)
	if _, err := kapi.Delete(ctx, key, nil); err != nil {
		return err
	}

	return nil
}

// ProvideOptions ...
type ProvideOptions struct {
	NamePrefix        string             // ServerName Prefix if "" means no prefix
	SetOpts           *client.SetOptions // etcd client SetOptions
	TTLDuration       time.Duration      // ttl time.Duration
	HeartbeatDuration time.Duration      // HeartbeatDuration duration
}

// ProviderHeartbeat ...
// ref to: http://ralphbupt.github.io/2017/05/04/etcd-%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0/
func ProviderHeartbeat(ctx context.Context,
	provider ServerProvider, opt *ProvideOptions,
) {
	for {
		select {
		case <-ctx.Done():
			if isDebug {
				fmt.Println("context done")
			}
			return
		default:
			if err := provider.Provide(opt); err != nil {
				fmt.Println(err)
			}
		}
		if isDebug {
			fmt.Println("provider heartbeat doing")
		}
		time.Sleep(opt.HeartbeatDuration)
	}
}

// NewProvider ...
func NewProvider(kapi client.KeysAPI, name, addr string) ServerProvider {
	return defaultServerProvider{
		name: name,
		addr: addr,
		kapi: kapi,
	}
}

// handleWithProvideOption ...
// deal key/value with opt *ProvideOptions
func handleWithProvideOption(
	key, value string,
	opt *ProvideOptions,
) (string, string, error) {
	if opt == nil {
		return key, value, errNilProvideOpt
	}
	// TODO: more opt handling
	key = fmt.Sprintf("/%s/%s", opt.NamePrefix, key)
	return key, value, nil
}
