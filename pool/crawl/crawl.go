package crawl

import (
	"time"
)

type Crawler interface {
	Name() string
	Work() (*CrawlResult, error)
}

type CrawlResult struct {
	CrawlerName string
	Addrs       []string
	CostTime    time.Duration
	Err         error
}

var AllCrawlers = []Crawler{
	&HensonGetter{},
}

func RunAllCrawlers() (ret []*CrawlResult) {
	for _, crawler := range AllCrawlers {
		result, err := crawler.Work()
		if err != nil {
			result.Err = err
		}
		ret = append(ret, result)
	}
	return
}
