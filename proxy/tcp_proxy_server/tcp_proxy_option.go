package tcp_proxy_server

import "time"

type ProxyOption func(*TCPReverseProxy) *TCPReverseProxy

func WithDialTimeout(d time.Duration) ProxyOption {
	return func(tcp *TCPReverseProxy) *TCPReverseProxy {
		tcp.dialTimeout = d
		return tcp
	}
}

func WithKeepAlivePeriod(d time.Duration) ProxyOption {
	return func(tcp *TCPReverseProxy) *TCPReverseProxy {
		tcp.keepAlivePeriod = d
		return tcp
	}
}
