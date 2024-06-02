package pocScan

import (
	"changeme/services/publicCode"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Client           *http.Client
	ClientNoRedirect *http.Client
	dialTimout       = 5 * time.Second
	keepAlive        = 5 * time.Second
)

func CreatPocScanClient() {

	// 设置超时时间
	timeout := time.Duration(publicCode.TimeOut) * time.Second

	type DialContext = func(ctx context.Context, network, addr string) (net.Conn, error)
	dialer := &net.Dialer{
		Timeout:   dialTimout,
		KeepAlive: keepAlive,
	}

	//transport := &http.Transport{
	//	DialContext:         dialer.DialContext,
	//	MaxConnsPerHost:     5,
	//	MaxIdleConns:        0,
	//	MaxIdleConnsPerHost: publicCode.ThreadNum * 2,
	//	IdleConnTimeout:     keepAlive,
	//	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	//	TLSHandshakeTimeout: 5 * time.Second,
	//	DisableKeepAlives:   false,
	//	DisableCompression:  true,
	//}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxConnsPerHost:     5,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: publicCode.ThreadNum * 2,
		IdleConnTimeout:     keepAlive,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS13, // 或者指定一个更高的版本
			InsecureSkipVerify: true,             // 跳过服务器证书验证，慎用
		},
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  true,
	}

	if publicCode.EnableProxy == true {

		var urlProxy string
		if publicCode.UserNameProxy != "" {
			urlProxy = strings.ToLower(publicCode.ModeProxy) + "://" + publicCode.UserNameProxy + ":" + publicCode.PassWordProxy + "@" + publicCode.HostProxy + ":" + publicCode.PortProxy
		} else {
			urlProxy = strings.ToLower(publicCode.ModeProxy) + "://" + publicCode.HostProxy + ":" + publicCode.PortProxy
		}

		proxyUrl, err := url.Parse(urlProxy)
		if err != nil {
			fmt.Printf("erroe: %v\n", err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	Client = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	ClientNoRedirect = &http.Client{
		Transport:     transport,
		Timeout:       timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }, // 禁止重定向，直接返回错误
	}

}
