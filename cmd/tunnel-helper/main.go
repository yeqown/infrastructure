package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/yeqown/infrastructure/pkg/cfgutil"
	"github.com/yeqown/log"
)

var (
	logger = log.NewLogger()

	cfgDir = flag.String("c", ".tunnelrc", "specified config file to load, default is ./.tunnelrc in json format")
)

type config struct {
	SSH     *SSHConfig      `json:"ssh"`
	Tunnels []*TunnelConfig `json:"tunnels"`
}

func main() {
	flag.Parse()

	cfg, err := loadConfig(*cfgDir)
	if err != nil {
		panic(err)
	}

	// start tunnels with cfg
	startTunnels(cfg)
}

func loadConfig(dir string) (*config, error) {
	cfg := new(config)
	r, err := cfgutil.Open(dir)
	if err != nil {
		return nil, err
	}
	err = cfgutil.LoadJSON(r, cfg)
	return cfg, err
}

// support context cancel
func startTunnels(cfg *config) {
	var (
		errChan    = make(chan errTunnel, 1) // 异常channel
		wg         = sync.WaitGroup{}        // 同步组
		tunnelChan = make(chan int, 1)       // 运行中tunnels 计数
	)
	wg.Add(len(cfg.Tunnels))

	// create and ssh tunnel and goto work
	for idx, v := range cfg.Tunnels {
		if v.SSH == nil {
			v.SSH = cfg.SSH
		}

		go func(idx int, cfg *TunnelConfig, errChan chan<- errTunnel) {
			defer wg.Done()
			defer func() { tunnelChan <- -1 }()
			tunnelChan <- 1

			// valid tunnel config
			if err := cfg.Valid(); err != nil {
				errChan <- newErrTunnel(idx, "invalid config, err=%v", err)
				return
			}

			// open tunnel and prepare
			tunnel := NewSSHTunnel(cfg)
			if err := tunnel.Start(); err != nil {
				errChan <- newErrTunnel(idx, "tunnel broken, err=%v", err)
				return
			}
		}(idx, v, errChan)
	}

	// log errors
	go func() {
		for err := range errChan {
			logger.Errorf("tunnelIdx=%d: %s", err.Idx, err.Errmsg)
		}
	}()

	// log tunnel changes
	go func() {
		running := 0
		msg := ""
		for cntChange := range tunnelChan {
			// if runningTunnelsCnt changed to notify
			running += cntChange
			if cntChange >= 0 {
				// true: starting
				msg = fmt.Sprintf("%d tunnel starting, current: %d", cntChange, running)
			} else {
				// true: quit
				msg = fmt.Sprintf("%d tunnel break, current: %d", 0-cntChange, running)
			}

			logger.Infof(msg)
		}

	}()

	wg.Wait()
	close(errChan)
	close(tunnelChan)
	// wait for all error message outputing
	time.Sleep(100 * time.Millisecond)
	logger.Infof("tunnel-helper finished")
}

// errTunnel .
type errTunnel struct {
	Idx    int
	Errmsg string
}

func newErrTunnel(idx int, format string, args ...interface{}) errTunnel {
	return errTunnel{
		Idx:    idx,
		Errmsg: fmt.Sprintf(format, args...),
	}
}

func (err errTunnel) Error() string {
	return err.Errmsg
}
