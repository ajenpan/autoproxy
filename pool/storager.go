package pool

type Storage interface {
	Save(*ProxyAddr) error
	GetAll() ([]*ProxyAddr, error)
	Get(string) (*ProxyAddr, error)
	Delete(string)
	Count() int
	RandOne() (*ProxyAddr, error)
	Flash() error
}
