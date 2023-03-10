package pool

type Source interface {
	Reader() <-chan string
	Close()
}

type SrouceGroup struct {
	sources []Source
	que     chan string
	closer  chan bool
}

func NewSrouceGroup(ss ...Source) Source {
	ret := &SrouceGroup{
		sources: ss,
		closer:  make(chan bool, 1),
		que:     make(chan string, 100*len(ss)+1),
	}
	ret.init()
	return ret
}

func (s *SrouceGroup) init() {
	for _, source := range s.sources {
		go func(source Source) {
			rd := source.Reader()
			defer source.Close()
			for {
				select {
				case <-s.closer:
					return
				case v := <-rd:
					s.que <- v
				}
			}
		}(source)
	}
}

func (s *SrouceGroup) Close() {
	close(s.closer)
}

func (s *SrouceGroup) Reader() <-chan string {
	return s.que
}
