package sign

import (
    `context`

    `github.com/chromedp/chromedp`
    log `github.com/sirupsen/logrus`
    `github.com/storezhang/gos/chromedps`

    `songjiang/utils`
)

// Hao4k Hao4k对象
type Hao4k struct {
    SignSelector   string `default:"'#JD_sign'"`
    SignedSelector string `default:"//span[contains(@class, 'btn btnvisted')]"`
    SignUrl        string `default:"https://www.hao4k.cn/k_misign-sign.html"`
    ScoreUrl       string `default:"https://www.hao4k.cn/home.php?mod=spacecp&ac=credit&showcredit=1"`
    KBSelector     string `default:"//em[contains(text(), 'K币')]/.."`
}

func (hao4k *Hao4k) AutoSign(ctx context.Context, cookies string) (result AutoSignResult) {
    if err := chromedp.Run(ctx); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法启动浏览器实例")
    } else {
        log.Info("启动浏览器成功")
    }

    // 等待签到界面完成
    if err := chromedp.Run(
        ctx,
        chromedps.DefaultVisit(hao4k.SignUrl, cookies),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法载入签到界面")
    } else {
        log.Info("成功进入签到界面")
    }

    // 签到前的K币
    result.Before = getKB(ctx, hao4k)
    // 确认是否已经签到
    if err := chromedp.Run(
        ctx,
        chromedps.DefaultSleep(),
        chromedps.TasksWithTimeOut(&ctx, "10s", chromedp.Tasks{
            chromedp.Navigate(hao4k.SignUrl),
            chromedp.WaitVisible(hao4k.SignedSelector),
        }),
    ); nil != err {
        log.Info("还没有签到，继续执行自动签到任务")
    } else {
        // 签到后的K币
        result.Success = true
        result.After = result.Before
        result.Msg = "已签到，明天再来签到吧"

        log.WithFields(log.Fields{
            "cookies": cookies,
        }).Info("已签到，明天再来签到吧")

        return
    }

    // 点击签到按扭
    if err := chromedp.Run(
        ctx,
        chromedps.DefaultSleep(),
        chromedp.Navigate(hao4k.SignUrl),
        chromedps.DefaultSleep(),
        chromedp.Click(hao4k.SignSelector, chromedp.NodeVisible),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法点击签到按扭")
    } else {
        // 签到后的K币
        result.Success = true
        result.After = getKB(ctx, hao4k)
        result.Msg = "签到成功"
        log.Info("成功点击签到按扭")
    }

    return
}

func getKB(ctx context.Context, hao4k *Hao4k) (kb string) {
    if err := chromedp.Run(
        ctx,
        utils.Sleep(),
        chromedp.Navigate(hao4k.ScoreUrl),
        chromedp.Text(hao4k.KBSelector, &kb),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法获得当前K币数据")
    } else {
        log.WithFields(log.Fields{
            "currentKB": kb,
        }).Info("成功获得当前K币数据")
    }

    return
}
