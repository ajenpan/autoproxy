package crawl

import (
	"fmt"
	"strings"
	"time"

	"github.com/henson/proxypool/getter"
	"github.com/henson/proxypool/pkg/models"
)

type HensonGetter struct {
}

type HensonGetterFunc func() []*models.IP

func (g *HensonGetter) Name() string {
	return "github.com/henson/proxypool-getter"
}

func (g *HensonGetter) Work() (*CrawlResult, error) {
	result := &CrawlResult{
		CrawlerName: g.Name(),
	}

	startAt := time.Now()
	defer func() {
		result.CostTime = time.Since(startAt)
	}()

	wraper := func(f HensonGetterFunc) (ret []*models.IP) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("panic:", err)
			}
		}()
		return f()
	}

	funs := []func() []*models.IP{
		getter.FQDL,  //新代理
		getter.PZZQZ, //新代理
		//getter.Data5u,
		//getter.Feiyi,
		//getter.IP66, //need to remove it
		getter.IP3306,
		getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		//getter.Xici,
		// getter.XDL,
		//getter.IP181,  // 已经无法使用
		//getter.YDL,	//失效的采集脚本，用作系统容错实验
		// getter.PLP, //need to remove it
		getter.PLPSSL,
		getter.IP89,
	}

	for _, fun := range funs {
		ips := wraper(fun)
		for _, item := range ips {
			addr := strings.ToLower(item.Type1) + "://" + item.Data
			result.Addrs = append(result.Addrs, addr)
		}
	}

	return result, nil
}
