package lib

import (
	"github.com/storezhang/ddns/conf"
	"github.com/storezhang/ddns/core"

	"github.com/pangum/dns"
	"github.com/pangum/logging"
	"github.com/pangum/wanip"
)

// A A地址解析
type A struct {
	dns    *dns.Client
	secret *conf.Secret
	domain *core.Domain
	wan    *wanip.Agent
	logger *logging.Logger
}

// NewA 创建A解析
func NewA(dns *dns.Client, secret *conf.Secret, domain *core.Domain, wan *wanip.Agent, logger *logging.Logger) *A {
	return &A{
		dns:    dns,
		secret: secret,
		domain: domain,
		wan:    wan,
		logger: logger,
	}
}
