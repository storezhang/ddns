package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/storezhang/gos/nets"

	"ddns/common"
	"ddns/sign"
)

// DDNSJob 动态域名解析任务
type DDNSJob struct {
	resolver sign.Signer
	domain   common.Domain
}

// Run 动态域名解析任务真正执行的方法
func (ddnsJob *DDNSJob) Run() {
	refreshDns(ddnsJob.resolver, ddnsJob.domain)
}

func refreshDns(resolver sign.Signer, domain common.Domain) {
	for _, dnsType := range strings.Split(domain.DNSTypes, ",") {
		value := getValue(dnsType, domain)
		for _, subDomain := range strings.Split(domain.SubDomains, ",") {
			if "" != strings.TrimSpace(domain.SubDomainPrefix) {
				subDomain = fmt.Sprintf("%s.%s", domain.SubDomainPrefix, subDomain)
			}
			if "" != strings.TrimSpace(domain.SubDomainStaff) {
				subDomain = fmt.Sprintf("%s.%s", subDomain, domain.SubDomainStaff)
			}

			log.WithFields(log.Fields{
				"name":      domain.Name,
				"subDomain": subDomain,
				"dnsType":   dnsType,
				"type":      domain.Type,
				"value":     value,
			}).Info("执行DNS解析更新")

			resolver.Resolve(
				domain.Name,
				subDomain,
				value,
				dnsType,
				domain.TTL,
			)
		}
	}
}

func getValue(dnsType string, domain common.Domain) string {
	var value string

	switch dnsType {
	case "A":
		if ip, err := nets.GetPublicIp(); nil != err {
			log.WithFields(log.Fields{
				"name":      domain.Name,
				"subDomain": domain.SubDomains,
				"dnsTypes":  domain.DNSTypes,
				"err":       err,
			}).Error("解析本机IP地址出错")
		} else {
			value = ip
		}
	default:
		value = domain.Value
	}

	return value
}
