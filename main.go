package main

import (
    "fmt"
    "io/ioutil"
    "time"

    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v2"

    "ddns/common"
)

func main() {
    conf := &common.Config{}

    configData, err := ioutil.ReadFile("config.yml")
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

    crontab := cron.New()
    defer crontab.Stop()

    now := time.Now()
    spec := fmt.Sprintf("%d %d */1 * * *", now.Second()+5, now.Minute())
    if _, err := crontab.AddFunc(spec, func() {
        running(conf)
    }); nil != err {
        log.WithFields(log.Fields{
            "spec": spec,
            "err":  err,
        }).Error("添加DDNS任务失败")
    }

    crontab.Start()
}

func running(conf *common.Config) {
    for _, domain := range conf.Domains {
        switch domain.Type {
        case "aliyun":
            log.WithFields(log.Fields{
                "name":      domain.Name,
                "subDomain": domain.SubDomains,
                "dnsTypes":  domain.DNSTypes,
                "type":      domain.Type,
            }).Debug("执行DNS解析更新")
        }
    }
}
