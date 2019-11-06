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

    if record, add, err := getRecordId(client, domain, rr, value, dnsType, ttl); nil == err {
        if add {
            return
        }

        if value == record.Value {
            log.WithFields(log.Fields{
                "domain": domain,
                "rr":     rr,
                "type":   dnsType,
                "value":  value,
                "ttl":    ttl,
            }).Info("已有相同记录，未做修改")
            return
        }

        req := alidns.CreateUpdateDomainRecordRequest()
        req.RecordId = record.RecordId
        req.RR = rr
        req.Type = dnsType
        req.Value = value
        req.TTL = requests.NewInteger(ttl)
        if rsp, err := client.UpdateDomainRecord(req); nil != rsp && !rsp.IsSuccess() {
            log.WithFields(log.Fields{
                "domain": domain,
                "rr":     rr,
                "type":   dnsType,
                "value":  value,
                "ttl":    requests.NewInteger(ttl),
                "err":    err,
            }).Error("修改解析记录出错")
        } else {
            log.WithFields(log.Fields{
                "domain": domain,
                "rr":     rr,
                "type":   dnsType,
                "value":  value,
                "ttl":    requests.NewInteger(ttl),
            }).Trace("修改解析记录成功")
        }
    } else {
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

var recordCache map[string]*alidns.Record

func getRecordId(
    client *alidns.Client,
    domain string,
    rr string,
    value string,
    dnsType string,
    ttl int,
) (record *alidns.Record, add bool, err error) {
    if nil == recordCache {
        recordCache = make(map[string]*alidns.Record)
    }

    recordIdKey := stringsx.Contract("-", domain, rr, dnsType)
    if cacheRecord, ok := recordCache[recordIdKey]; !ok {
        req := alidns.CreateDescribeDomainRecordsRequest()
        req.DomainName = domain
        req.RRKeyWord = rr
        req.TypeKeyWord = dnsType
        if prQueryRsp, queryErr := client.DescribeDomainRecords(req); nil == queryErr {
            for _, serverRecord := range prQueryRsp.DomainRecords.Record {
                if domain == serverRecord.DomainName && dnsType == serverRecord.Type && rr == serverRecord.RR {
                    record = &serverRecord
                }
            }
        } else {
            err = queryErr
            return
        }

        if nil == record {
            req := alidns.CreateAddDomainRecordRequest()
            req.DomainName = domain
            req.RR = rr
            req.Type = dnsType
            req.Value = value
            req.TTL = requests.NewInteger(ttl)
            if addRsp, err := client.AddDomainRecord(req); nil == err {
                record = &alidns.Record{
                    Value:      value,
                    TTL:        int64(ttl),
                    Remark:     "",
                    DomainName: domain,
                    RR:         rr,
                    Priority:   0,
                    RecordId:   addRsp.RecordId,
                    Status:     "",
                    Locked:     false,
                    Weight:     0,
                    Line:       "",
                    Type:       dnsType,
                }
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
        recordCache[recordIdKey] = record
    } else {
        record = cacheRecord
    }

    return
}

var clientCache map[string]*alidns.Client

func getClient(appKey string, secret string) (client *alidns.Client) {
    if nil == clientCache {
        clientCache = make(map[string]*alidns.Client)
    }

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
