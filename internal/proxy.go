// MIT License
//
// Copyright (c) 2023 TTBT Enterprises LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	yaml "gopkg.in/yaml.v3"
)

// Proxy receives TLS connections and forwards them to the configured
// backends.
type Proxy struct {
	cm       *autocert.Manager
	cfg      *Config
	ctx      context.Context
	cancel   func()
	listener net.Listener

	mu      sync.Mutex
	mapping map[string]*Backend
	numOpen int
}

// New returns a new initialized Proxy.
func New(cfg *Config) (*Proxy, error) {
	p := &Proxy{
		cm: &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(cfg.CacheDir),
			Email:  cfg.Email,
			// HostPolicy is enforced by GetConfigForClient.
		},
	}
	p.Reconfigure(cfg)
	return p, nil
}

// Reconfigure updates the proxy's configuration. Some parameters cannot be
// changed after Start has been called, e.g. HTTPAddr, TLSAddr, CacheDir.
func (p *Proxy) Reconfigure(cfg *Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	a, _ := yaml.Marshal(p.cfg)
	b, _ := yaml.Marshal(cfg)
	if bytes.Equal(a, b) {
		return nil
	}
	if p.cfg != nil {
		log.Print("INFO Configuration changed")
	}

	mapping := make(map[string]*Backend, len(cfg.Backends))
	for _, be := range cfg.Backends {
		for _, sn := range be.ServerNames {
			mapping[sn] = be
		}
		if be.ClientAuth {
			tc := &tls.Config{
				GetCertificate: p.cm.GetCertificate,
				ClientAuth:     tls.RequireAndVerifyClientCert,
			}
			if be.ClientCAs != "" {
				c, err := loadCerts(be.ClientCAs)
				if err != nil {
					return err
				}
				tc.ClientCAs = c
			}
			be.tlsConfig = tc
		}
		if be.ForwardRootCAs != "" {
			c, err := loadCerts(be.ForwardRootCAs)
			if err != nil {
				return err
			}
			be.forwardRootCAs = c
		}
	}
	p.mapping = mapping
	p.cfg = cfg
	return nil
}

// Start starts a TLS proxy with the given configuration. The proxy runs
// in background until the context is canceled.
func (p *Proxy) Start(ctx context.Context) error {
	var httpServer *http.Server
	if p.cfg.HTTPAddr != "" {
		httpServer = &http.Server{
			Handler: p.cm.HTTPHandler(nil),
		}
		httpListener, err := net.Listen("tcp", p.cfg.HTTPAddr)
		if err != nil {
			return err
		}
		go func() {
			if err := httpServer.Serve(httpListener); err != http.ErrServerClosed {
				log.Fatalf("http: %v", err)
			}
		}()
	}

	listener, err := tls.Listen("tcp", p.cfg.TLSAddr, &tls.Config{
		GetCertificate: p.cm.GetCertificate,
		GetConfigForClient: func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			be, err := p.backend(hello.ServerName)
			if err != nil {
				return nil, err
			}
			return be.tlsConfig, nil
		},
	})
	if err != nil {
		return err
	}
	p.listener = listener
	p.ctx, p.cancel = context.WithCancel(ctx)

	go func() {
		<-p.ctx.Done()
		p.cancel()
		if httpServer != nil {
			httpServer.Close()
		}
		listener.Close()
	}()

	go func() {
		log.Printf("INFO Accepting TLS connections on %s", p.listener.Addr())
		for {
			conn, err := p.listener.Accept()
			if err != nil {
				log.Printf("ERR Accept: %v", err)
				continue
			}
			go p.handleConnection(conn.(*tls.Conn))
		}
	}()
	return nil
}

// Stop signals the background goroutines to exit.
func (p *Proxy) Stop() {
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
}

func (p *Proxy) incNumOpen(n int) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.numOpen += n
	return p.numOpen
}

func (p *Proxy) handleConnection(extConn *tls.Conn) {
	start := time.Now()
	numOpen := p.incNumOpen(1)
	defer func() {
		extConn.Close()
		p.incNumOpen(-1)
	}()
	if numOpen >= p.cfg.MaxOpen {
		log.Printf("ERR %s: too many open connections: %d >= %d", extConn.RemoteAddr(), numOpen, p.cfg.MaxOpen)
		return
	}
	setKeepAlive(extConn.NetConn())
	hsCtx, cancel := context.WithTimeout(p.ctx, time.Minute)
	defer cancel()
	if err := extConn.HandshakeContext(hsCtx); err != nil {
		log.Printf("ERR %s: %v", extConn.RemoteAddr(), err)
		return
	}
	if err := extConn.Handshake(); err != nil {
		log.Printf("ERR %s: %v", extConn.RemoteAddr(), err)
		return
	}
	tlsTime := time.Since(start)
	cs := extConn.ConnectionState()
	sni := cs.ServerName
	be, err := p.backend(sni)
	if err != nil {
		log.Printf("ERR %s: received unexpected SNI %q", extConn.RemoteAddr(), sni)
		return
	}
	intConn, err := be.dial()
	if err != nil {
		return
	}
	dialTime := time.Since(start)
	defer intConn.Close()
	setKeepAlive(intConn)
	peer := "-"
	if len(cs.PeerCertificates) > 0 {
		peer = cs.PeerCertificates[0].Subject.String()
	}
	desc := fmt.Sprintf("[%s] %s ➔ %s ➔ %s", peer, extConn.RemoteAddr(), sni, intConn.RemoteAddr())
	log.Printf("BEGIN %s", desc)
	ch := make(chan error)
	go func() {
		ch <- forward(extConn, intConn)
	}()
	if err := forward(intConn, extConn); err != nil {
		log.Printf("ERR %s [ext ➔ int]: %v", desc, err)
	}
	if err := <-ch; err != nil {
		log.Printf("ERR %s [int ➔ ext]: %v", desc, err)
	}
	log.Printf("END   %s; HS:%s Dial:%s Total:%s", desc, tlsTime, dialTime-tlsTime, time.Since(start))
}

func setKeepAlive(conn net.Conn) {
	switch c := conn.(type) {
	case *tls.Conn:
		setKeepAlive(c.NetConn())
	case *net.TCPConn:
		c.SetKeepAlivePeriod(time.Minute)
		c.SetKeepAlive(true)
	}
}

func loadCerts(s string) (*x509.CertPool, error) {
	var b []byte
	if len(s) > 0 && s[0] == '/' {
		var err error
		if b, err = os.ReadFile(s); err != nil {
			return nil, err
		}
	} else {
		b = []byte(s)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(b) {
		return nil, errors.New("invalid certs")
	}
	return pool, nil
}

type timeout struct {
	net.Conn
}

func (t timeout) Read(b []byte) (int, error) {
	t.Conn.SetReadDeadline(time.Now().Add(15 * time.Minute))
	return t.Conn.Read(b)
}

func (t timeout) Write(b []byte) (int, error) {
	t.Conn.SetWriteDeadline(time.Now().Add(1 * time.Minute))
	return t.Conn.Write(b)
}

func forward(out net.Conn, in net.Conn) error {
	if _, err := io.Copy(timeout{out}, timeout{in}); err != nil {
		if c, ok := out.(io.Closer); ok {
			c.Close()
		}
		if c, ok := in.(io.Closer); ok {
			c.Close()
		}
		return err
	}

	type closeWriter interface {
		CloseWrite() error
	}
	type closeReader interface {
		CloseRead() error
	}
	if c, ok := out.(closeWriter); ok {
		c.CloseWrite()
	}
	if c, ok := in.(closeReader); ok {
		c.CloseRead()
	}
	return nil
}

func (p *Proxy) backend(sni string) (*Backend, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	be, ok := p.mapping[sni]
	if !ok {
		return nil, fmt.Errorf("unexpected SNI: %q", sni)
	}
	return be, nil
}

func (be *Backend) dial() (net.Conn, error) {
	var max int
	for {
		be.mu.Lock()
		sz := len(be.Addresses)
		if max == 0 {
			max = sz
		}
		addr := be.Addresses[be.next]
		be.next = (be.next + 1) % sz
		be.mu.Unlock()

		dialer := &net.Dialer{
			Timeout:   be.ForwardTimeout,
			KeepAlive: 30 * time.Second,
		}
		var c net.Conn
		var err error
		if be.UseTLS {
			tc := &tls.Config{
				InsecureSkipVerify: be.InsecureSkipVerify,
				ServerName:         be.ForwardServerName,
			}
			if be.forwardRootCAs != nil {
				tc.RootCAs = be.forwardRootCAs
			}
			c, err = tls.DialWithDialer(dialer, "tcp", addr, tc)
		} else {
			c, err = dialer.Dial("tcp", addr)
		}
		if err != nil {
			log.Printf("ERR dial %q: %v", addr, err)
			max--
			if max > 0 {
				continue
			}
			return nil, err
		}
		return c, nil
	}
}
