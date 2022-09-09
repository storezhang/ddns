package dns

import (
	"strings"

	"ddns/conf"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pangum/dns"
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
		dns       *dns.Client
		wan       *wanip.Agent
		logger    *logging.Logger
	}

	agentIn struct {
		pangu.In

		Config    *pangu.Config
		Scheduler *schedule.Scheduler
		Dns       *dns.Client
		Wan       *wanip.Agent
		Logger    *logging.Logger
	}
)

func newAgent(in agentIn) *Agent {
	return &Agent{
		config:    in.Config,
		scheduler: in.Scheduler,
		dns:       in.Dns,
		wan:       in.Wan,
		logger:    in.Logger,
	}
}

func (a *Agent) Start() (err error) {
	config := new(conf.Config)
	// 加载配置文件
	if err = a.loadTask(config); nil != err {
		return
	}
	// 监控配置文件
	if err = a.config.Watch(config, a); nil != err {
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

func (a *Agent) OnChanged(path string, from any, to any) {
	var diff string
	if diff = cmp.Diff(from, to, cmpopts.IgnoreFields(conf.Config{}, `secrets`)); `` == diff {
		return
	}

	fields := gox.Fields{
		field.String(`diff`, strings.TrimSpace(diff)),
		field.String(`path`, path),
	}
	a.logger.Info(`检测到配置有更新，重新加载任务`, fields...)
	if err := a.loadTask(to.(*conf.Config)); nil != err {
		a.logger.Error(`装载任务失败`, fields.Connect(field.Error(err))...)
	}
}

func (a *Agent) OnDeleted(path string) {
	a.logger.Error(`配置文件被删除，严重错误`, field.String(`path`, path))
}

func (a *Agent) OnError(path string, err error) {
	a.logger.Warn(`配置文件有错误`, field.String(`path`, path), field.Error(err))
}
