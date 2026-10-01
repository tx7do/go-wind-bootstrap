// Package tlsutil builds [tls.Config] values from the bootstrap conf TLS
// messages, shared by the registry (and other client-side) adapters.
//
// The conf-level TLS message accepts either file paths or inline PEM bytes;
// when both are provided the PEM bytes win.
package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// ClientTLS builds a client-side [tls.Config] from a [v1.Registry_TLS].
// Returns (nil, nil) when t is nil or carries no certificate material,
// letting the underlying SDK fall back to its own defaults.
func ClientTLS(t *v1.Registry_TLS) (*tls.Config, error) {
	if t == nil {
		return nil, nil
	}

	var certPEM, keyPEM, caPEM []byte

	if f := t.GetFile(); f != nil {
		for _, load := range []struct {
			path string
			dst  *[]byte
		}{
			{f.GetCertPath(), &certPEM},
			{f.GetKeyPath(), &keyPEM},
			{f.GetCaPath(), &caPEM},
		} {
			if load.path == "" {
				continue
			}
			b, err := os.ReadFile(load.path)
			if err != nil {
				return nil, fmt.Errorf("tlsutil: read %q: %w", load.path, err)
			}
			*load.dst = b
		}
	}

	// Inline PEM overrides file content.
	if c := t.GetConfig(); c != nil {
		if len(c.GetCertPem()) > 0 {
			certPEM = c.GetCertPem()
		}
		if len(c.GetKeyPem()) > 0 {
			keyPEM = c.GetKeyPem()
		}
		if len(c.GetCaPem()) > 0 {
			caPEM = c.GetCaPem()
		}
	}

	if len(certPEM) == 0 && len(keyPEM) == 0 && len(caPEM) == 0 && !t.GetInsecureSkipVerify() {
		return nil, nil
	}

	cfg := &tls.Config{
		InsecureSkipVerify: t.GetInsecureSkipVerify(), //nolint:gosec // 显式由配置开启
	}

	if len(certPEM) > 0 || len(keyPEM) > 0 {
		pair, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			return nil, fmt.Errorf("tlsutil: parse key pair: %w", err)
		}
		cfg.Certificates = []tls.Certificate{pair}
	}

	if len(caPEM) > 0 {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("tlsutil: append CA certificate: no valid PEM block")
		}
		cfg.RootCAs = pool
	}

	return cfg, nil
}
