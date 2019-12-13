package utils

import (
    "strings"

	"github.com/chromedp/chromedp"
)

func SelectorType(selectorType string) chromedp.QueryOption {
	selectType := chromedp.ByID

	switch strings.ToLower(strings.TrimSpace(selectorType)) {
	case "id":
		selectType = chromedp.ByID
	default:
		selectType = chromedp.ByID
	}

	return selectType
}
