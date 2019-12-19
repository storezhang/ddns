package utils

import (
    `context`
    `fmt`

    `github.com/chromedp/chromedp`
    `github.com/storezhang/gos/chromedps`
)

func Sleep() chromedp.Action {
    return chromedps.Sleep("4s")
}

func CheckDisplay(ctx context.Context, url string, selector string) (display bool, err error) {
    var displayTypes []string
    if err = chromedp.Run(
        ctx,
        chromedp.Navigate(url),
        chromedp.Evaluate(
            fmt.Sprintf(`getComputedStyle(document.querySelector("%s"), null).display;`, selector),
            &displayTypes,
        ),
    ); nil == err {
        if nil == displayTypes || 1 != len(displayTypes) {
            display = false
        } else {
            display = "none" != displayTypes[0]
        }
    }

    return
}
