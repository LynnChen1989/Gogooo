package main

import "github.com/robfig/cron"

func PauseTimer() {
	// 停服定时器，于每天23:59:50定时发布停服公告
	c := cron.New()

	// 23点50分关闭服务入口
	spec01 := "0 50 23 * * ?"
	c.AddFunc(spec01, func() {
		PauseSrv()
	})

	// 0点2分检查日切状态，日切成功后断开主从
	spec02 := "0 2 0 * * ?"
	c.AddFunc(spec02, func() {
		cd := &CutDate{}
		if cd.CheckCutDateStatus() {
			BatchHandleDbSlave("stop")
		}
	})

	// 每隔30秒检查一次日终状态
	spec03 := "*/30 3-59 0 * * ?"
	c.AddFunc(spec03, func() {
		cd := &CutDate{}
		if cd.CheckCutEndStatus() {
			BatchHandleDbSlave("stop")
		}
	})

	spec04 := "*/5 * * * * ?"
	c.AddFunc(spec04, func() {
		getTodayOdsAllStatus()
	})

	c.Start()
	select {}
}
