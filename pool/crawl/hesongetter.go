package crawl

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func init() {
	RegTask("PZZQZ", PZZQZ)
	RegTask("IP3306", IP3306)
	RegTask("KDL", KDL)
	RegTask("PLPSSL", PLPSSL)
	RegTask("89ip", IP89)
}

func newRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	return req
}

func loadDoc(url string) (doc *html.Node, err error) {
	resp, err := http.DefaultClient.Do(newRequest(url))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return htmlquery.Parse(resp.Body)
}

func PZZQZ() (result []string) {
	pollURL := "http://pzzqz.com/"
	doc, err := loadDoc(pollURL)
	if err != nil {
		return
	}
	trNode := htmlquery.Find(doc, "//table[@class='table table-hover']//tbody//tr")
	for i := 0; i < len(trNode); i++ {
		tdNode := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		scheme := htmlquery.InnerText(tdNode[4])
		result = append(result, fmt.Sprintf("%s://%s:%s", scheme, ip, port))
	}
	return
}

func IP3306() (result []string) {
	pollURL := "http://www.ip3366.net/free/?stype=1&page=1"
	doc, err := loadDoc(pollURL)
	if err != nil {
		return
	}
	trNode := htmlquery.Find(doc, "//div[@id='list']//table//tbody//tr")

	for i := 1; i < len(trNode); i++ {
		tdNode := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		result = append(result, fmt.Sprintf("%s://%s:%s", strings.ToLower(Type), ip, port))
	}

	return
}

func KDL() (result []string) {
	pollURL := "http://www.kuaidaili.com/free/inha/"
	doc, err := loadDoc(pollURL)
	if err != nil {
		return
	}
	trNode := htmlquery.Find(doc, "//table[@class='table table-bordered table-striped']//tbody//tr")
	for i := 0; i < len(trNode); i++ {
		tdNode := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		result = append(result, fmt.Sprintf("%s://%s:%s", strings.ToLower(Type), ip, port))
	}
	return
}

func PLPSSL() (result []string) {
	pollURL := "https://list.proxylistplus.com/SSL-List-1"
	doc, err := loadDoc(pollURL)
	if err != nil {
		return
	}
	trNode := htmlquery.Find(doc, "//div[@class='hfeed site']//table[@class='bg']//tbody//tr")

	for i := 3; i < len(trNode); i++ {
		tdNode := htmlquery.Find(trNode[i], "//td")
		if len(tdNode) < 6 {
			return
		}
		ip := htmlquery.InnerText(tdNode[1])
		port := htmlquery.InnerText(tdNode[2])
		Type := htmlquery.InnerText(tdNode[6])
		if Type == "yes" {
			Type = "https"
		} else {
			Type = "http"
		}
		result = append(result, fmt.Sprintf("%s://%s:%s", Type, ip, port))
	}
	return
}

func IP89() (result []string) {
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	pollURL := "http://www.89ip.cn/tqdl.html?api=1&num=100&port=&address=%E7%BE%8E%E5%9B%BD&isp="

	resp, err := http.DefaultClient.Do(newRequest(pollURL))
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	bodyIPs := string(body)
	ips := ExprIP.FindAllString(bodyIPs, 100)

	for index := 0; index < len(ips); index++ {
		data := strings.TrimSpace(ips[index])
		result = append(result, fmt.Sprintf("http://%s", data))
	}
	return
}
