package dns

import (
	"strings"
	"time"

	"github.com/storezhang/ddns/conf"
	"github.com/storezhang/ddns/core"
	"github.com/storezhang/ddns/lib"

	"github.com/goexl/dns"
	"github.com/goexl/gox/field"
	"github.com/pangum/schedule"
)

func (a *Agent) loadTask(config *conf.Config) (err error) {
	if err = a.config.Load(config); nil != err {
		return
	}

	// 先删除原来的任务
	a.scheduler.Clear()

	// 加载最新的任务
	for _, domain := range config.Resolves {
		for _, subdomain := range domain.Subdomains {
			_domain := core.NewDomain(
				domain.Name, subdomain, dns.TypeCname, domain.Value, domain.Ttl,
				domain.Prefix, domain.Staff,
			)

			secret := config.Secret(domain.Label)
			id := schedule.StringId(_domain.Final())
			switch {
			case domain.Contains(dns.TypeCname) && `` != strings.TrimSpace(domain.Value):
				executor := lib.NewCname(a.dns, secret, _domain, a.logger)
				_, err = a.scheduler.Add(executor, schedule.DurationTime(time.Second), id)
			case domain.Contains(dns.TypeA):
				executor := lib.NewA(a.dns, secret, _domain, a.wan, a.logger)
				_, err = a.scheduler.Add(executor, schedule.Duration(config.Ddns.Interval), id)
			default:
				a.logger.Error(`配置有误`, field.String(`domain`, _domain.Final()), field.Duration(`ttl`, domain.Ttl))
			}
		}
	}

	return
}
