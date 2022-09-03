package lib

import (
	"context"

	"ddns/conf"
	"ddns/core"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/uda"
	"github.com/pangum/logging"
)

// Cname CNAME解析
type Cname struct {
	uda    uda.Resolver
	secret *conf.Secret
	domain *core.Domain
	logger *logging.Logger
}

// NewCname 创建CNAME解析
func NewCname(secret *conf.Secret, domain *core.Domain, logger *logging.Logger) *Cname {
	return &Cname{
		uda:    uda.New(),
		secret: secret,
		domain: domain,
		logger: logger,
	}
}

func (c *Cname) Run() (err error) {
	options := uda.NewOptions(uda.Secret(c.secret.Ak, c.secret.Sk), uda.Ttl(c.domain.Ttl()), uda.CNAME())
	switch c.secret.Type {
	case conf.TypeAliyun:
		options.Add(uda.Aliyun())
	}

	fields := gox.Fields{
		field.String(`domain`, c.domain.Final()),
		field.Duration(`ttl`, c.domain.Ttl()),
	}
	if original, do, udaErr := c.uda.Resolve(
		context.Background(),
		c.domain.Name(), c.domain.Subdomain(), c.domain.Value(),
		options...,
	); nil != udaErr {
		err = udaErr
	} else if do {
		c.logger.Info(`域名解析成功`, fields.Connect(field.String(`original`, original))...)
	} else {
		c.logger.Info(`未做域名解析`, fields.Connect(field.String(`original`, original))...)
	}

	return
}
