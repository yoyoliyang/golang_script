package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yoyoliyang/gomod/getpubip"
)

var accessKeyID = os.Getenv("ALIYUN_ACCESSKEYID")
var accessSecret = os.Getenv("ALIYUN_ACCESSSECRET")
var domainName = os.Getenv("DOMAINNAME")

// response struct用来作为json的解析数据
type response struct {
	RequesetId    string `json:"RequestId"`
	DomainRecords struct {
		Record []request
	} `json:"DomainRecords"`
}

// 一个request的接口，用来以后扩展方法
type requester interface {
	getRecord()
}

// request结构体，保存request的参数
type request struct {
	// 非json解析字段 发送请求的数据
	Client     *alidns.Client
	IP         net.IP // 公网获取到的最新ip保存到此处
	Scheme     string // 默认https
	Type       string // @
	DomainName string // 解析的域名,环境变量取得

	// json解析字段 // 解析response返回的数据
	RR       string `json:"RR"`
	RecordId string `json:"RecordId"`
	Value    string `json:"Value"` // 已经映射的ip地址
}

// 获取域名Record,返回结构体,提供RecordId, Value(IP), RR(@)
func (r *request) getRecord() (request request, err error) {
	bytes, err := r.describe()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
	resp := &response{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return request, err
	}
	records := resp.DomainRecords.Record
	for _, record := range records {
		// 只处理@根域名解析
		if record.RR == "@" {
			return record, nil
		}
	}
	// 若没有找到@解析，那么返回空结构体
	return request, nil
}

// 获取域名解析记录信息
func (r *request) describe() (b []byte, err error) {
	req := alidns.CreateDescribeDomainRecordsRequest()

	req.Scheme = r.Scheme
	req.DomainName = domainName

	resp, err := r.Client.DescribeDomainRecords(req)
	if err != nil {
		return b, err
	}
	b = resp.GetHttpContentBytes()
	return b, nil

}

// 添加domain记录
// https://help.aliyun.com/document_detail/29772.html?spm=a2c4g.11186623.6.668.cdac5eb4iEuHVt
func (r *request) add() (b []byte, err error) {

	req := alidns.CreateAddDomainRecordRequest()
	req.Scheme = r.Scheme

	req.DomainName = r.DomainName
	req.RR = r.RR
	req.Type = r.Type
	req.Value = r.IP.String()

	resp, err := r.Client.AddDomainRecord(req)
	if err != nil {
		return b, err
	}
	return resp.GetHttpContentBytes(), nil

}

// 更新domain记录
func (r *request) update() (b []byte, err error) {

	req := alidns.CreateUpdateDomainRecordRequest()
	req.Scheme = r.Scheme

	req.RR = r.RR
	req.RecordId = r.RecordId
	req.Type = r.Type
	req.Value = r.IP.String()

	resp, err := r.Client.UpdateDomainRecord(req)
	if err != nil {
		return b, err
	}

	return resp.GetHttpContentBytes(), nil

}

func main() {
	if accessKeyID != "" && accessSecret != "" && domainName != "" {

		// 获取公网ip，并与record中的value对比，如果不同，则进行更新，否则跳过
		ip, err := getpubip.GetIP()
		if err != nil {
			fmt.Fprintln(os.Stdout, "Failed to get public IP")
			return
		}

		// 解析客户端初始化
		client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessSecret)
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
		}

		// RecordId   string `json:"RecordId`
		// 初始化request struct数据, 因为RecordId还没有获取到，所以先不进行赋值
		var req = &request{
			Client:     client,
			Scheme:     "https",
			Type:       "A",
			RR:         "@",
			DomainName: domainName,
			IP:         ip,
		}

		// 获取@解析记录的Record
		record, err := req.getRecord()
		if err != nil {
			fmt.Fprintf(os.Stdout, "get record error %v\n", err)
		}

		req.RecordId = record.RecordId

		// 没有找到@解析记录，使用add方法方法添加
		if req.RecordId == "" {
			resp, err := req.add()
			if err != nil {
				fmt.Fprintf(os.Stdout, "add domain record error %v\n", err)
				return
			}

			fmt.Fprintln(os.Stdout, string(resp))
			fmt.Fprintf(os.Stdout, "add domain record success! %v -> %v", req.DomainName, req.IP.String())
			return
		}

		// 存在相同的ip，不去执行操作
		if req.IP.String() == record.Value {
			fmt.Fprintf(os.Stdout, "nochange %v %v\n", req.DomainName, ip.String())
			return
		}

		// 已经发生了ip变化，要进行解析操作处理
		updateResult, err := req.update()
		if err != nil {
			fmt.Fprintf(os.Stdout, "update domain record error %v\n", err)
			return
		}
		// 显示更新后的结果
		fmt.Fprintln(os.Stdout, string(updateResult))
		fmt.Fprintf(os.Stdout, "ip has changed! current : %v\n", ip.String())

	} else {
		// 没有环境变量的处理
		fmt.Fprintln(os.Stdout, "ERR: Problems with Environment Variables")
		return
	}
}
