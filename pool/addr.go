package pool

import (
	"net/url"
)

type ProxyAddr string

func (p *ProxyAddr) Url() *url.URL {
	u, err := url.Parse(string(*p))
	if err != nil {
		//TODO:
		panic(err)
	}
	return u
}

type ProxyAddrItem struct {
	Scheme string
	Host   string
	Port   string
}

func SplitProxyAddr(raw string) (*ProxyAddrItem, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	return &ProxyAddrItem{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Port:   u.Port(),
	}, nil
}

func (p *ProxyAddrItem) String() string {
	return p.Scheme + "://" + p.Host + ":" + p.Port
}
