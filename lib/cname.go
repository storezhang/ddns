package lib

import (
	"github.com/storezhang/ddns/conf"
	"github.com/storezhang/ddns/core"

	"github.com/pangum/dns"
	"github.com/pangum/logging"
)

// Cname CNAME解析
type Cname struct {
	dns    *dns.Client
	secret *conf.Secret
	domain *core.Domain
	logger *logging.Logger
}

// NewCname 创建CNAME解析
func NewCname(dns *dns.Client, secret *conf.Secret, domain *core.Domain, logger *logging.Logger) *Cname {
	return &Cname{
		dns:    dns,
		secret: secret,
		domain: domain,
		logger: logger,
	}
}
