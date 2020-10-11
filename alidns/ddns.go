package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yoyoliyang/gomod/getpubip"
)

var accessKeyID = os.Getenv("ALIYUN_ACCESSKEYID")
var accessSecret = os.Getenv("ALIYUN_ACCESSSECRET")
var domainName = os.Getenv("DOMAINNAME")
var domainid = os.Getenv("DOMAINID")

// 更新domain记录
func updateDC(ip net.IP) {

	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessSecret)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.Value = ip.String()
	request.Type = "A"
	request.RR = "@"
	request.RecordId = domainid

	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error for UpdateDomainRecord: %v", err)
	}
	fmt.Fprintf(os.Stdout, "response is %#v\n", response)
}

func main() {
	if accessKeyID != "" && accessSecret != "" && domainName != "" {

		filename := "ip.txt"

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stdout, "Failed to open file: \n", file.Name())
			return
		}
		defer file.Close()

		oldIP, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error for reading ip file %v\n", err)
			return
		}

		ip, err := getpubip.GetIP()
		if err != nil {
			fmt.Fprintln(os.Stdout, "Failed to get public IP")
			return
		}

		// 存在相同的ip，不去执行操作
		if ip.String() == string(oldIP) {
			fmt.Fprintf(os.Stdout, "nochange %v", ip.String())
			return
		}

		fmt.Fprintf(os.Stdout, "current ip: %v\n", ip.String())

		// Truncate用法 https://golang.org/pkg/os/#File.Truncate
		if err = file.Truncate(0); err != nil {
			fmt.Fprintf(os.Stdout, "Failed to clear ip file %v", err)
			return
		}
		_, err = file.WriteString(ip.String())
		if err != nil {
			fmt.Fprintf(os.Stdout, "Failed to write IP address, %v\n", err)
			return
		}

		// 解析操作处理
		updateDC(ip)

	} else {
		fmt.Fprintln(os.Stdout, "ERR: Problems with Environment Variables")
		return
	}
}
