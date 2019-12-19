package main

import (
    `context`
    `flag`
    `fmt`
    `strings`
    `time`

    `github.com/chromedp/chromedp`
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

    for _, app := range conf.Apps {
        var signer sign.Signer

        switch strings.ToLower(app.Type) {
        case "hao4k":
            signer = &conf.Hao4k
        default:
            log.WithFields(log.Fields{"type": app.Type}).Fatal("不支持该类型，请重新配置")
        }

        // 增加启动立即执行
        songjiangJob := &SongjiangJob{ctx: ctx(&songjiang), signer: signer, songjiang: &songjiang, app: &app}
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

func ctx(songjiang *common.Songjiang) context.Context {
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

    allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
    // defer cancel()

    ctx, _ := chromedp.NewContext(
        allocCtx,
        chromedp.WithLogf(log.Printf),
        chromedp.WithDebugf(log.Debugf),
        chromedp.WithErrorf(log.Errorf),
    )
    // defer cancel()

    if !songjiang.Debug {
        if duration, err := time.ParseDuration(songjiang.BrowserTimeout); nil != err {
            log.WithFields(log.Fields{
                "browserTimeout": songjiang.BrowserTimeout,
            }).Warn("browserTimeout配置有错误")
        } else {
            ctx, _ = context.WithTimeout(ctx, duration)
            // defer cancel()
        }
    }

    return ctx
}
