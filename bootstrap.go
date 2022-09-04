package main

import (
	"ddns/dns"

	"github.com/pangum/pangu"
)

type (
	bootstrap struct {
		application *pangu.Application
		agent       *dns.Agent
	}

	bootstrapIn struct {
		pangu.In

		Application *pangu.Application
		Agent       *dns.Agent
	}
)

func newBootstrap(in bootstrapIn) pangu.Bootstrap {
	return &bootstrap{
		application: in.Application,
		agent:       in.Agent,
	}
}

func (b *bootstrap) Startup() error {
	return b.application.AddServes(b.agent)
}
