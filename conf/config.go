package conf

import (
	"sync"
)

// Config 配置
type Config struct {
	// 授权列表
	Secrets []*Secret `json:"secrets" yaml:"secrets" xml:"secrets" toml:"secrets" validate:"required"`
	// 域名列表
	Resolves []*resolve `json:"resolves" yaml:"resolves" xml:"resolves" toml:"resolves" validate:"required"`

	secrets *sync.Map
}

func (c *Config) Secret(label string) (secret *Secret) {
	if nil == c.secrets {
		c.secrets = new(sync.Map)
		for _, _secret := range c.Secrets {
			c.secrets.Store(_secret.Label, _secret)
		}
	}

	if cached, ok := c.secrets.Load(label); ok {
		secret = cached.(*Secret)
	}

	return
}
