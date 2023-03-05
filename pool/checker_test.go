package pool

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestChecker(t *testing.T) {

	raw, err := os.ReadFile("proxy.json")

	if err != nil {
		return
	}

	rows := make([]map[string]interface{}, 0)

	err = json.Unmarshal(raw, &rows)
	if err != nil {
		t.Error(err)
	}

	checker := &Checker{}

	for _, row := range rows {
		schema := row["schema"].(string)
		proxy := row["proxy"].(string)

		addrs := strings.ToLower(schema) + "://" + proxy

		report, err := checker.HealthCheck((*ProxyAddr)(&addrs))
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		fmt.Println("report:", report)
	}
}

func TestChecker2(t *testing.T) {
	addrs := []string{"socks5://127.0.0.1:10808", "http://127.0.0.1:10809"}
	checker := &Checker{}
	for _, v := range addrs {
		report, err := checker.HealthCheck((*ProxyAddr)(&v))
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("report:", report)
	}
}
