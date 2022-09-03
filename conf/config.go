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

	secretCache *sync.Map
}

func (c *Config) Secret(label string) (secret *Secret) {
	if nil == c.secretCache {
		c.secretCache = new(sync.Map)
		for _, _secret := range c.Secrets {
			c.secretCache.Store(_secret.Label, _secret)
		}
	}

	if _secret, ok := c.secretCache.Load(label); ok {
		secret = _secret.(*Secret)
	}

	return
}
