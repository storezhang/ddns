package dns

import (
    `github.com/robfig/cron/v3`
    log `github.com/sirupsen/logrus`
)

func init() {
    crontab := cron.New(cron.WithSeconds())
    defer crontab.Stop()
    if id, err := crontab.AddFunc("@every 15m", func() {
        recordCache = nil
        log.Info("成功清空缓存")
    }); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Fatal("设置清空缓存任务失败")
    } else {
        log.WithFields(log.Fields{
            "id": id,
        }).Info("设置清空缓存任务成功")
    }
    crontab.Start()
}
