package conf

import (
	"github.com/goexl/gox"
)

// Secret 授权
type Secret struct {
	gox.Secret `yaml:",inline"`

	// 类型
	Type typ `default:"aliyun" json:"type" yaml:"type" xml:"type" toml:"type" validate:"required,oneof=aliyun"`
	// 标签
	Label string `default:"aliyun" json:"label" yaml:"label" xml:"label" toml:"label" validate:"required"`
}
