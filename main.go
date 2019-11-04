package main

import (
    "fmt"
    "io/ioutil"
    "time"

    "ddns/common"

    "github.com/robfig/cron"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v2"
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

    if logLevel, err := log.ParseLevel(conf.Niulang.LogLevel); nil != err {
        log.SetLevel(log.InfoLevel)
        log.WithFields(log.Fields{
            "err":      err,
            "logLevel": conf.Niulang.LogLevel,
        }).Fatal("日志级别配置有误")
    } else {
        log.SetLevel(logLevel)
    }

    crontab := cron.New()
    defer crontab.Stop()

    now := time.Now()
    spec := fmt.Sprintf("%d %d */1 * * *", now.Second()+5, now.Minute())
    if err := crontab.AddFunc(spec, func() {
        // running(conf)
    }); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("添加SSL更新计划任务失败")
    }

    crontab.Start()
    web.StartServer(conf)
}
