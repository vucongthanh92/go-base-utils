package client

import (
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	maxIdleConns              = 100
	maxConnsPerHost           = 100
	maxIdleConnsPerHost       = 100
	clientTimeout             = 10 * time.Second
	dialContextTimeout        = 10 * time.Second
	clientTLSHandshakeTimeout = 10 * time.Second
	clientRetryWaitTime       = 300 * time.Millisecond
	retryCount                = 3
)

func NewHttpClient() *resty.Client {
	transport := &http.Transport{
		DialContext:         (&net.Dialer{Timeout: dialContextTimeout}).DialContext,
		MaxIdleConns:        maxIdleConns,
		MaxConnsPerHost:     maxConnsPerHost,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		TLSHandshakeTimeout: clientTLSHandshakeTimeout,
	}

	client := resty.New().
		SetTimeout(clientTimeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetTransport(transport)

	return client
}
