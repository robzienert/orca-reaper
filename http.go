package orcareaper

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"net/http"
	"time"

	"github.com/pkg/errors"
)

const httpClientTimeout = time.Second * 10

var httpClient *http.Client

func initHTTPClient(config *Config) error {
	if !config.TLSEnabled() {
		log.Println("TLS not enabled!")
		httpClient = &http.Client{
			Timeout: httpClientTimeout,
		}
		return nil
	}

	cert, err := tls.LoadX509KeyPair(config.X509CertPath, config.X509KeyPath)
	if err != nil {
		return errors.Wrap(err, "loading x509 keypair")
	}

	clientCACert, err := ioutil.ReadFile(config.X509CertPath)
	if err != nil {
		return errors.Wrap(err, "loading client ca cert")
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()

	httpClient = &http.Client{
		Timeout: httpClientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return nil
}
