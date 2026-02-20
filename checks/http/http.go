// Package http provides an HTTP check for checking the status of a web endpoint
package http

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	// UserAgent default
	UserAgent = "Go Up Checker"
)

// Check returns a test
func Check(v *viper.Viper) error {

	url := v.GetString("Endpoint")

	if url == "" {
		return errors.New("no endpoint set")
	}

	method := strings.ToUpper(v.GetString("Method"))
	if method == "" {
		method = "HEAD"
	}

	expectedCode := v.GetInt("Status")
	if expectedCode == 0 {
		expectedCode = 200
	}

	client := DefaultHTTPClient

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain and discard the response body to properly close the connection
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != expectedCode {
		return fmt.Errorf("expected status %d, received %d", expectedCode, resp.StatusCode)
	}

	return nil
}

// DefaultHTTPClient is used when no other http.Client
// is specified on a Checker.
var DefaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 0,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   1,
		DisableCompression:    true,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: 5 * time.Second,
	},
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Timeout: 10 * time.Second,
}
