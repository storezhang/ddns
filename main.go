package main

import (
	"github.com/pangum/pangu"
)

func main() {
	panic(pangu.New(
		pangu.Named(`ddns`),
		pangu.Banner(`DDNS`, pangu.BannerTypeAscii),
		pangu.Author(`storezhang`, `storezhang@gmail.com`),
		pangu.Usage(`动态载解析，使用请参数https://github.com/storezhang/ddns的相关说明`),
	).Run(newBootstrap))
}
