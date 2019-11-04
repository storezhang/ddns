package common

import (
    "ddns/dns"
)

type Config struct {
    Debug      bool       `yaml:"debug"`
    LogLevel   string     `yaml:"logLevel"`
    TimeFormat string     `yaml:"timeFormat"`
    DDNS       []DDNS     `yaml:"ddns"`
    Aliyun     dns.Aliyun `yaml:"aliyun"`
}

type DDNS struct {
    Domain    string `yaml:"domain"`
    SubDomain string `yaml:"subDomain"`
    Redo      int    `yaml:"redo"`
    Type      string `yaml:"type"`
}

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
