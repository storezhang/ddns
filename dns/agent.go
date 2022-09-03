package dns

import (
	"strings"
	"time"

	"ddns/conf"
	"ddns/core"
	"ddns/lib"

	"github.com/goexl/gox/field"
	"github.com/goexl/uda"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"github.com/pangum/schedule"
	"github.com/pangum/wanip"
)

type (
	// Agent 执行器
	Agent struct {
		config    *pangu.Config
		scheduler *schedule.Scheduler
		wan       *wanip.Agent
		logger    *logging.Logger
	}

	agentIn struct {
		pangu.In

		Config    *pangu.Config
		Scheduler *schedule.Scheduler
		Wan       *wanip.Agent
		Logger    *logging.Logger
	}
)

func newAgent(in agentIn) *Agent {
	return &Agent{
		config:    in.Config,
		scheduler: in.Scheduler,
		wan:       in.Wan,
		logger:    in.Logger,
	}
}

func (a *Agent) Start() (err error) {
	if err = a.Run(); nil != err {
		return
	}
	a.scheduler.Start()

	return
}

func (a *Agent) Stop() (err error) {
	a.scheduler.Stop()

	return
}

func (a *Agent) Name() string {
	return `域名解析`
}

func (a *Agent) Run() (err error) {
	config := new(conf.Config)
	if err = a.config.Load(config); nil != err {
		return
	}

	for _, domain := range config.Resolves {
		for _, subdomain := range domain.Subdomains {
			_domain := core.NewDomain(
				domain.Name, subdomain, uda.TypeCname, domain.Value, domain.Ttl,
				domain.Prefix, domain.Staff,
			)

			secret := config.Secret(domain.Label)
			switch {
			case domain.Contains(uda.TypeCname) && `` != strings.TrimSpace(domain.Value):
				executor := lib.NewCname(secret, _domain, a.logger)
				_, err = a.scheduler.Add(executor, schedule.DurationTime(time.Second))
			case domain.Contains(uda.TypeA):
				executor := lib.NewA(secret, _domain, a.wan, a.logger)
				_, err = a.scheduler.Add(executor, schedule.Duration(5*time.Second))
			default:
				a.logger.Error(`配置有误`, field.String(`domain`, _domain.Final()), field.Duration(`ttl`, domain.Ttl))
			}
		}
	}

	return
}
