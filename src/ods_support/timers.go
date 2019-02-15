package main

import "github.com/robfig/cron"

func PauseTimer() {
	// 停服定时器，于每天23:59:50定时发布停服公告
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		getTodayOdsAllStatus()
	})
	c.Start()
	select {}
}
