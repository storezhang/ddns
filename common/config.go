package common

import (
    `ddns/dns`
)

// Config 程序整体配置
type Config struct {
    DDNS    DDNS
    Aliyun  dns.Aliyun
    Domains []Domain
}

// DDNS 程序整体配置
type DDNS struct {
    Debug    bool   `default:"false"`
    LogLevel string `default:"info" yaml:"logLevel" toml:"logLevel"`
    Chans    []ServerChan
    Template Template
}

// Domain 每个域名的配置
type Domain struct {
    Name            string `required:"true"`
    Chans           []ServerChan
    Template        Template
    SubDomains      []string `required:"true" yaml:"subDomains" toml:"subDomains"`
    SubDomainPrefix string   `yaml:"subDomainPrefix" toml:"subDomainPrefix"`
    SubDomainStaff  string   `yaml:"subDomainStaff" toml:"subDomainStaff"`
    Type            string   `default:"aliyun" required:"true"`
    DNSTypes        []string `required:"true" yaml:"dnsTypes" toml:"dnsTypes"`
    Value           string
    TTL             int    `default:"600"`
    Redo            string `default:"1m"`
}

// Template 模板配置
// 用于：推送
type Template struct {
    Title   string `default:"'解析后：{{.Result.After}}，解析前{{.Result.Before}}'"`
    Content string `default:"'解析域名：{{.SubDomain}}.{{.Domain.Name}}'"`
}
