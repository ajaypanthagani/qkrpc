package qkrpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// LoadTLSConfig loads a TLS config for the server using a cert and key file.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"qkrpc"},
	}, nil
}

// LoadClientTLS loads a TLS config for the client that trusts the provided cert.
func LoadClientTLS(certFile string) (*tls.Config, error) {
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(certPEM); !ok {
		return nil, err
	}

	return &tls.Config{
		RootCAs:    pool,
		NextProtos: []string{"qkrpc"},
		ServerName: "localhost",
	}, nil
}
