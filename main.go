package main

import (
    `flag`
    `fmt`
    `io/ioutil`
    `time`

    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
    `gopkg.in/yaml.v2`

    `ddns/common`
    `ddns/dns`
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

    for _, domain := range conf.Domains {
        var resolver dns.Resolver

        switch domain.Type {
        case "aliyun":
            resolver = &conf.Aliyun
        default:
            log.WithFields(log.Fields{"type": domain.Type}).Fatal("不支持该类型，请重新配置")
        }

        // 增加启动立即执行
        ddnsJob := &DDNSJob{resolver: resolver, domain: domain}
        now := time.Now()
        now = now.Add(time.Second * 3)
        spec := fmt.Sprintf(
            "%d %d %d %d %d %d",
            now.Second(), now.Minute(), now.Hour(),
            now.Day(), now.Month(), now.Weekday(),
        )
        if id, err := crontab.AddJob(spec, ddnsJob); nil != err {
            log.WithFields(log.Fields{
                "domain": domain.Name,
                "spec":   spec,
                "err":    err,
            }).Error("添加DDNS任务失败")
        } else {
            log.WithFields(log.Fields{
                "domain":     domain.Name,
                "subDomains": domain.SubDomains,
                "spec":       spec,
                "id":         id,
            }).Info("添加DDNS任务成功")
        }
        // 真正的执行任务
        spec = fmt.Sprintf("@every %s", domain.Redo)
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
