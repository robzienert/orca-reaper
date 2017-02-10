package orcareaper

import (
	"errors"
	"flag"
	"time"
)

type Config struct {
	APIBaseURL       string
	Region           string
	Credentials      string
	InactiveDuration time.Duration
}

func ParseConfig() (*Config, error) {
	var apiBaseURL string
	var region string
	var credentials string

	flag.StringVar(&apiBaseURL, "apiBaseURL", "", "The base URL for Gate")
	flag.StringVar(&region, "region", "us-west-2", "The region where Spinnaker is running")
	flag.StringVar(&credentials, "credentials", "", "The account Spinnaker is running in")
	flag.Parse()

	if apiBaseURL == "" {
		return nil, errors.New("apiBaseURL must be defined")
	}
	if credentials == "" {
		return nil, errors.New("credentials must be defined")
	}

	return &Config{
		APIBaseURL:       apiBaseURL,
		Region:           region,
		Credentials:      credentials,
		InactiveDuration: time.Hour * 4,
	}, nil
}
