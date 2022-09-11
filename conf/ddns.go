package conf

import "time"

type ddns struct {
	// 间隔
	// nolint:lll
	Interval time.Duration `default:"5s" json:"interval" yaml:"interval" xml:"interval" toml:"interval" validate:"required"`
}
