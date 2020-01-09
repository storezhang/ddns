package main

import (
    `crypto/tls`
    `fmt`
    `strings`
    `time`

    `github.com/parnurzeal/gorequest`
    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
    `github.com/storezhang/gos/nets`
    `github.com/storezhang/gos/tpls`

    `ddns/common`
    `ddns/dns`
)

var crontab *cron.Cron
var req *gorequest.SuperAgent

func init() {
    // 初始化Http客户端
    req = gorequest.New()
    req.Timeout(60 * time.Second)
    // 忽略TLS证书
    req.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

// DDNSJob 动态域名解析任务
type DDNSJob struct {
    resolver dns.Resolver
    domain   common.Domain
    ddns     *common.DDNS
}

// Run 动态域名解析任务真正执行的方法
func (job *DDNSJob) Run() {
    refreshDns(job.resolver, job.ddns, job.domain)
}

func refreshDns(resolver dns.Resolver, ddns *common.DDNS, domain common.Domain) {
    for _, dnsType := range domain.DNSTypes {
        value := getValue(dnsType, domain)

        for _, subDomain := range domain.SubDomains {
            rr := subDomain
            if "" != strings.TrimSpace(domain.SubDomainPrefix) {
                rr = fmt.Sprintf("%s.%s", domain.SubDomainPrefix, subDomain)
            }
            if "" != strings.TrimSpace(domain.SubDomainStaff) {
                rr = fmt.Sprintf("%s.%s", subDomain, domain.SubDomainStaff)
            }

            // 对每个解析开启一个协程，增加性能
            go resolve(resolver, domain, ddns, rr, value, dnsType)
        }
    }
}

func resolve(resolver dns.Resolver, domain common.Domain, ddns *common.DDNS, rr string, value string, dnsType string) {
    if result, err := resolver.Resolve(
        domain.Name,
        rr,
        value,
        dnsType,
        domain.TTL,
    ); nil != err {
        log.WithFields(log.Fields{
            "name":    domain.Name,
            "rr":      rr,
            "dnsType": dnsType,
            "type":    domain.Type,
            "value":   value,
            "error":   err,
        }).Info("执行DNS解析更新失败")
    } else {
        log.WithFields(log.Fields{
            "name":    domain.Name,
            "rr":      rr,
            "dnsType": dnsType,
            "type":    domain.Type,
            "value":   value,
        }).Info("执行DNS解析更新成功")

        // 成功解析，发推送
        if result.Success {
            notify(domain, rr, ddns, result)
        }
    }
}

func getValue(dnsType string, domain common.Domain) string {
    var value string

    switch dnsType {
    case "A":
        if ip, err := nets.GetPublicIp(); nil != err {
            log.WithFields(log.Fields{
                "name":       domain.Name,
                "subDomains": domain.SubDomains,
                "dnsTypes":   domain.DNSTypes,
                "err":        err,
            }).Error("解析本机IP地址出错")
        } else {
            value = ip
        }
    default:
        value = domain.Value
    }

    return value
}

func notify(domain common.Domain, subDomain string, ddns *common.DDNS, result dns.ResolveResult) {
    var serverChans []common.ServerChan
    if nil != domain.Chans && 0 < len(domain.Chans) {
        serverChans = domain.Chans
    } else if nil != ddns.Chans && 0 < len(ddns.Chans) {
        serverChans = ddns.Chans
    } else {
        return
    }

    var titleTemplate string
    var contentTemplate string
    if "" != strings.TrimSpace(domain.Template.Title) && "" != strings.TrimSpace(domain.Template.Content) {
        titleTemplate = domain.Template.Title
        contentTemplate = domain.Template.Content
    } else if "" != strings.TrimSpace(ddns.Template.Title) && "" != strings.TrimSpace(ddns.Template.Content) {
        titleTemplate = ddns.Template.Title
        contentTemplate = ddns.Template.Content
    } else {
        return
    }

    notifyToUser(serverChans, titleTemplate, contentTemplate, domain, subDomain, result)
}

type notifyData struct {
    Domain    common.Domain
    SubDomain string
    Result    dns.ResolveResult
}

func notifyToUser(
    chans []common.ServerChan,
    titleTemplate string,
    contentTemplate string,
    domain common.Domain,
    subDomain string,
    result dns.ResolveResult,
) {
    data := &notifyData{
        Domain:    domain,
        SubDomain: subDomain,
        Result:    result,
    }
    title := tpls.Render("title", titleTemplate, data)
    desp := tpls.Render("desp", contentTemplate, data)

    // 真正发推送
    for _, ch := range chans {
        rsp, body, err := req.Post(fmt.Sprintf("https://sc.ftqq.com/%s.send", ch.Key)).
            Type("form").
            Send(common.ServerChanRequest{
                Text: title,
                Desp: desp,
            }).End()

        if nil != err {
            log.WithFields(log.Fields{
                "name": domain.Name,
                "chan": ch.Key,
                "rsp":  rsp,
                "body": body,
                "err":  err,
            }).Info("ServerChan推送消息失败")
        } else {
            log.WithFields(log.Fields{
                "name": domain.Name,
                "chan": ch.Key,
            }).Info("ServerChan推送消息成功")
        }
    }
}
