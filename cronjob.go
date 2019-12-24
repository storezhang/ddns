package main

import (
    `bytes`
    `context`
    `fmt`
    `html/template`
    `math/rand`
    `strings`
    `time`

    `github.com/chromedp/chromedp`
    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
    `github.com/tj/go-naturaldate`

    `github.com/parnurzeal/gorequest`

    `songjiang/common`
    `songjiang/sign`
)

// SongjiangJob 动态域名解析任务
type SongjiangJob struct {
    signer    sign.Signer
    songjiang *common.Songjiang
    app       *common.App
}

var crontab *cron.Cron
var jobCache map[string]cron.EntryID

var req *gorequest.SuperAgent

func initJobCache() {
    jobCache = make(map[string]cron.EntryID)
}
func init() {
    // 初始化Http客户端
    req = gorequest.New()
    // 初始化随机数
    rand.Seed(time.Now().UnixNano())
    // 初始化缓存
    initJobCache()

    crontab = cron.New(cron.WithSeconds())
    defer crontab.Stop()
    if id, err := crontab.AddFunc("59 59 23 * * *", func() {
        initJobCache()
        log.Info("成功清空缓存")
    }); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Fatal("设置每日清空缓存任务失败")
    } else {
        log.WithFields(log.Fields{
            "id": id,
        }).Info("设置每日清空缓存任务成功")
    }
    crontab.Start()
}

var base = time.Now()

// Run 动态域名解析任务真正执行的方法
func (job *SongjiangJob) Run() {
    if _, ok := jobCache[job.app.Cookies]; !ok {
        jobId := addJob(job)
        jobCache[job.app.Cookies] = jobId
    }
}

func addJob(job *SongjiangJob) (jobId cron.EntryID) {
    var startNano int64
    var endNano int64
    if start, err := naturaldate.Parse(job.app.StartTime, base); nil != err {
        log.WithFields(log.Fields{
            "name": job.app.Name,
            "time": job.app.StartTime,
            "err":  err,
        }).Fatal("应用开始时间设置有问题")
    } else {
        startNano = start.UnixNano()
    }
    if end, err := naturaldate.Parse(job.app.EndTime, base); nil != err {
        log.WithFields(log.Fields{
            "name": job.app.Name,
            "time": job.app.EndTime,
            "err":  err,
        }).Fatal("应用结束时间设置有问题")
    } else {
        endNano = end.UnixNano()
    }

    now := time.Now()
    nowNano := now.UnixNano()

    var runTime time.Time
    if job.songjiang.Debug {
        runTime = now.Add(3 * time.Second)
    } else {
        if nowNano > endNano {
            jobId = cron.EntryID(now.Nanosecond())
            log.WithFields(log.Fields{
                "name":  job.app.Name,
                "start": job.app.StartTime,
                "end":   job.app.EndTime,
                "jobId": jobId,
            }).Warn("应用运行时机已过，不再执行")
            return
        }

        if nowNano > startNano {
            runTime = now.Add(time.Duration(rand.Int63n(endNano - nowNano)))
        } else {
            runTime = now.Add(time.Duration(startNano - nowNano + rand.Int63n(endNano-startNano)))
        }
    }
    spec := fmt.Sprintf(
        "%d %d %d %d %d %d",
        runTime.Second(), runTime.Minute(), runTime.Hour(),
        runTime.Day(), runTime.Month(), runTime.Weekday(),
    )

    if id, err := crontab.AddJob(spec, &AutoSignJob{
        signer:    job.signer,
        app:       job.app,
        songjiang: job.songjiang,
        cookies:   job.app.Cookies,
    }); nil != err {
        log.WithFields(log.Fields{
            "name":  job.app.Name,
            "start": job.app.StartTime,
            "end":   job.app.EndTime,
            "spec":  spec,
            "err":   err,
        }).Error("添加签到任务失败")
    } else {
        jobId = id
        log.WithFields(log.Fields{
            "name":  job.app.Name,
            "start": job.app.StartTime,
            "end":   job.app.EndTime,
            "spec":  spec,
            "jobId": jobId,
        }).Info("添加签到任务成功")
    }

    return
}

// AutoSignJob 自动签到Job
type AutoSignJob struct {
    signer    sign.Signer
    app       *common.App
    songjiang *common.Songjiang
    cookies   string
}

// Run 自动签到执行任务
func (job *AutoSignJob) Run() {
    opts := append(
        chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoDefaultBrowserCheck,
        chromedp.NoSandbox,
        chromedp.Flag("ignore-certificate-errors", true),
    )
    if job.songjiang.Debug {
        opts = append(opts, chromedp.Flag("headless", false))
        opts = append(opts, chromedp.Flag("start-maximized", true))
    } else {
        opts = append(opts, chromedp.Headless)
        opts = append(opts, chromedp.WindowSize(job.songjiang.BrowserWidth, job.songjiang.BrowserHeight))
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

    log.WithFields(log.Fields{
        "name":  job.app.Name,
        "start": job.app.StartTime,
        "end":   job.app.EndTime,
        "type":  job.app.Type,
    }).Info("开始执行签到任务")

    result := job.signer.AutoSign(ctx, job.app.Cookies)
    // 通知用户，如果有设置消息推送
    notify(job.app, job.songjiang, &result)

    log.WithFields(log.Fields{
        "name":   job.app.Name,
        "start":  job.app.StartTime,
        "end":    job.app.EndTime,
        "type":   job.app.Type,
        "before": result.Before,
        "after":  result.After,
    }).Info("成功签到任务")
}

func notify(app *common.App, songjiang *common.Songjiang, result *sign.AutoSignResult) {
    var serverChans []common.ServerChan
    if nil != app.Chans && 0 < len(app.Chans) {
        serverChans = app.Chans
    } else if nil != songjiang.Chans && 0 < len(songjiang.Chans) {
        serverChans = songjiang.Chans
    } else {
        return
    }

    var titleTemplate string
    var contentTemplate string
    if "" == strings.TrimSpace(app.Template.Title) && "" == strings.TrimSpace(app.Template.Context) {
        titleTemplate = app.Template.Title
        contentTemplate = app.Template.Context
    } else if "" == strings.TrimSpace(songjiang.Template.Title) && "" == strings.TrimSpace(songjiang.Template.Context) {
        titleTemplate = songjiang.Template.Title
        contentTemplate = songjiang.Template.Context
    } else {
        return
    }

    notifyToUser(serverChans, titleTemplate, contentTemplate, app, result)
}

type notifyData struct {
    App    *common.App
    Result *sign.AutoSignResult
}

func notifyToUser(
    chans []common.ServerChan,
    titleTemplate string,
    contentTemplate string,
    app *common.App,
    result *sign.AutoSignResult,
) {
    data := &notifyData{
        App:    app,
        Result: result,
    }
    title := render("title", titleTemplate, data)
    desp := render("desp", contentTemplate, data)
    // 真正发推送
    for _, ch := range chans {
        rsp, body, errs := req.Post(fmt.Sprintf("https://sc.ftqq.com/%s.send", ch.Key)).
            Type("form").
            Send(common.ServerChanRequest{
                Text: title,
                Desp: desp,
            }).End()

        if nil != errs {
            log.WithFields(log.Fields{
                "name": app.Name,
                "chan": ch.Key,
                "rsp":  rsp,
                "body": body,
                "errs": errs,
            }).Info("ServerChan推送消息失败")
        } else {
            log.WithFields(log.Fields{
                "name": app.Name,
                "chan": ch.Key,
            }).Info("ServerChan推送消息成功")
        }
    }
}

func render(name string, tpl string, data interface{}) (result string) {
    tmpl, err := template.New(name).Parse(tpl)
    if err != nil {
        result = ""
    }
    buf := new(bytes.Buffer)
    err = tmpl.Execute(buf, data)
    if err != nil {
        result = ""
    }
    result = buf.String()

    return
}
