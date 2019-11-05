package main

import (
    `flag`
    `fmt`
    `io/ioutil`
    `strings`
    `time`

    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
    `gopkg.in/yaml.v2`

    `github.com/storezhang/gos/nets`

    `ddns/common`
)

func main() {
    var confFilepath = flag.String("conf", "ddns.yml", "配置文件路径")
    flag.Parse()
    conf := &common.Config{}

    configData, err := ioutil.ReadFile(*confFilepath)
    if nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Fatal("加载配置文件出错")
    }

    err = yaml.Unmarshal(configData, &conf)
    if nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Fatal("配置文件出错")
    }

    if logLevel, err := log.ParseLevel(conf.DDNS.LogLevel); nil != err {
        log.SetLevel(log.InfoLevel)
        log.WithFields(log.Fields{
            "err":      err,
            "logLevel": conf.DDNS.LogLevel,
        }).Fatal("日志级别配置有误")
    } else {
        log.SetLevel(logLevel)
    }

    crontab := cron.New(cron.WithSeconds())
    defer crontab.Stop()

    // 增加启动执行
    now := time.Now()
    now.Add(time.Second * 5)
    spec := fmt.Sprintf(
        "%d %d %d %d %d %d",
        now.Second(), now.Minute(), now.Hour(),
        now.Day(), now.Month(), now.Weekday(),
    )
    if _, err := crontab.AddFunc(spec, func() {
        running(conf)
    }); nil != err {
        log.WithFields(log.Fields{
            "spec": spec,
            "err":  err,
        }).Error("添加DDNS任务失败")
    }
    // 真正的执行任务
    if _, err := crontab.AddFunc(fmt.Sprintf("@every %s", conf.DDNS.Redo), func() {
        running(conf)
    }); nil != err {
        log.WithFields(log.Fields{
            "spec": spec,
            "err":  err,
        }).Error("添加DDNS任务失败")
    }

    crontab.Start()
    select {}
}

func running(conf *common.Config) {
    for _, domain := range conf.Domains {
        for _, dnsType := range strings.Split(domain.DNSTypes, ",") {
            value := getValue(dnsType, domain)
            for _, subDomain := range strings.Split(domain.SubDomains, ",") {
                log.WithFields(log.Fields{
                    "name":      domain.Name,
                    "subDomain": subDomain,
                    "dnsType":   dnsType,
                    "type":      domain.Type,
                    "value":     value,
                }).Debug("执行DNS解析更新")

                switch domain.Type {
                case "aliyun":
                    conf.Aliyun.Resolve(
                        domain.Name,
                        subDomain,
                        value,
                        dnsType,
                        domain.TTL,
                    )
                }
            }
        }
    }
}

func getValue(dnsType string, domain common.Domain) string {
    var value string

    switch dnsType {
    case "A":
        if ip, err := nets.GetPublicIp(); nil != err {
            log.WithFields(log.Fields{
                "name":      domain.Name,
                "subDomain": domain.SubDomains,
                "dnsTypes":  domain.DNSTypes,
                "err":       err,
            }).Error("解析本机IP地址出错")
        } else {
            value = ip
        }
    default:
        value = domain.Value
    }

    return value
}
