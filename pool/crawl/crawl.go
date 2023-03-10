package crawl

import (
	"fmt"
	"log"
	"time"
)

type Crawler struct {
	que chan string

	tk *time.Ticker
}

func NewCrawler() *Crawler {
	ret := &Crawler{que: make(chan string, 100)}
	ret.tk = time.NewTicker(time.Minute * 20)
	go func() {
		for {
			ret.work()
			<-ret.tk.C
		}
	}()

	return ret
}

func (c *Crawler) Reader() <-chan string {
	return c.que
}

func (c *Crawler) Close() {
	c.tk.Stop()
	close(c.que)
}

func (c *Crawler) work() {
	wrap := func(name string, f func() []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("crawl get error:", name, err)
			}
		}()
		addrs := f()
		if len(addrs) == 0 {
			fmt.Println("crawl get empty:", name)
		}
		for _, addr := range addrs {
			c.que <- addr
		}
	}

	for name, task := range Task {
		wrap(name, task)
	}
}

var Task = make(map[string]func() []string)

func RegTask(name string, task func() []string) {
	Task[name] = task
}
