package main

import (
	"songjiang/common"
	"songjiang/sign"
)

// SongjiangJob 动态域名解析任务
type SongjiangJob struct {
	signer sign.Signer
	app    common.App
}

// Run 动态域名解析任务真正执行的方法
func (songjiangJob *SongjiangJob) Run() {

}
