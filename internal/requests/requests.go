package requests

// TODO: Handle http:// & http:// checks in urls
// TODO: Handle default to https if nothing found

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	REQUEST_TIMEOUT         = 10
	MAX_IDLE_CONNS          = 100
	MAX_IDLE_CONNS_PER_HOST = 100
	IDLE_CONN_TIMEOUT       = 90
)

type ApiRequestClient struct {
	client *http.Client
}

func NewApiRequestClient() *ApiRequestClient {
	return &ApiRequestClient{
		client: &http.Client{
			Timeout: REQUEST_TIMEOUT * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        MAX_IDLE_CONNS,
				MaxIdleConnsPerHost: MAX_IDLE_CONNS_PER_HOST,
				IdleConnTimeout:     IDLE_CONN_TIMEOUT * time.Second,
			},
		},
	}
}

func (c *ApiRequestClient) Fetch(url string) (respContent string, err error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("error closing body: %v\n", closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rawResp), err
}
