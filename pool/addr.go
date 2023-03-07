package pool

import (
	"fmt"
	"net/http"
	"net/url"
)

// type ProxyAddr string

type ProxyAddr struct {
	Header http.Header

	address string
	url     *url.URL
}

func Addrs2ProxyAddrs(s []string) []*ProxyAddr {
	addrs := make([]*ProxyAddr, 0)
	for _, v := range s {
		addr, err := NewProxyAddrFromStr(v)
		if err != nil {
			fmt.Println(err)
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func NewProxyAddr(u *url.URL) *ProxyAddr {
	ret := &ProxyAddr{
		url:     u,
		address: u.String(),
		Header:  http.Header{},
	}
	return ret
}

func NewProxyAddrFromStr(addr string) (*ProxyAddr, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return &ProxyAddr{
		address: u.String(),
		url:     u,
		Header:  http.Header{},
	}, nil
}

func (p *ProxyAddr) URL() *url.URL {
	return p.url
}

func (p *ProxyAddr) Address() string {
	return p.address
}
