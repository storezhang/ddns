package common

import (
	"ddns/sign"
)

// Config 程序整体配置
type Config struct {
	DDNS    DDNS        `yaml:"ddns"`
	Aliyun  sign.Aliyun `yaml:"aliyun"`
	Domains []Domain    `yaml:"domains"`
}

// DDNS DDNS的配置
type DDNS struct {
	Debug      bool   `yaml:"debug"`
	LogLevel   string `yaml:"logLevel"`
	TimeFormat string `yaml:"timeFormat"`
}

// Domain 每个域名的配置
type Domain struct {
	Name            string `yaml:"name"`
	SubDomains      string `yaml:"subDomains"`
	SubDomainPrefix string `yaml:"subDomainPrefix"`
	SubDomainStaff  string `yaml:"subDomainStaff"`
	Type            string `yaml:"type"`
	DNSTypes        string `yaml:"dnsTypes"`
	Value           string `yaml:"value"`
	TTL             int    `yaml:"ttl"`
	Redo            string `yaml:"redo"`
}

// UnmarshalYAML 从YAML反序列化成DDNS对象时的默认值处理
func (ddns *DDNS) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawType DDNS
	raw := rawType{
		Debug:    false,
		LogLevel: "info",
	}
	if err := unmarshal(&raw); nil != err {
		return err
	}

	*ddns = DDNS(raw)

	return nil
}

// UnmarshalYAML 从YAML反序列化成域名对象时的默认值处理
func (domain *Domain) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawType Domain
	raw := rawType{
		SubDomainPrefix: "",
		SubDomainStaff:  "",
		Value:           "",
		TTL:             600,
		Redo:            "1m",
	}
	if err := unmarshal(&raw); nil != err {
		return err
	}

	*domain = Domain(raw)

	return nil
}
