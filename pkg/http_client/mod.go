package http_client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type HttpConf struct {
	HttpClientTimeout             int  `mapstructure:"HTTP_CLIENT_TIMEOUT" json:"http_client_timeout"`
	HttpClientIdleConnTimeout     int  `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT" json:"http_client_idle_conn_timeout"`
	HttpClientMaxIdleConns        int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS" json:"http_client_max_idle_conns"`
	HttpClientMaxIdleConnsPerHost int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PERHOST" json:"http_client_max_idle_conns_per_host"`
	HttpClientMaxConnsPerHost     int  `mapstructure:"HTTP_CLIENT_MAX_CONNS_PERHOST" json:"http_client_max_conns_per_host"`
	HttpClientPrintCurl           bool `mapstructure:"HTTP_CLIENT_PRINT_CURL" json:"http_client_print_curl"`
}

var Settings = &HttpConf{
	HttpClientTimeout:             1500,
	HttpClientIdleConnTimeout:     90,
	HttpClientMaxIdleConns:        100,
	HttpClientMaxIdleConnsPerHost: 100,
	HttpClientMaxConnsPerHost:     100,
	HttpClientPrintCurl:           true,
}

func Init() {
	Log.Infof("HttpConf: timeout=%ds", Settings.HttpClientTimeout)
}

func defaultTransport() http.RoundTripper {
	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(Settings.HttpClientTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(Settings.HttpClientIdleConnTimeout) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        Settings.HttpClientMaxIdleConns,
		IdleConnTimeout:     time.Duration(150) * time.Second,
		MaxIdleConnsPerHost: Settings.HttpClientMaxIdleConnsPerHost,
		MaxConnsPerHost:     Settings.HttpClientMaxConnsPerHost,
		DisableCompression:  true,
		DisableKeepAlives:   false, // 默认选项
	}
}
