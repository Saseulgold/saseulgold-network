package swift

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

type SecurityConfig struct {
	CertFile   string
	KeyFile    string
	CACertFile string
	UseTLS     bool
}

func newTLSConfig(config SecurityConfig) (*tls.Config, error) {
	if !config.UseTLS {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(config.CACertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
