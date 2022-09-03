package lib

import (
	"context"

	"ddns/conf"
	"ddns/core"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/uda"
	"github.com/pangum/logging"
	"github.com/pangum/wanip"
)

// A A地址解析
type A struct {
	uda    *uda.Client
	secret *conf.Secret
	domain *core.Domain
	wan    *wanip.Agent
	logger *logging.Logger
}

// NewA 创建A解析
func NewA(secret *conf.Secret, domain *core.Domain, wan *wanip.Agent, logger *logging.Logger) *A {
	return &A{
		uda:    uda.New(),
		secret: secret,
		domain: domain,
		wan:    wan,
		logger: logger,
	}
}

func (a *A) Run() (err error) {
	var ip string
	if ip, err = a.wan.Get(); nil != err {
		return
	}

	fields := gox.Fields{
		field.String(`domain`, a.domain.Final()),
		field.String(`ip`, ip),
		field.Duration(`ttl`, a.domain.Ttl()),
	}
	// 判断外网地址是否改变
	if ip == a.domain.Value() {
		a.logger.Info(`未做域名解析`, fields.Connect(field.String(`original`, a.domain.Value()))...)
	}

	options := uda.NewOptions(uda.Secret(a.secret.Ak, a.secret.Sk), uda.Ttl(a.domain.Ttl()), uda.A())
	switch a.secret.Type {
	case conf.TypeAliyun:
		options.Add(uda.Aliyun())
	}

	if original, do, udaErr := a.uda.Resolve(
		context.Background(),
		a.domain.Name(), a.domain.Subdomain(), ip,
		options...,
	); nil != udaErr {
		err = udaErr
	} else if do {
		a.logger.Info(`域名解析成功`, fields.Connect(field.String(`original`, original))...)
	} else {
		a.logger.Info(`未做域名解析`, fields.Connect(field.String(`original`, original))...)
	}

	return
}
