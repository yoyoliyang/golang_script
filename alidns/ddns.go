package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yoyoliyang/gomod/getpubip"
)

var accessKeyID = os.Getenv("ALIYUN_ACCESSKEYID")
var accessSecret = os.Getenv("ALIYUN_ACCESSSECRET")
var domainName = os.Getenv("DOMAINNAME")

func main() {
	if accessKeyID != "" && accessSecret != "" && domainName != "" {

		filename := "ip.txt"
		file, err := os.Open(filename)
		if err != nil {
			file, _ = os.Create(filename)
		}
		defer file.Close()

		oldIP, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error for reading ip file %v", err)
			return
		}

		ip, _ := getpubip.GetIP()

		if ip.String() == string(oldIP) {
			fmt.Fprintf(os.Stdout, "nochange %v", ip.String())
			return
		}

		fmt.Fprintf(os.Stdout, "current ip: %v\n", ip.String())
		defer ioutil.WriteFile(filename, []byte(ip.String()), 0644)

		client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessSecret)

		request := alidns.CreateAddDomainRecordRequest()
		request.Scheme = "https"

		request.Value = ip.String()
		request.Type = "A"
		request.RR = "@"
		request.DomainName = domainName

		response, err := client.AddDomainRecord(request)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error for AddDomainRecord: %v", err)
		}
		fmt.Fprintf(os.Stdout, "response is %#v\n", response)

	} else {
		fmt.Fprintln(os.Stdout, "ERR: Problems with Environment Variables")
		return
	}
}
