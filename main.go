package main

import (
	"github.com/pangum/pangu"
)

func main() {
	panic(pangu.New(
		pangu.Banner(`DDNS`, pangu.BannerTypeAscii),
	).Run(newBootstrap))
}
