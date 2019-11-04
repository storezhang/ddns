package dns

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	log "github.com/sirupsen/logrus"
	"github.com/storezhang/gos/stringsx"
)

type Aliyun struct {
	AppKey string `yaml:"appKey"`
	Secret string `yaml:"secret"`
}

func (aliyun *Aliyun) Resolve(domain string, rr string, value string, dnsType string, ttl int) {
	client := getClient(aliyun.AppKey, aliyun.Secret)

	if recordId, add := getRecordId(client, domain, rr, value, dnsType, ttl); !add {
		if _, err := client.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
			RecordId: recordId,
			RR:       rr,
			Type:     dnsType,
			Value:    value,
			TTL:      requests.NewInteger(ttl),
		}); nil != err {
			log.WithFields(log.Fields{
				"domain": domain,
				"rr":     rr,
				"type":   dnsType,
				"value":  value,
				"ttl":    requests.NewInteger(ttl),
				"err":    err,
			}).Error("修改解析记录出错")
		}
	}
}

var recordIdCache map[string]string

func getRecordId(
	client *alidns.Client,
	domain string,
	rr string,
	dnsType string,
	value string,
	ttl int,
) (recordId string, add bool) {
	recordIdKey := stringsx.Contract("-", domain, rr, dnsType)
	if cacheRecordId, ok := recordIdCache[recordIdKey]; !ok {
		if prQueryRsp, err := client.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
			DomainName:  domain,
			RRKeyWord:   rr,
			TypeKeyWord: dnsType,
		}); nil != err {
			for _, record := range prQueryRsp.DomainRecords.Record {
				if domain == record.DomainName && dnsType == record.Type && rr == record.RR {
					recordId = record.RecordId
				}
			}
		} else {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("查询阿里域名解析出错")
		}

		if "" == recordId {
			if addRsp, err := client.AddDomainRecord(&alidns.AddDomainRecordRequest{
				DomainName: domain,
				RR:         rr,
				Type:       dnsType,
				Value:      value,
				TTL:        requests.NewInteger(ttl),
			}); nil == err {
				recordId = addRsp.RecordId
			} else {
				log.WithFields(log.Fields{
					"domain": domain,
					"rr":     rr,
					"type":   dnsType,
					"value":  value,
					"ttl":    requests.NewInteger(ttl),
					"err":    err,
				}).Error("添加解析记录出错")
			}
		}
		// 将recordId放入缓存
		recordIdCache[recordIdKey] = recordId
	} else {
		recordId = cacheRecordId
	}

	return
}

var clientCache map[string]*alidns.Client

func getClient(appKey string, secret string) (client *alidns.Client) {
	clientKey := stringsx.Contract("-", appKey, secret)
	if cacheClient, ok := clientCache[clientKey]; !ok {
		newClient, err := alidns.NewClientWithAccessKey("cn-hangzhou", appKey, secret)
		if nil != err {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("创建阿里云客户端出错")
		}
		client = newClient
		clientCache[clientKey] = newClient
	} else {
		client = cacheClient
	}

	return
}
