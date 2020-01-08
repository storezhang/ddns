package dns

import (
    "github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
    log "github.com/sirupsen/logrus"
    "github.com/storezhang/gos/stringsx"
)

// Aliyun 阿里云对象
type Aliyun struct {
    AppKey string `yaml:"appKey" toml:"appKey"`
    Secret string
}

// Resolve 阿里云域名解析
func (aliyun *Aliyun) Resolve(
    domain string,
    rr string,
    value string,
    dnsType string,
    ttl int,
) (result ResolveResult, err error) {
    client := getClient(aliyun.AppKey, aliyun.Secret)

    record, queryErr := queryRecord(client, domain, rr, dnsType)
    if nil != queryErr {
        err = queryErr
        log.WithFields(log.Fields{
            "domain": domain,
            "rr":     rr,
            "type":   dnsType,
            "value":  value,
            "ttl":    requests.NewInteger(ttl),
            "error":  queryErr,
        }).Error("查询解析记录出错")

        return
    }

    if nil == record { // 无记录，增加记录
        if addErr := add(client, domain, rr, value, dnsType, ttl); nil == addErr {
            result.Success = true
            result.Before = value
            result.After = value
        } else {
            err = addErr
            result.Success = false
        }
    } else if value != record.Value { // 有记录，和当前值不一致，更新
        if updateErr := update(client, record, domain, rr, value, dnsType, ttl); nil == updateErr {
            result.Success = true
            result.Before = record.Value
            result.After = value
        } else {
            err = updateErr
            result.Success = false
        }
    } else { // 不需要做更新
        result.Success = false
        result.Before = value
        result.After = value
        log.WithFields(log.Fields{
            "domain": domain,
            "rr":     rr,
            "type":   dnsType,
            "value":  value,
            "ttl":    ttl,
        }).Info("已有相同记录，未做修改")
    }

    return
}

func update(
    client *alidns.Client,
    record *Record,
    domain string,
    rr string,
    value string,
    dnsType string,
    ttl int,
) (err error) {
    req := alidns.CreateUpdateDomainRecordRequest()
    req.RecordId = record.ID
    req.RR = rr
    req.Type = dnsType
    req.Value = value
    req.TTL = requests.NewInteger(ttl)
    rsp, updateErr := client.UpdateDomainRecord(req)
    if nil != updateErr {
        err = updateErr
        log.WithFields(log.Fields{
            "domain": domain,
            "rr":     rr,
            "type":   dnsType,
            "value":  value,
            "ttl":    requests.NewInteger(ttl),
            "error":  updateErr,
        }).Error("修改解析记录出错")

        return
    }
    if nil != rsp && !rsp.IsSuccess() {
        log.WithFields(log.Fields{
            "domain": domain,
            "rr":     rr,
            "type":   dnsType,
            "value":  value,
            "ttl":    requests.NewInteger(ttl),
            "error":  updateErr,
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

    return
}

func add(
    client *alidns.Client,
    domain string,
    rr string,
    value string,
    dnsType string,
    ttl int,
) (err error) {
    req := alidns.CreateAddDomainRecordRequest()
    req.DomainName = domain
    req.RR = rr
    req.Type = dnsType
    req.Value = value
    req.TTL = requests.NewInteger(ttl)
    if addRsp, addErr := client.AddDomainRecord(req); nil == addErr {
        record := &Record{
            ID:     addRsp.RecordId,
            Domain: domain,
            Value:  value,
            RR:     rr,
        }
        // 将recordId放入缓存
        recordCache[key(domain, rr, dnsType)] = record
    } else {
        err = addErr
        log.WithFields(log.Fields{
            "domain": domain,
            "rr":     rr,
            "type":   dnsType,
            "value":  value,
            "ttl":    requests.NewInteger(ttl),
            "error":  addErr,
        }).Error("添加解析记录出错")
    }

    return
}

func queryRecord(
    client *alidns.Client,
    domain string,
    rr string,
    dnsType string,
) (record *Record, err error) {
    if nil == recordCache {
        recordCache = make(map[string]*Record)
    }

    recordIdKey := key(domain, rr, dnsType)
    if cacheRecord, ok := recordCache[recordIdKey]; !ok {
        req := alidns.CreateDescribeDomainRecordsRequest()
        req.DomainName = domain
        req.RRKeyWord = rr
        req.TypeKeyWord = dnsType
        if prQueryRsp, queryErr := client.DescribeDomainRecords(req); nil == queryErr {
            for _, serverRecord := range prQueryRsp.DomainRecords.Record {
                if domain == serverRecord.DomainName && dnsType == serverRecord.Type && rr == serverRecord.RR {
                    record = &Record{
                        ID:     serverRecord.RecordId,
                        Domain: domain,
                        Value:  serverRecord.Value,
                        RR:     serverRecord.RR,
                    }
                }
            }
        } else {
            err = queryErr
            return
        }
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

func key(domain string, rr string, dnsType string) string {
    return stringsx.Contract("-", domain, rr, dnsType)
}
