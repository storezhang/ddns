package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    "time"

    "github.com/chromedp/chromedp"
    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v2"

    "songjiang/common"
    "songjiang/sign"
)

func main() {
    var confFilepath = flag.String("conf", "songjiang.yml", "配置文件路径")
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

    songjiang := conf.Songjiang
    if logLevel, err := log.ParseLevel(songjiang.LogLevel); nil != err {
        log.SetLevel(log.InfoLevel)
        log.WithFields(log.Fields{
            "err":      err,
            "logLevel": conf.Songjiang.LogLevel,
        }).Fatal("日志级别配置有误")
    } else {
        log.SetLevel(logLevel)
    }

    crontab := cron.New(cron.WithSeconds())
    defer crontab.Stop()

    for _, app := range conf.Apps {
        var signer sign.Signer

        switch strings.ToLower(app.Type) {
        case "hao4k":
            signer = &conf.Hao4k
        default:
            log.WithFields(log.Fields{"type": app.Type}).Fatal("不支持该类型，请重新配置")
        }

        // 增加启动立即执行
        songjiangJob := &SongjiangJob{signer: signer, app: app}
        now := time.Now()
        now = now.Add(time.Second * 3)
        spec := fmt.Sprintf(
            "%d %d %d %d %d %d",
            now.Second(), now.Minute(), now.Hour(),
            now.Day(), now.Month(), now.Weekday(),
        )
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
        // 真正的执行任务
        spec = fmt.Sprintf("@every %s", songjiang.Redo)
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

func context(songjiang *common.Songjiang) context.Context {
    opts := append(
        chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoDefaultBrowserCheck,
        chromedp.NoSandbox,
        chromedp.Flag("ignore-certificate-errors", true),
    )
    if songjiang.Debug {
        opts = append(opts, chromedp.Flag("headless", false))
        opts = append(opts, chromedp.Flag("start-maximized", true))
    } else {
        opts = append(opts, chromedp.Headless)
        opts = append(opts, chromedp.WindowSize(songjiang.BrowserWidth, songjiang.BrowserHeight))
    }

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    ctx, cancel := chromedp.NewContext(
        allocCtx,
        chromedp.WithLogf(log.Printf),
        chromedp.WithDebugf(log.Debugf),
        chromedp.WithErrorf(log.Errorf),
    )
    defer cancel()

    if duration, err := time.ParseDuration(songjiang.BrowserTimeout); nil != err {
        ctx, cancel = context.WithTimeout(ctx, duration)
        defer cancel()
    } else {
        log.WithFields(log.Fields{
            "browserTimeout": songjiang.BrowserTimeout,
        }).Warn("browserTimeout配置有错误")
    }

    return ctx
}
