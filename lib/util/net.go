package util

import (
	"crypto/tls"
	"net/http"
	"time"
)

const defaultHttpTimeout = 10 * time.Second

func DefaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: defaultHttpTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
