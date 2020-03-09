package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

// SSHConfig .
type SSHConfig struct {
	Host           string `json:"host"`
	User           string `json:"user"`
	Port           int    `json:"port"`
	Secret         string `json:"secret"`
	PrivateKeyFile string `json:"privateKeyFile"`
}

// TunnelConfig .
type TunnelConfig struct {
	Ident      string     `json:"ident"`
	SSH        *SSHConfig `json:"ssh"`
	LocalPort  int        `json:"localPort"`
	RemoteHost string     `json:"remoteHost"`
	RemotePort int        `json:"remotePort"`
}

var (
	// this will never be released manually,
	// but it would be cleared when the program finished
	identUnique = make(map[string]bool, 4)
)

// Valid .
func (c *TunnelConfig) Valid() error {
	if c.Ident == "" {
		return errors.New("empty tunnel identify")
	}

	if _, ok := identUnique[c.Ident]; ok {
		return errors.New("duplicate tunnel ident=" + c.Ident)
	}
	identUnique[c.Ident] = true

	// if c.SSH == nil {
	// 	return errors.New("empty ssh config")
	// }

	if c.LocalPort == 0 {
		// pass
	}

	if c.RemoteHost == "" {
		return errors.New("empty remote host")
	}

	if c.RemotePort == 0 {
		return errors.New("empty remote port")
	}

	return nil
}

// SSHTunnel .
type SSHTunnel struct {
	LocalAddr  string            // format: "host:port"
	ServerAddr string            // format: "host:port"
	RemoteAddr string            // format: "host:port"
	SSHConfig  *ssh.ClientConfig // ssh client config
}

// output: tunnel=(localhost:6379)
func (tunnel *SSHTunnel) name() string {
	return "tunnel=(" + tunnel.LocalAddr + ")"
}

// NewSSHTunnel .
func NewSSHTunnel(tunnelConfig *TunnelConfig) *SSHTunnel {
	var (
		auth      ssh.AuthMethod
		sshConfig = tunnelConfig.SSH
	)

	if sshConfig.PrivateKeyFile != "" {
		// true: privateKey specified
		auth = loadPrivateKeyFile(sshConfig.PrivateKeyFile)
	} else {
		auth = ssh.Password(sshConfig.Secret)
	}

	return &SSHTunnel{
		SSHConfig: &ssh.ClientConfig{
			User: sshConfig.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// logger.Infof("accept hostkey callback")
				//  Always accept key.
				return nil
			},
		},
		LocalAddr:  assembleAddr("localhost", tunnelConfig.LocalPort),
		ServerAddr: assembleAddr(tunnelConfig.SSH.Host, tunnelConfig.SSH.Port),
		RemoteAddr: assembleAddr(tunnelConfig.RemoteHost, tunnelConfig.RemotePort),
	}
}

// format liek "host:port"
func assembleAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Start .
// TODO: support random port by using localhost:0
func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.LocalAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	// tunnel.Local.Port = listener.Addr().(*net.TCPAddr).Port
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		logger.Infof(tunnel.name() + " accepted connection")
		go tunnel.forward(conn)
	}
}

// just do the work like proxy to transfer data from local to remote
func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverSSHClient, err := ssh.Dial("tcp", tunnel.ServerAddr, tunnel.SSHConfig)
	if err != nil {
		logger.Infof(tunnel.name()+" server dial error: %s", err)
		return
	}
	logger.Infof(tunnel.name()+" connected to server=%s (1 of 2)", tunnel.ServerAddr)
	remoteConn, err := serverSSHClient.Dial("tcp", tunnel.RemoteAddr)
	if err != nil {
		logger.Infof(tunnel.name()+" remote dial error: %s", err)
		return
	}
	logger.Infof(tunnel.name()+" connected to remote=%s (2 of 2)", tunnel.RemoteAddr)

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			logger.Infof(tunnel.name()+" io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// loadPrivateKeyFile . load privare file by @dir
func loadPrivateKeyFile(dir string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
