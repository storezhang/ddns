package sign

import (
    `context`

    `github.com/chromedp/chromedp`
    log `github.com/sirupsen/logrus`
    `github.com/storezhang/gos/chromedps`
    `github.com/storezhang/gos/webs`

    `songjiang/utils`
)

// Hao4k Hao4k对象
type Hao4k struct {
    SignSelector string `default:"'#JD_sign'"`
    SignUrl      string `default:"https://www.hao4k.cn/k_misign-sign.html"`
    ScoreUrl     string `default:"https://www.hao4k.cn/home.php?mod=spacecp&ac=credit&showcredit=1"`
    KBSelector   string `default:"//em[contains(text(), 'K币')]/.."`
}

func (hao4k *Hao4k) AutoSign(ctx context.Context, cookies string) (result AutoSignResult) {
    if err := chromedp.Run(ctx); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法启动浏览器实例")
    } else {
        log.Info("启动浏览器成功")
    }

    // 等待登录界面完成
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

    // 签到前的积分
    result.Before = getKB(ctx, hao4k)

    // 点击签到按扭
    if err := chromedp.Run(
        ctx,
        chromedp.Click(webs.ID(hao4k.SignSelector), chromedp.NodeEnabled),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法点击签到按扭")
    } else {
        log.Info("成功点击签到按扭")
    }

    // 签到后的积分
    result.After = getKB(ctx, hao4k)

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
