package utils

import (
    `github.com/chromedp/chromedp`
    `github.com/storezhang/gos/chromedps`
)

func Sleep() chromedp.Action {
    return chromedps.Sleep("4s")
}
