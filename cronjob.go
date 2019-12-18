package main

import (
    `context`
    `fmt`
    `math/rand`
    `time`

    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
    `github.com/tj/go-naturaldate`

    `github.com/parnurzeal/gorequest`

    `songjiang/common`
    `songjiang/sign`
)

// SongjiangJob 动态域名解析任务
type SongjiangJob struct {
    ctx       context.Context
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

    if id, err := crontab.AddFunc(spec, func() {
        log.WithFields(log.Fields{
            "name":  job.app.Name,
            "start": job.app.StartTime,
            "end":   job.app.EndTime,
            "type":  job.app.Type,
            "spec":  spec,
            "jobId": jobId,
        }).Info("开始执行签到任务")

        result := job.signer.AutoSign(job.ctx, job.app.Cookies)
        // 通知用户，如果有设置消息推送
        notify(job.app, result)

        log.WithFields(log.Fields{
            "name":   job.app.Name,
            "start":  job.app.StartTime,
            "end":    job.app.EndTime,
            "type":   job.app.Type,
            "spec":   spec,
            "jobId":  jobId,
            "before": result.Before,
            "after":  result.After,
        }).Info("成功签到任务")
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

func notify(app *common.App, result sign.AutoSignResult) {
    for _, ch := range app.ServerChans {
        req.Post(fmt.Sprintf("https://sc.ftqq.com/%s.send", ch)).
            Type("multipart").
            Send(common.ServerChanRequest{
                Text: fmt.Sprintf("任务执行完成：%s - %s", app.Name, result.Msg),
                Desp: fmt.Sprintf("执行后的结果：%s\n执行前的状态：%s", result.After, result.Before),
            }).
            End()
    }
}
