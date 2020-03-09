package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/yeqown/infrastructure/pkg/cfgutil"
	"github.com/yeqown/log"
)

const (
	// defaultCfgName .
	defaultCfgName = ".tunnelrc"
)

var (
	logger = log.NewLogger()

	cfgDir             = flag.String("c", "", "specified config file to load, default is ./.tunnelrc in json format")
	initFlag           = flag.Bool("i", false, "create an default format config in current folder")
	tunnelIdentPattern = flag.String("p", `.*`, "pattern to match with specified tunnel to open") // FIXME: defaule pattern should be all
)

type config struct {
	SSH     *SSHConfig      `json:"ssh"`
	Tunnels []*TunnelConfig `json:"tunnels"`
}

func main() {
	flag.Parse()
	if shouldContinue := commandHelper(); !shouldContinue {
		return
	}

	cfg, err := loadConfig(*cfgDir)
	if err != nil {
		panic(err)
	}

	// start tunnels with cfg
	startTunnels(cfg)
}

// load config from file
// check duplicate config of tunnel
// support tunnel identify
func loadConfig(dir string) (*config, error) {
	var (
		cfg = new(config)
		err error
		r   io.ReadCloser
	)
	if r, err = cfgutil.Open(dir); err != nil {
		return nil, err
	}
	defer r.Close()

	if err = cfgutil.LoadJSON(r, cfg); err != nil {
		return nil, err
	}

	// valid tunnel config
	for _, v := range cfg.Tunnels {
		if err = v.Valid(); err != nil {
			logger.Errorf("invalid config, err=%v", err)
			return nil, err
		}
	}

	return cfg, err
}

func generateConfig() error {
	var defaultConfig = config{
		SSH: &SSHConfig{
			Host:           "172.168.1.1",
			User:           "username",
			Port:           22,
			Secret:         "SSH-PASSWORD",
			PrivateKeyFile: "/path/to/.ssh/id_rsa",
		},
		Tunnels: []*TunnelConfig{
			&TunnelConfig{
				Ident:      "tunnel-config-without-ssh",
				SSH:        nil,
				LocalPort:  8081,
				RemoteHost: "172.168.1.1",
				RemotePort: 8081,
			},
			&TunnelConfig{
				Ident: "tunnel-config-with-ssh",
				SSH: &SSHConfig{
					Host:           "172.168.1.2",
					User:           "username2",
					Port:           22,
					Secret:         "SSH-PASSWORD2",
					PrivateKeyFile: "/path/to/.ssh/id_rsa",
				},
				LocalPort:  8080,
				RemoteHost: "172.168.1.1",
				RemotePort: 8080,
			},
		},
	}
	fd, err := os.OpenFile(defaultCfgName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("could not open file=%s, err=%v", defaultCfgName, err)
		return err
	}

	defer fd.Close()

	byts, err := json.MarshalIndent(defaultConfig, "", "\t")
	if err != nil {
		logger.Errorf("could not marshal defaultConfig=%+v, err=%v", defaultConfig, err)
		return err
	}

	_, err = fd.Write(byts)
	return err
}

// do some commands here
// 1. handle init command
// 2. default load config from user home
func commandHelper() (shouldContinue bool) {
	shouldContinue = true

	if *initFlag {
		shouldContinue = false
		// generate config file in current folder
		if err := generateConfig(); err != nil {
			logger.Errorf("could not generate .tunnelrc config, err=%v", err)
		}
		return
	}

	// handle default config path
	if strings.Compare(*cfgDir, "") == 0 {
		// true: default config filepath
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Errorf("could not get user home dir, err=%v", err)
			return
		}
		*cfgDir = path.Join(home, defaultCfgName)
	}

	return
}

// support context cancel
// upport tunnel open by ident matcher
func startTunnels(cfg *config) {
	var (
		errChan    = make(chan errTunnel, 1) // 异常channel
		wg         = sync.WaitGroup{}        // 同步组
		tunnelChan = make(chan int, 1)       // 运行中tunnels 计数

		exp *regexp.Regexp // pattern regexp
	)

	// compile pattern
	exp = regexp.MustCompile(*tunnelIdentPattern)

	wg.Add(len(cfg.Tunnels))
	// create and ssh tunnel and goto work
	for idx, v := range cfg.Tunnels {
		if v.SSH == nil {
			v.SSH = cfg.SSH
		}

		// matche with ident pattern
		if !exp.MatchString(v.Ident) {
			logger.Warnf("tunnel ident=%s, not matched with pattern=%s, so skipped", v.Ident, *tunnelIdentPattern)
			continue
		}

		go func(idx int, tunnelCfg *TunnelConfig, errChan chan<- errTunnel) {
			defer wg.Done()
			defer func() { tunnelChan <- -1 }()
			tunnelChan <- 1

			// // valid tunnel config
			// // has been moved to loadConfig
			// if err := tunnelCfg.Valid(); err != nil {
			// 	errChan <- newErrTunnel(idx, "invalid config, err=%v", err)
			// 	return
			// }

			// open tunnel and prepare
			tunnel := NewSSHTunnel(tunnelCfg)
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

	// record changes of opening-tunnel count
	go func() {
		running := 0
		msg := ""
		for cntChange := range tunnelChan {
			// if runningTunnelsCnt changed to notify
			// FIXME: atomic op with running
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
	logger.Infof("tunnel-helper is finished")
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
