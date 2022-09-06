package lib

import (
	"context"

	"ddns/conf"

	"github.com/goexl/dns"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (c *Cname) Run() (err error) {
	options := dns.NewOptions(dns.Secret(c.secret.Ak, c.secret.Sk), dns.Ttl(c.domain.Ttl()), dns.CNAME())
	switch c.secret.Type {
	case conf.TypeAliyun:
		options.Add(dns.Aliyun())
	case conf.TypeTencentCloud, conf.TypeDnspod:
		options.Add(dns.TencentCloud())
	}

	fields := gox.Fields{
		field.String(`domain`, c.domain.Final()),
		field.Duration(`ttl`, c.domain.Ttl()),
	}
	if original, do, de := c.dns.Resolve(
		context.Background(),
		c.domain.Name(), c.domain.Subdomain(), c.domain.Value(),
		options...,
	); nil != de {
		err = de
	} else if do {
		c.logger.Info(`域名解析成功`, fields.Connect(field.String(`original`, original))...)
	} else {
		c.logger.Info(`未做域名解析`, fields.Connect(field.String(`original`, original))...)
	}

	return
}
