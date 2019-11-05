package common

import (
    "ddns/dns"
)

type Config struct {
    DDNS    DDNS       `yaml:"ddns"`
    Aliyun  dns.Aliyun `yaml:"aliyun"`
    Domains []Domain   `domains`
}

type DDNS struct {
    Debug      bool   `yaml:"debug"`
    LogLevel   string `yaml:"logLevel"`
    TimeFormat string `yaml:"timeFormat"`
    Redo       string `yaml:"redo"`
}

type Domain struct {
    Name       string `yaml:"name"`
    SubDomains string `yaml:"subDomains"`
    Type       string `yaml:"type"`
    DNSTypes   string `yaml:"dnsTypes"`
    Value      string `yaml:"value"`
    TTL        int    `yaml:"ttl"`
}

func (ddns *DDNS) UnmarshalYAML(unmarshal func(interface{}) error) error {
    type rawType DDNS
    raw := rawType{
        Debug:    false,
        LogLevel: "info",
        Redo:     "1m",
    }
    if err := unmarshal(&raw); nil != err {
        return err
    }

    *ddns = DDNS(raw)

    return nil
}

func (domain *Domain) UnmarshalYAML(unmarshal func(interface{}) error) error {
    type rawType Domain
    raw := rawType{
        Value: "",
        TTL:   600,
    }
    if err := unmarshal(&raw); nil != err {
        return err
    }

    *domain = Domain(raw)

    return nil
}
