package utils

import (
    `github.com/chromedp/chromedp`
    `github.com/storezhang/gos/chromedps`
)

// Sleep 步骤之间的休眠
func Sleep() chromedp.Action {
    return chromedps.Sleep("10s")
}
