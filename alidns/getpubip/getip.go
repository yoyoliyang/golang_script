package getpubip

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
)

// GetIP 获取公网ip模块
func GetIP() (ip net.IP, err error) {
	var url string
	flag.StringVar(&url, "m", "", "ip138.com mirror")
	flag.Parse()
	if url == "" {
		url = "http://2021.ip138.com"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stdout, "error for new request")
		return ip, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36")
	req.Header.Set("Accept-Charset", "gb2312")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stdout, "error for request")
		return ip, err
	}

	scanner := bufio.NewScanner(resp.Body)
	ipf := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	var adr string
	for scanner.Scan() {
		ipArr := ipf.FindAllString(scanner.Text(), -1)
		if len(ipArr) == 1 {
			adr = ipArr[0]
		}
	}
	if adr == "" {
		fmt.Fprintln(os.Stdout, "not found IP in ip138.com")
		return ip, err
	}

	return net.ParseIP(adr), nil

}
