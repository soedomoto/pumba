package container

import (
	"net/url"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

func HTTPClient(daemonUrl string, tlsConfig *tls.Config) (*http.Client, error) {
	u, err := url.Parse(daemonUrl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Scheme == "tcp" {
		if tlsConfig == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}

	return newHTTPClient(u, tlsConfig, time.Duration(defaultTimeout))
}

func newHTTPClient(url *url.URL, tlsConfig *tls.Config, timeout time.Duration) (*http.Client, error) {
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	switch url.Scheme {
	default:
		httpTransport.Dial = func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout(proto, addr, timeout)
		}
	case "unix":
		socketPath := url.Path
		unixDial := func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, timeout)
		}
		httpTransport.Dial = unixDial
		// Override the main URL object so the HTTP lib won't complain
		url.Scheme = "http"
		url.Host = "unix.sock"
		url.Path = ""
	}
	return &http.Client{Transport: httpTransport}, nil
}
