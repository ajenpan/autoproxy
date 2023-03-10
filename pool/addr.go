package pool

import (
	"fmt"
	"net/url"
	"time"
)

type ProxyAddr struct {
	Addr    string                 `json:"addr"`
	Speed   string                 `json:"speed"`
	Health  int8                   `json:"health"`
	CheckAt time.Time              `json:"check_at"`
	Extra   map[string]interface{} `json:"extra"`

	url *url.URL
}

func Addrs2ProxyAddrs(s []string) []*ProxyAddr {
	addrs := make([]*ProxyAddr, 0)
	for _, v := range s {
		addr, err := NewProxyAddrWithStr(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func NewProxyAddr(u *url.URL) *ProxyAddr {
	ret := &ProxyAddr{
		url:   u,
		Addr:  u.String(),
		Extra: make(map[string]interface{}),
	}
	return ret
}

func NewProxyAddrWithStr(addr string) (*ProxyAddr, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return &ProxyAddr{
		Addr:  u.String(),
		url:   u,
		Extra: make(map[string]interface{}),
	}, nil
}

func (p *ProxyAddr) URL() *url.URL {
	return p.url
}

func (p *ProxyAddr) Address() string {
	return p.Addr
}

func (p *ProxyAddr) Clone() *ProxyAddr {
	n, _ := NewProxyAddrWithStr(p.Addr)
	return n
}
