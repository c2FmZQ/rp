# TLSPROXY Release Notes

## next

### :star2: New feature

* Add support for Encrypted Client Hello. This feature improves privacy by allowing the clients to encrypt the Server Name to which they are connecting. Without ECH, this information is actually transmitted in plaintext. When `ech:` is set in `config.yaml`, tlsproxy handles ECH as a Client-Facing Server with a Split Mode Topology as specified in https://datatracker.ietf.org/doc/html/draft-ietf-tls-esni/. See [ECH](https://github.com/c2FmZQ/tlsproxy/blob/main/docs/ECH.md)

### :wrench: Bug fixes

* Handle quic.ErrTransportClosed correctly.

### :wrench: Misc

* Use os.Root for static file isolation
* Allow customized permission denied messages.
* Update go: 1.24rc3
* Update go dependencies:
  * upgraded github.com/quic-go/quic-go v0.48.2 => v0.49.0
  * upgraded github.com/beevik/etree v1.4.1 => v1.5.0
  * upgraded github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad => v0.0.0-20250202011525-fc3143867406
  * upgraded golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 => v0.0.0-20250128182459-e0ece0dbea4c
  * upgraded golang.org/x/mod v0.22.0 => v0.23.0
  * upgraded golang.org/x/sync v0.10.0 => v0.11.0
  * upgraded golang.org/x/sys v0.29.0 => v0.30.0
  * upgraded golang.org/x/text v0.21.0 => v0.22.0
  * upgraded golang.org/x/time v0.9.0 => v0.10.0

## v0.14.2

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/c2FmZQ/storage v0.2.3 => v0.2.4
  * upgraded github.com/c2FmZQ/tpm v0.3.1 => v0.4.0
  * upgraded github.com/google/go-tpm v0.9.1 => v0.9.3
  * upgraded github.com/jonboulle/clockwork v0.4.0 => v0.5.0
  * upgraded github.com/onsi/ginkgo/v2 v2.22.0 => v2.22.2
  * upgraded golang.org/x/crypto v0.31.0 => v0.32.0
  * upgraded golang.org/x/exp v0.0.0-20241210194714-1829a127f884 => v0.0.0-20250106191152-7588d65b2ba8
  * upgraded golang.org/x/net v0.32.0 => v0.34.0
  * upgraded golang.org/x/sys v0.28.0 => v0.29.0
  * upgraded golang.org/x/time v0.8.0 => v0.9.0
  * upgraded golang.org/x/tools v0.28.0 => v0.29.0

## v0.14.1

### :wrench: Misc

* Simplify the PKI client-side code and remove the service worker.

## v0.14.0

### :star2: New feature

* Add a built-in certificate authority for SSH. It is enabled with the top-level `sshCertificateAuthorities` field in the config file.
  * This CA issues SSH user certificates with the current user's email address as both `Key ID` and `Principal`. It only works when SSO is enabled.
  * User authorization is done by adding a line like this to the user's `.ssh/authorized_keys` file:
```
cert-authority,principals="<email>" <CA's public key>
```

### :star: Feature improvement

* Add an option to exclude some path prefixes from SSO enforcement.

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20241203143554-1e3fdc7de467 => v0.0.0-20241210010833-40e02aabc2ad
  * upgraded golang.org/x/crypto v0.29.0 => v0.31.0
  * upgraded golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f => v0.0.0-20241210194714-1829a127f884
  * upgraded golang.org/x/net v0.31.0 => v0.32.0
  * upgraded golang.org/x/sync v0.9.0 => v0.10.0
  * upgraded golang.org/x/sys v0.27.0 => v0.28.0
  * upgraded golang.org/x/text v0.20.0 => v0.21.0
  * upgraded golang.org/x/tools v0.27.0 => v0.28.0

## v0.13.2

### :wrench: Misc

* Update go: 1.23.4
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20241101162523-b92577c0c142 => v0.0.0-20241203143554-1e3fdc7de467
  * upgraded github.com/onsi/ginkgo/v2 v2.21.0 => v2.22.0
  * upgraded github.com/quic-go/quic-go v0.48.1 => v0.48.2

## v0.13.1

### :wrench: Bug fix

* Fix goroutine and connection leak with websockets.

## v0.13.0

### :star2: New feature

* Add support for forwarding WebSocket requests to arbitrary TCP servers. WebSockets were already forwarded transparently to backends before, and that is not changing. The new feature lets tlsproxy itself handle the WebSocket request and forward them to any TCP servers. The content of BinaryMessages is streamed to the remote server, and data received from the server is sent back to the client also as BinaryMessages.
  * This is used by [SSH Term](https://github.com/c2FmZQ/sshterm).

## v0.12.0

### :star: Feature improvement

* Improve passkey authentication to support WebAuthn devices without discoverable credentials.

### :wrench: Misc

* Update go dependencies:
  * upgraded golang.org/x/crypto v0.28.0 => v0.29.0
  * upgraded golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c => v0.0.0-20241108190413-2d47ceb2692f
  * upgraded golang.org/x/mod v0.21.0 => v0.22.0
  * upgraded golang.org/x/net v0.30.0 => v0.31.0
  * upgraded golang.org/x/sync v0.8.0 => v0.9.0
  * upgraded golang.org/x/sys v0.26.0 => v0.27.0
  * upgraded golang.org/x/text v0.19.0 => v0.20.0
  * upgraded golang.org/x/time v0.7.0 => v0.8.0
  * upgraded golang.org/x/tools v0.26.0 => v0.27.0

## v0.11.1

### :wrench: Bug fix

* Add missing lock that might affect the use of `tlsCertificates`.

## v0.11.0

### :star2: New features

* Add a configuration option to filter out logged data. See `LogFilter` in [config.go](https://github.com/c2FmZQ/tlsproxy/blob/main/proxy/config.go). This can be set at the top level config, the backend level config, or both.
* Add support for TLS certificates stored locally. See `tlsCertificates` in [config.go](https://github.com/c2FmZQ/tlsproxy/blob/main/proxy/config.go).

### :wrench: Misc

* Update go: 1.23.3
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20241008150032-332c0e1a4a34 => v0.0.0-20241101162523-b92577c0c142
  * upgraded github.com/onsi/ginkgo/v2 v2.20.2 => v2.21.0
  * upgraded go.uber.org/mock v0.4.0 => v0.5.0
  * upgraded golang.org/x/exp v0.0.0-20241004190924-225e2abe05e6 => v0.0.0-20241009180824-f66d83c29e7c
* Testing:
  * When testing with the `--use-ephemeral-certificate-manager` flag, the ephemeral CA cert and key will be saved if the `CERTMANAGER_STATE_FILE` environment variable is set.

## v0.10.9

### :wrench: Bug fix

* Another fix for `content-length: -1`. The change in v0.10.7 broke HTTP POST requests.

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/quic-go/quic-go v0.47.0 => v0.48.1

## v0.10.8

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240910150728-a0b0bb1d4134 => v0.0.0-20241008150032-332c0e1a4a34
  * upgraded github.com/pires/go-proxyproto v0.7.0 => v0.8.0
  * upgraded golang.org/x/crypto v0.27.0 => v0.28.0
  * upgraded golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 => v0.0.0-20241004190924-225e2abe05e6
  * upgraded golang.org/x/net v0.29.0 => v0.30.0
  * upgraded golang.org/x/sys v0.25.0 => v0.26.0
  * upgraded golang.org/x/text v0.18.0 => v0.19.0
  * upgraded golang.org/x/time v0.6.0 => v0.7.0
  * upgraded golang.org/x/tools v0.25.0 => v0.26.0

## v0.10.7

### :wrench: Bug fix

* Don't send `content-length: -1` to backends. This caused `400` errors in some configurations. This bug was introduced in v0.10.6.

### :wrench: Misc

* Update go: 1.23.2

## v0.10.6

### :wrench: Misc

* Update to quic-go [v0.47.0](https://github.com/quic-go/quic-go/releases/tag/v0.47.0). The release notes point out that a bug in go 1.23 is causing problems with quic. So, we're also setting go version in `go.mod` back to `1.22.0` for now.
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240903155634-a8630aee4ab9 => v0.0.0-20240910150728-a0b0bb1d4134
  * upgraded github.com/quic-go/qpack v0.5.0 => v0.5.1
  * upgraded github.com/quic-go/quic-go v0.46.0 => v0.47.0
  * upgraded golang.org/x/exp v0.0.0-20240904232852-e7e105dedf7e => v0.0.0-20240909161429-701f63a606c0
  * upgraded golang.org/x/tools v0.24.0 => v0.25.0

## v0.10.5

### :wrench: Misc

* Update go: 1.23.1
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 => v0.0.0-20240903155634-a8630aee4ab9
  * upgraded github.com/onsi/ginkgo/v2 v2.19.1 => v2.20.2
  * upgraded github.com/quic-go/qpack v0.4.0 => v0.5.0
  * upgraded golang.org/x/crypto v0.26.0 => v0.27.0
  * upgraded golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 => v0.0.0-20240904232852-e7e105dedf7e
  * upgraded golang.org/x/mod v0.20.0 => v0.21.0
  * upgraded golang.org/x/net v0.28.0 => v0.29.0
  * upgraded golang.org/x/sys v0.23.0 => v0.25.0
  * upgraded golang.org/x/text v0.17.0 => v0.18.0
  * upgraded software.sslmate.com/src/go-pkcs12 v0.4.0 => v0.5.0

## v0.10.4

### :wrench: Misc

* Update go: 1.22.6
* Update go dependencies:
  * upgraded github.com/beevik/etree v1.4.0 => v1.4.1
  * upgraded github.com/google/pprof v0.0.0-20240625030939-27f56978b8b0 => v0.0.0-20240727154555-813a5fbdbec8
  * upgraded github.com/onsi/ginkgo/v2 v2.19.0 => v2.19.1
  * upgraded github.com/quic-go/quic-go v0.45.2 => v0.46.0
  * upgraded golang.org/x/crypto v0.25.0 => v0.26.0
  * upgraded golang.org/x/exp v0.0.0-20240707233637-46b078467d37 => v0.0.0-20240719175910-8a7402abbf56
  * upgraded golang.org/x/mod v0.19.0 => v0.20.0
  * upgraded golang.org/x/net v0.27.0 => v0.28.0
  * upgraded golang.org/x/sync v0.7.0 => v0.8.0
  * upgraded golang.org/x/sys v0.22.0 => v0.23.0
  * upgraded golang.org/x/text v0.16.0 => v0.17.0
  * upgraded golang.org/x/time v0.5.0 => v0.6.0
  * upgraded golang.org/x/tools v0.23.0 => v0.24.0

## v0.10.3

### :wrench: Misc

* Pick up bug fixes in the quic-go package.
* Update go dependencies:
  * upgraded github.com/quic-go/quic-go v0.45.1 => v0.45.2
  * upgraded golang.org/x/crypto v0.24.0 => v0.25.0
  * upgraded golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 => v0.0.0-20240707233637-46b078467d37
  * upgraded golang.org/x/mod v0.18.0 => v0.19.0
  * upgraded golang.org/x/net v0.26.0 => v0.27.0
  * upgraded golang.org/x/sys v0.21.0 => v0.22.0
  * upgraded golang.org/x/tools v0.22.0 => v0.23.0

## v0.10.2

### :wrench: Misc

* Same as v0.10.0. The v0.10.0 docker image was created before the golang linux/amd64 was ready. So, it used linux/386 instead. v0.10.2 should be OK.

## v0.10.1

### :wrench: Misc

* Same as v0.10.0. ~The v0.10.0 docker image was created before the golang linux/amd64 was ready. So, it used linux/386 instead. v0.10.1 should be OK.~

## v0.10.0

### :star: Feature improvements

* When `forwardHttpHeaders` is used, special keywords are automatically expanded from the header values:
  * `${NETWORK}` is either tcp or udp.
  * `${LOCAL_ADDR}` is the local address of the network connection.
  * `${REMOTE_ADDR}` is the remote address of the network connection.
  * `${LOCAL_IP}` is the local IP address of the network connection.
  * `${REMOTE_IP}` is the remote IP address of the network connection.
  * `${SERVER_NAME}` is the server name requested by the client.
  * `${JWT:xxxx}` expands to the value of claim `xxxx` from the ID token.

### :wrench: Misc

* Update go: 1.22.5
* Update go dependencies:
  * upgraded github.com/fxamacker/cbor/v2 v2.6.0 => v2.7.0
  * upgraded github.com/google/pprof v0.0.0-20240528025155-186aa0362fba => v0.0.0-20240625030939-27f56978b8b0
  * upgraded golang.org/x/exp v0.0.0-20240604190554-fc45aab8b7f8 => v0.0.0-20240613232115-7f521ea00fb8

## v0.9.1

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/quic-go/quic-go v0.45.0 => v0.45.1

## v0.9.0

### :star: Feature improvements

* Add `forwardHttpHeaders` to set HTTP headers in the forwarded HTTP requests. Headers that already exist are overwritten.

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/c2FmZQ/storage v0.2.2 => v0.2.3
  * upgraded github.com/c2FmZQ/tpm v0.3.0 => v0.3.1
  * upgraded github.com/google/go-tpm-tools v0.4.3 => v0.4.4
  * upgraded github.com/google/go-tpm v0.9.0 => v0.9.1
  * upgraded github.com/quic-go/quic-go v0.44.0 => v0.45.0

## v0.8.3

### :wrench: Misc

* Update go: 1.22.4
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240509144519-723abb6459b7 => v0.0.0-20240528025155-186aa0362fba
  * upgraded github.com/onsi/ginkgo/v2 v2.17.3 => v2.19.0
  * upgraded golang.org/x/crypto v0.23.0 => v0.24.0
  * upgraded golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 => v0.0.0-20240604190554-fc45aab8b7f8
  * upgraded golang.org/x/mod v0.17.0 => v0.18.0
  * upgraded golang.org/x/net v0.25.0 => v0.26.0
  * upgraded golang.org/x/sys v0.20.0 => v0.21.0
  * upgraded golang.org/x/text v0.15.0 => v0.16.0
  * upgraded golang.org/x/tools v0.21.0 => v0.22.0

## v0.8.2

### :wrench: Bug fix

* Sign OCSP responses with RSA or ECDSA keys.

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/beevik/etree v1.3.0 => v1.4.0
  * upgraded github.com/google/pprof v0.0.0-20240507183855-6f11f98ebb1c => v0.0.0-20240509144519-723abb6459b7
  * upgraded github.com/quic-go/quic-go v0.43.1 => v0.44.0

## v0.8.1

### :star: Feature improvements

* Only allow GET and HEAD methods for static files.
* Sanitize the request path before comparing to local endpoints, e.g. `//.sso` redirects to `/.sso`
* Add a `sanitizePath` option to backends. When true (default), request paths are sanitized before they are sent to the backends.

### :wrench: Misc

* Update go: 1.22.3
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240422182052-72c8669ad3e7 => v0.0.0-20240507183855-6f11f98ebb1c
  * upgraded github.com/onsi/ginkgo/v2 v2.17.1 => v2.17.3
  * upgraded golang.org/x/crypto v0.22.0 => v0.23.0
  * upgraded golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f => v0.0.0-20240506185415-9bf2ced13842
  * upgraded golang.org/x/net v0.24.0 => v0.25.0
  * upgraded golang.org/x/sys v0.19.0 => v0.20.0
  * upgraded golang.org/x/text v0.14.0 => v0.15.0
  * upgraded golang.org/x/tools v0.20.0 => v0.21.0

## v0.8.0

### :star2: New features

* Serve static files from a local filesystem when `documentRoot:` is set.

### :wrench: Misc

* Upgrade github.com/quic-go/quic-go v0.42.0 => v0.43.1

## v0.7.2

### :wrench: Misc

* Update the tpm library to pick up a bug fix. The saved TPM keys would become invalid after a reboot. This only affected configurations with `hwBacked: true`.

## v0.7.1

### :wrench: Misc

* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20240402174815-29b9bb013b0f => v0.0.0-20240422182052-72c8669ad3e7
  * upgraded golang.org/x/crypto v0.21.0 => v0.22.0
  * upgraded golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 => v0.0.0-20240416160154-fe59bbe5cc7f
  * upgraded golang.org/x/mod v0.16.0 => v0.17.0
  * upgraded golang.org/x/net v0.23.0 => v0.24.0
  * upgraded golang.org/x/sys v0.18.0 => v0.19.0
  * upgraded golang.org/x/tools v0.19.0 => v0.20.0

## v0.7.0

### :star2: New features

* Add `hwBacked` option. When enabled, hardware-backed cryptographic keys are used to:
  * encrypt local data (the data cannot be used or recovered on a different device),
  * sign authentication tokens,
  * sign the PKI certificates, OCSP responses, and CRLs.
* Add `--quiet` flag. When set (or the `TLSPROXY_QUIET` env variable is `true`), logging is turned off after tlsproxy starts.

### :wrench: Misc

* Release binaries and container images are now signed.
* Update go: 1.22.2
* Update go dependencies:
  * upgraded github.com/quic-go/quic-go v0.41.0 => v0.42.0

## v0.6.4

* Update go: 1.22.1
* Update go dependencies

## v0.6.3

* Reduce lock contention in passkey manager.

## v0.6.2

* Update go: 1.22.0
* Log HTTP protocol upgrades.
* Fix bug that prevented PKI admins from revoking certificates.
* Improve heap & lock profiling.
* Reduce lock contention.

## v0.6.1

* Fix WebSocket connections when `backendProto` is set to something other than `http/1.1`.

## v0.6.0

* Starting with v0.6.0, tlsproxy is built with QUIC & HTTP/3 support by default. Binaries without QUIC can still be built with `-tags noquic`.
* Add support for QUIC datagrams.

## v0.5.2

* Report outbound connections on the metrics page.
* Improve how http3 connections are handled.
* Update go dependencies:
  * upgraded github.com/google/pprof v0.0.0-20231229205709-960ae82b1e42 => v0.0.0-20240117000934-35fc243c5815
  * upgraded github.com/onsi/ginkgo/v2 v2.13.2 => v2.14.0
  * upgraded github.com/quic-go/quic-go v0.40.1 => v0.41.0
  * upgraded golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc => v0.0.0-20240112132812-db7319d0e0e3
  * upgraded golang.org/x/tools v0.16.1 => v0.17.0

## v0.5.1

* Prevent backend connections from being re-used by other clients when PROXY protocol is enabled.

## v0.5.0

* Add support for the PROXY protocol for incoming connections. See the `AcceptProxyHeaderFrom` config option.
* Fix bug with the handling of the PROXY protocol header with TLS backends. The header was sent after the TLS handshake instead of before.
* Fix bug that prevented logins with passkeys on non default ports when ForceReAuth is set. (introduced in v0.4.4)
* Log aborted ReverseProxy requests more gracefully.
* Update go: 1.21.6
* Update go dependencies:
  * upgraded github.com/beevik/etree v1.2.0 => v1.3.0
  * upgraded github.com/google/pprof v0.0.0-20231212022811-ec68065c825e => v0.0.0-20231229205709-960ae82b1e42
  * upgraded golang.org/x/crypto v0.17.0 => v0.18.0
  * upgraded golang.org/x/exp v0.0.0-20231226003508-02704c960a9b => v0.0.0-20240103183307-be819d1f06fc
  * upgraded golang.org/x/net v0.19.0 => v0.20.0
  * upgraded golang.org/x/sys v0.15.0 => v0.16.0

## v0.4.4

* Improve user experience with the passkey login screen.
* Fix small accounting bug in ingress metrics.
* Apply the forward rate limit to http requests. Before, the rate limit was only applied to incoming connections.

## v0.4.3

* Add support for the OIDC `hd` param (https://developers.google.com/identity/openid-connect/openid-connect#hd-param).
* Return the proper TLS error when a QUIC client requests an unknown server name and/or alpn protocol.
* Return the proper TLS error when a client certificate is revoked.
* Refactor / improve the handling of QUIC connections.

## v0.4.2

* Minor changes to the metrics page: combine runtime and memory profile, add config.
* Remove the /config endpoint.
* Refactor HTTP handlers. No functional change.

## v0.4.1

* Add a confirmation page before redirecting to the remote identity provider's OIDC or SAML login page.
* Add a nonce cookie to the OIDC login flow.
* Convert the metrics page to HTML.

## v0.4.0

* This version changes these config options:
  * `enableQUIC` now defaults to `true` if the binary is compiled with QUIC support, e.g. `+quic` releases.
  * `alpnProtos` now includes `h3` by default when QUIC is enabled and mode is one of `HTTP`, `HTTPS`, `QUIC`, `LOCAL`, or `CONSOLE`.
  * `revokeUnusedCertificates` now defaults to `true`, and revocation is only done at the time the proxy starts.

## v0.3.5

* Add missing lock in revoke.go
* Fix bug that could close legitimate connections after a config change.
* Fixed some flaky tests.
* Add backendProto and buildinfo to metrics page.

## v0.3.4

* Add an option to automatically revoke unused _Let's encrypt_ certificates.
* Add a flag to revoke all cached certificates. This is useful when a server is decommissioned or compromised.
* IDN fixes and tests

## v0.3.3

* Implement OCSP stapling.
* Verify the status of client and server certificates when OCSP servers are specified.
* Fix bug that prevented the PKI WASM client from loading on chrome 119.
* Handle Internationalized Domain Names correctly.
* Add an option to refresh the identity used with passkeys (`refreshInterval`).

## v0.3.2

* Use the certificate provided by _Let's Encrypt_ to authenticate with backends when backends require client authentication.
* go: upgraded github.com/quic-go/quic-go v0.40.0 => v0.40.1 (CVE-2023-49295)

## v0.3.1

* Decouple frontend and backend protocols in `mode: http` and `mode: https`. The ALPN protocols that TLSPROXY accepts are specified in `alpnProtos: [...]`, and the protocol to use with the request to the backend is specified with `backendProto: ...`, which defaults to `http/1.1`. The previous behavior was to use the same protocol that the client used. It is still possible to do that by setting `backendProto` explicitly to an empty string.
* Fix QUIC stream direction in metrics.
* Fix persistent config reload every 30 sec in some conditions.

## v0.3.0

* Add support for the PROXY protocol defined by HAProxy. https://github.com/haproxy/haproxy/blob/master/doc/proxy-protocol.txt
* Add QUIC and HTTP/3 support. See [docs/QUIC.md](https://github.com/c2FmZQ/tlsproxy/blob/main/docs/QUIC.md).
* Increase timeouts for proxied HTTP requests to 24 hours for large uploads or downloads.
* Update to go 1.21.5

## v0.2.0

* Add support for the HTTP/2 protocol in reverse proxy mode. This only works when the backend actually supports HTTP/2. So, it is not enabled by default. To enable it, set `alpnProtos: [h2, http/1.1]` on the backend explicitly.
* Add support for routing HTTP requests based on path. By default, all requests go to the same backend servers. When `pathOverrides` is set, some paths can be routed to other servers.

## v0.1.1

* Allow multiple backends with the same server name but different ALPN protos, e.g. one backend could have foo.example.com with the default ALPN protos, and another backend could have foo.example.com with `alpnProtos: [imap]`. If ALPN is not used by the client, the first backend with a matching server name will be used.

## v0.1.0

Let's call this the first _stable_ development release. We'll try to keep decent release notes going forward.

TLSPROXY is primarily a [TLS termination proxy](https://en.wikipedia.org/wiki/TLS_termination_proxy) that uses letsencrypt to provide TLS encryption for any number of TCP or HTTP servers, and any number of server names concurrently on the same port.

Its functionality is similar to an [stunnel](https://www.stunnel.org/) server, but without the need to configure and run [certbot](https://certbot.eff.org/) separately.

TLSPROXY can also be used as a [Reverse Proxy](https://en.wikipedia.org/wiki/Reverse_proxy) for HTTP(S) services, and optionally control access to these services with user authentication and authorization.

Overview of features:

* [x] Use [Let's Encrypt](https://letsencrypt.org/) automatically to get TLS certificates (http-01 & tls-alpn-01 challenges).
* [x] Terminate TLS connections, and forward the data to any TCP server in plaintext.
* [x] Terminate TLS connections, and forward the data to any TLS server. The data is encrypted in transit, but the proxy sees the plaintext.
* [x] Terminate _TCP_ connections, and forward the TLS connection to any TLS server (passthrough). The proxy doesn't see the plaintext.
* [x] Terminate HTTPS connections, and forward the requests to HTTP or HTTPS servers (http/1 only, not recommended with c2fmzq-server).
* [x] TLS client authentication & authorization (when the proxy terminates the TLS connections).
* [x] Built-in Certificate Authority for managing client and backend server TLS certificates.
* [x] User authentication with OpenID Connect, SAML, and/or passkeys (for HTTP and HTTPS connections). Optionally issue JSON Web Tokens (JWT) to authenticated users to use with the backend s
ervices and/or run a local OpenID Connect server for backend services.
* [x] Access control by IP address.
* [x] Routing based on Server Name Indication (SNI), with optional default route when SNI isn't used.
* [x] Simple round-robin load balancing between servers.
* [x] Support any ALPN protocol in TLS, TLSPASSTHROUGH, or TCP mode.
* [x] Use the same TCP address (IPAddr:port) for any number of server names, e.g. foo.example.com and bar.example.com on the same xxx.xxx.xxx.xxx:443.

