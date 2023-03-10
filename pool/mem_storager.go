package pool

import (
	"fmt"
	"sync"
)

func NewMemStorager() *MemStorager {
	return &MemStorager{
		addrs: map[string]*ProxyAddr{},
	}
}

type MemStorager struct {
	locker sync.RWMutex
	addrs  map[string]*ProxyAddr
}

func (s *MemStorager) Save(addr *ProxyAddr) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.addrs[addr.Addr] = addr
	return nil
}
func (s *MemStorager) GetAll() ([]*ProxyAddr, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	ret := []*ProxyAddr{}
	for _, v := range s.addrs {
		ret = append(ret, v)
	}
	return ret, nil
}

func (s *MemStorager) Get(addr string) (*ProxyAddr, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	v, ok := s.addrs[addr]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return v, nil
}

func (s *MemStorager) Delete(addr string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.addrs, addr)
}

func (s *MemStorager) Count() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.addrs)
}

func (s *MemStorager) RandOne() (*ProxyAddr, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	for _, v := range s.addrs {
		return v, nil
	}
	return nil, fmt.Errorf("no proxy addr")
}

func (s *MemStorager) Flash() error {
	return nil
}
