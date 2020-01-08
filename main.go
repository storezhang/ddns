package main

import (
    `flag`
    `fmt`
    `os`
    `time`

    `github.com/jinzhu/configor`
    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`

    `ddns/common`
    `ddns/dns`
)

func main() {
    var confFilepath = flag.String("conf", "application.yml", "配置文件路径")
    flag.Parse()
    conf := &common.Config{}
    configor.New(&configor.Config{
        AutoReload:         true,
        Silent:             true,
        AutoReloadInterval: time.Minute,
        AutoReloadCallback: func(config interface{}) {
            log.Info("配置解析的域名有变化，退出程序重新解析")
            os.Exit(0)
        },
    }).Load(conf, *confFilepath, "application.json", "application.toml")

    ddns := conf.DDNS
    if logLevel, err := log.ParseLevel(ddns.LogLevel); nil != err {
        log.SetLevel(log.InfoLevel)
        log.WithFields(log.Fields{
            "err":      err,
            "logLevel": ddns.LogLevel,
        }).Warn("日志级别配置有误，已修复成Info级别")
    } else {
        log.SetLevel(logLevel)
    }

    crontab = cron.New(cron.WithSeconds())
    defer crontab.Stop()
    for _, domain := range conf.Domains {
        var resolver dns.Resolver

        switch domain.Type {
        case "aliyun":
            resolver = &conf.Aliyun
        default:
            log.WithFields(log.Fields{"type": domain.Type}).Fatal("不支持该类型，请重新配置")
        }

        // 增加启动立即执行
        ddnsJob := &DDNSJob{resolver: resolver, domain: &domain, ddns: &ddns}
        // 真正的执行任务
        spec := fmt.Sprintf("@every %s", domain.Redo)
        if id, err := crontab.AddJob(spec, ddnsJob); nil != err {
            log.WithFields(log.Fields{
                "domain":     domain.Name,
                "subDomains": domain.SubDomains,
                "spec":       spec,
                "err":        err,
            }).Error("添加DDNS任务失败")
        } else {
            log.WithFields(log.Fields{
                "domain":     domain.Name,
                "subDomains": domain.SubDomains,
                "spec":       spec,
                "id":         id,
            }).Info("添加DDNS任务成功")
        }
    }

    crontab.Start()
    select {}
}
