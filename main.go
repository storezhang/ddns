package main

import (
    `flag`
    `fmt`
    `strings`
    `time`

    `github.com/jinzhu/configor`
    log `github.com/sirupsen/logrus`

    `songjiang/common`
    `songjiang/sign`
)

func main() {
    var confFilepath = flag.String("conf", "application.yml", "配置文件路径")
    flag.Parse()
    conf := &common.Config{}
    configor.New(&configor.Config{
        AutoReload:         true,
        Silent:             true,
        AutoReloadInterval: time.Minute,
    }).Load(conf, *confFilepath, "application.json", "application.ini")

    songjiang := conf.Songjiang
    if logLevel, err := log.ParseLevel(songjiang.LogLevel); nil != err {
        log.SetLevel(log.InfoLevel)
        log.WithFields(log.Fields{
            "err":      err,
            "logLevel": conf.Songjiang.LogLevel,
        }).Error("日志级别配置有误")
    } else {
        log.SetLevel(logLevel)
    }
    // 处理Debug
    if songjiang.Debug {
        songjiang.Redo = "5s"
    }

    for _, app := range conf.Apps {
        var signer sign.Signer

        switch strings.ToLower(app.Type) {
        case "hao4k":
            signer = &conf.Hao4k
        default:
            log.WithFields(log.Fields{"type": app.Type}).Fatal("不支持该类型，请重新配置")
        }

        // 增加启动立即执行
        songjiangJob := &SongjiangJob{signer: signer, songjiang: &songjiang, app: &app}
        // 真正的执行任务
        spec := fmt.Sprintf("@every %s", songjiang.Redo)
        if id, err := crontab.AddJob(spec, songjiangJob); nil != err {
            log.WithFields(log.Fields{
                "spec": spec,
                "err":  err,
            }).Error("添加Songjiang任务失败")
        } else {
            log.WithFields(log.Fields{
                "spec": spec,
                "id":   id,
            }).Info("添加Songjiang任务成功")
        }
    }

    crontab.Start()
    select {}
}
