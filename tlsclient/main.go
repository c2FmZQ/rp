// MIT License
//
// Copyright (c) 2023 TTBT Enterprises LLC
// Copyright (c) 2023 Robin Thellend <rthellend@rthellend.com>
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

// Command tlsproxy establishes a TLS connection with a TLS server and redirects
// the stream to its stdin and stdout.
package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/quic-go/quic-go"
	"golang.org/x/crypto/ocsp"
)

// Version is set with -ldflags="-X main.Version=${VERSION}"
var Version = "dev"

func main() {
	versionFlag := flag.Bool("v", false, "Show the version.")
	key := flag.String("key", "", "A file that contains the TLS key to use.")
	cert := flag.String("cert", "", "A file that contains the TLS certificate to use.")
	alpn := flag.String("alpn", "", "The ALPN proto to request.")
	ech := flag.String("ech", "", "Use this ECH ConfigList.")
	useQUIC := flag.Bool("quic", false, "Use QUIC.")
	verifyOCSP := flag.Bool("ocsp", false, "Require stapled OCSP response.")
	flag.Parse()

	if *versionFlag {
		os.Stdout.WriteString(Version + " " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n")
		return
	}
	if flag.NArg() != 1 || (*key == "") != (*cert == "") {
		os.Stderr.WriteString("Usage: tlsclient [-key=<keyfile> -cert=<certfile>] [-alpn=<proto>] host:port\n")
		os.Exit(1)
	}
	addr := flag.Arg(0)

	var certs []tls.Certificate
	if *key != "" && *cert != "" {
		c, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.Fatalf("ERR: %v", err)
		}
		certs = append(certs, c)
	}

	var protos []string
	if *alpn != "" {
		protos = append(protos, *alpn)
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
		port = "443"
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		res := &net.Resolver{
			PreferGo: true,
			Dial: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("udp", "8.8.8.8:53")
			},
		}
		addrs, err = res.LookupHost(context.Background(), host)
	}
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	if len(addrs) == 0 {
		log.Fatalf("ERR: cannot resolve %s", host)
	}
	target := net.JoinHostPort(addrs[0], port)

	tc := &tls.Config{
		Certificates: certs,
		NextProtos:   protos,
		ServerName:   host,
		VerifyConnection: func(cs tls.ConnectionState) error {
			if !*verifyOCSP {
				return nil
			}
			if len(cs.OCSPResponse) == 0 {
				return errors.New("no ocsp response")
			}
			cert := cs.PeerCertificates[0]
			issuer := cert
			if len(cs.PeerCertificates) > 1 {
				issuer = cs.PeerCertificates[1]
			}
			resp, err := ocsp.ParseResponseForCert(cs.OCSPResponse, cert, issuer)
			if err != nil {
				return err
			}
			if time.Now().After(resp.NextUpdate) {
				return errors.New("ocsp response is expired")
			}
			if resp.Status != ocsp.Good {
				return errors.New("ocsp response status is not good")
			}
			return nil
		},
	}
	if *ech != "" {
		configList, err := base64.StdEncoding.DecodeString(*ech)
		if err != nil {
			log.Fatalf("ERR: --ech decoding error: %v", err)
		}
		tc.EncryptedClientHelloConfigList = configList
	}

	if *useQUIC {
		conn, err := quic.DialAddr(context.Background(), target, tc, &quic.Config{})
		if err != nil {
			log.Fatalf("ERR: %v", err)
		}
		stream, err := conn.OpenStream()
		if err != nil {
			log.Fatalf("ERR: %v", err)
		}
		go func() {
			if _, err := io.Copy(stream, os.Stdin); err != nil && !errors.Is(err, net.ErrClosed) {
				log.Printf("ERR: %v", err)
			}
			stream.Close()
		}()
		if _, err := io.Copy(os.Stdout, stream); err != nil {
			stream.CancelRead(0)
			log.Printf("ERR: %v", err)
		}
		return
	}

	conn, err := tls.Dial("tcp", target, tc)
	if err != nil {
		var echErr *tls.ECHRejectionError
		if errors.As(err, &echErr) {
			log.Fatalf("ECH RetryConfigList: %s", echErr.RetryConfigList)
		}
		log.Fatalf("ERR Dial: %v", err)
	}
	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("ERR Stdin: %v", err)
		}
		conn.CloseWrite()
	}()
	if _, err := io.Copy(os.Stdout, conn); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("ERR Conn: %v", err)
	}
	conn.Close()
	return
}
