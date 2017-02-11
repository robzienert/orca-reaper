package orcareaper

import (
	"errors"
	"flag"
	"time"
)

type Config struct {
	APIBaseURL       string
	Cluster          string
	Region           string
	Credentials      string
	InactiveDuration time.Duration
	DryRun           bool
	X509CertPath     string
	X509KeyPath      string
}

func (c *Config) TLSEnabled() bool {
	return c.X509CertPath != "" && c.X509KeyPath != ""
}

func ParseConfig() (*Config, error) {
	c := Config{}

	flag.StringVar(&c.APIBaseURL, "apiBaseURL", "", "Base URL for Gate")
	flag.StringVar(&c.Region, "region", "us-west-2", "Region where Spinnaker is running")
	flag.StringVar(&c.Credentials, "credentials", "", "Account Spinnaker is running in")
	flag.StringVar(&c.X509CertPath, "x509CertPath", "", "x509 certificate path")
	flag.StringVar(&c.X509KeyPath, "x509KeyPath", "", "x509 key path")
	flag.StringVar(&c.Cluster, "cluster", "", "Limit scope of reaping to a cluster (useful for multiple Spinnaker deploys)")
	flag.BoolVar(&c.DryRun, "dryRun", false, "When set, no reap tasks will be submitted to Orca")
	flag.Parse()

	if c.APIBaseURL == "" {
		return nil, errors.New("apiBaseURL must be defined")
	}
	if c.Credentials == "" {
		return nil, errors.New("credentials must be defined")
	}
	if c.X509CertPath != "" && c.X509KeyPath == "" {
		return nil, errors.New("x509 key path must be supplied with cert")
	}
	if c.X509KeyPath != "" && c.X509CertPath == "" {
		return nil, errors.New("x509 cert path must be supplied with key")
	}

	return &c, nil
}
