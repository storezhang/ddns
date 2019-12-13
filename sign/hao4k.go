package sign

import (
    "context"

    "github.com/chromedp/chromedp"
    log "github.com/sirupsen/logrus"
    "github.com/storezhang/gos/chromedps"

    "songjiang/utils"
)

// Hao4k Hao4k对象
type Hao4k struct {
    SignSelector string `yaml:"signSelector"`
    SelectorType string `yaml:"selectorType"`
    SignUrl      string `yaml:"signUrl"`
}

// UnmarshalYAML 从YAML反序列化成域名对象时的默认值处理
func (hao4k *Hao4k) UnmarshalYAML(unmarshal func(interface{}) error) error {
    type rawType Hao4k
    raw := rawType{
        SignSelector: "#JD_sign",
        SignUrl:      "https://www.hao4k.cn//k_misign-sign.html",
    }
    if err := unmarshal(&raw); nil != err {
        return err
    }

    *hao4k = Hao4k(raw)

    return nil
}

func (hao4k *Hao4k) AutoSign(ctx context.Context, cookies string) {
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
        chromedps.DefaultCookies(hao4k.SignUrl, cookies),
        chromedp.WaitVisible(hao4k.SignSelector, utils.SelectorType(hao4k.SelectorType)),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法载入登录界面")
    } else {
        log.Info("成功进入登录界面")
    }

    // 点击签到按扭
    if err := chromedp.Run(
        ctx,
        chromedp.Click("#JD_sign", chromedp.ByID),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法载入登录界面")
    } else {
        log.Info("成功进入登录界面")
    }
}
