package conf

import (
	"time"

	"github.com/goexl/uda"
)

type resolve struct {
	// 域名名称
	Name string `json:"name" yaml:"name" xml:"name" toml:"name" validate:"required"`
	// 子域名列表
	Subdomains []string `json:"subdomains" yaml:"subdomains" toml:"subdomains" validate:"required"`
	// 标签
	Label string `default:"aliyun" json:"label" yaml:"label" xml:"label" toml:"label" validate:"required"`
	// 类型列表
	Types []uda.Type `json:"types" yaml:"types" xml:"types" toml:"types" validate:"required,dive,oneof=CNAME A AAAA"`
	// 记录值
	Value string `json:"value" yaml:"value" xml:"value" toml:"value"`
	// 生存时间
	Ttl time.Duration `default:"10m" json:"ttl" yaml:"ttl" xml:"ttl" toml:"ttl" validate:"required"`
	// 前缀
	Prefix string `json:"prefix" yaml:"prefix" xml:"prefix" toml:"prefix"`
	// 后缀
	Staff string `json:"staff" yaml:"staff" xml:"staff" toml:"staff"`
}

func (r *resolve) Contains(typ uda.Type) (contains bool) {
	for _, _typ := range r.Types {
		if _typ == typ {
			contains = true
		}

		if contains {
			break
		}
	}

	return
}
