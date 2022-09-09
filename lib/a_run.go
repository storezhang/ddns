package lib

import (
	"context"

	"github.com/storezhang/ddns/conf"

	"github.com/goexl/dns"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

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
		a.logger.Info(`地址未改变，不做更新处理`, fields.Connect(field.String(`original`, a.domain.Value()))...)
	}
	if ip == a.domain.Value() {
		return
	}

	options := dns.NewOptions(dns.Secret(a.secret.Ak, a.secret.Sk), dns.Ttl(a.domain.Ttl()), dns.A())
	switch a.secret.Type {
	case conf.TypeAliyun:
		options.Add(dns.Aliyun())
	case conf.TypeTencentCloud, conf.TypeDnspod:
		options.Add(dns.TencentCloud())
	}

	if original, do, de := a.dns.Resolve(
		context.Background(),
		a.domain.Name(), a.domain.Subdomain(), ip,
		options...,
	); nil != de {
		err = de
	} else if do { // 如果设置成功，则更新原值，防止下次再请求数据
		a.domain.SetValue(ip)
		a.logger.Info(`域名解析成功`, fields.Connect(field.String(`original`, original))...)
	} else if ip == original { // 记录值和原来一致，也要更新原来的记录值，防止后续继续调用解析方法
		a.domain.SetValue(ip)
		a.logger.Info(`记录值和原来一致，解析记录未更新`, fields.Connect(field.String(`original`, original))...)
	} else {
		a.logger.Warn(`遇到未知原因，解析记录未更新`, fields.Connect(field.String(`original`, original))...)
	}

	return
}
