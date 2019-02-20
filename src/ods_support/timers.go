package main

import (
	"github.com/robfig/cron"
	"os"
	"time"
)

func OdsTimer() {
	// 停服定时器，于每天23:59:50定时发布停服公告
	today := NowFormatDate("20060102")
	client := RedisClient(
		os.Getenv("CACHE_REDIS_HOST"),
		os.Getenv("CACHE_REDIS_PASSWORD"),
		os.Getenv("CACHE_REDIS_DB"))

	c := cron.New()
	// 23点50分关闭服务入口
	//spec01 := "0 50 23 * * ?"
	//DEBUG
	spec01 := os.Getenv("CRON_STOP_SRV")
	c.AddFunc(spec01, func() {
		Info.Println("-------- [*]STEP: start close CROS service entry point -------- ")
		PauseSrv()
		status := PauseSrv()
		if status == "success" {
			SendNotify("general", "ODS_SUPPORT:STOP CROS"+NowFormatDate("20060102"), "关闭服务入口[成功]")
		} else if status == "failed" {
			SendNotify("fatal", "ODS_SUPPORT:STOP CROS"+NowFormatDate("20060102"), "关闭服务入口[失败]")
		}
	})

	// 0点2分检查日切状态，日切成功后断开核心和管理主从关系
	//spec02 := "0 2 0 * * ?"
	spec02 := os.Getenv("CRON_CUT_DATE")
	c.AddFunc(spec02, func() {
		cd := &CutDate{}
		if cd.CheckCutDateStatus() {
			client.Set("cutdate.status."+today, "success", 0)
			//屏蔽告警
			//api := getAPI()
			Info.Println("-------- [*]STEP: prevent zabbix alert about MySQL master-slave -------- ")
			//api.PauseAlert()
			SendNotify("general", "ODS_SUPPORT:屏蔽告警"+NowFormatDate("20060102"), "[BEGIN]屏蔽主从异常告警")

			time.Sleep(30 * time.Second)
			//断开核心，管理主从
			Info.Println("-------- [*]STEP: pause loan CAS and CMS MySQL master-slave -------- ")
			SendNotify("general", "ODS_SUPPORT:断开主从"+NowFormatDate("20060102"), "[BEGIN]开始断开信贷核心{CAS}信贷管理{CMS}主从")
			//BatchHandleDbSlave("cascms", "stop")

			//打开服务入口
			Info.Println("-------- [*]STEP: pause loan CAS and CMS MySQL master-slave -------- ")
			SendNotify("general", "ODS_SUPPORT:START CROS"+NowFormatDate("20060102"), "[BEGIN]恢复服务入口")
			status := RestoreSrv()
			if status == "success" {
				SendNotify("general", "ODS_SUPPORT:START CROS"+NowFormatDate("20060102"), "[SUCCESS]开启服务入口成功")
			} else if status == "failed" {
				SendNotify("fatal", "ODS_SUPPORT:START CROS"+NowFormatDate("20060102"), "[FAILURE]开启服务入口成功")
			}
		} else {
			client.Set("cutdate.status."+today, "failure", 0)
			SendNotify("fatal", "ODS_SUPPORT:CUT DATE"+NowFormatDate("20060102"), "[FAILURE]日切失败")
		}
	})

	// 每隔30秒检查一次日终状态, 日终完成断开核算主从，后通知下游服务, 这个要避免周期内的重复推送
	//spec03 := "*/60 3-59 0 * * ?"
	spec03 := os.Getenv("CRON_CUT_END")
	c.AddFunc(spec03, func() {
		cd := &CutDate{}
		if client.Get("cutdate.status."+today).Val() != "success" {
			// 日切失败就不做日终检查
			return
		}
		if cd.CheckCutEndStatus() {
			//Redis: set "20190214cutend" "ok"
			key := today + "cutend.finsh"
			client := RedisClient(
				os.Getenv("CACHE_REDIS_HOST"),
				os.Getenv("CACHE_REDIS_PASSWORD"),
				os.Getenv("CACHE_REDIS_DB"))
			val := client.Get(key)
			if val.Val() == "" {
				// 断开会计核算主从
				Info.Println("-------- [*]STEP: pause loan ACT MySQL master-slave -------- ")
				SendNotify("general", "ODS_SUPPORT:断开主从"+NowFormatDate("20060102"), "[BEGIN]断开会计核算主从{ACT}")
				//BatchHandleDbSlave("act", "stop")

				//通知下游系统抽数
				SendNotify("general", "ODS_SUPPORT:通知下游系统抽数"+NowFormatDate("20060102"), "[BEGIN]通知下游系统抽数{python,ods,bigdata}")
				Info.Println("-------- [*]STEP: notify downstream service[Python,ODS,Bigdata] to start job -------- ")
				PushCutBatchMsg()
				client.Set(key, "success", 0)
			} else if val.Val() == "success" {
				Info.Println("cutdate and cutend finish aleady push status to other job platform, ignore")
			}
		}
	})

	// 0-1点每分钟检查一次各个系统的状态后恢复主从, 恢复告警，这个要避免周期内的重复操作
	//spec04 := "*/60 3-59 0-1 * * ?"
	spec04 := os.Getenv("CRON_RESTORE_SRV")
	c.AddFunc(spec04, func() {
		if GetTodayOdsAllStatusFormRedis() {
			key := NowFormatDate("20060102") + "restore"
			client := RedisClient(os.Getenv("CACHE_REDIS_HOST"),
				os.Getenv("CACHE_REDIS_PASSWORD"),
				os.Getenv("CACHE_REDIS_DB"))
			val := client.Get(key)
			if val.Val() == "" {
				// 恢复会计核算主从状态
				SendNotify("general", "ODS_SUPPORT:恢复主从"+NowFormatDate("20060102"), "[BEGIN]开始恢复会计核算主从{ACT}")
				Info.Println("-------- [*]STEP: restore loan ACT MySQL master-slave -------- ")
				BatchHandleDbSlave("act", "start")
				// 恢复信贷核心、信贷管理主从状态
				SendNotify("general", "ODS_SUPPORT:恢复主从"+NowFormatDate("20060102"), "[BEGIN]开始恢复信贷核心{CAS}、信贷管理主从{CMS}")
				Info.Println("-------- [*]STEP: restore loan CAS and CMS MySQL master-slave -------- ")
				BatchHandleDbSlave("cascms", "start")
				// 恢复告警
				api := getAPI()
				SendNotify("general", "ODS_SUPPORT:告警恢复"+NowFormatDate("20060102"), "[BEGIN]解除主从异常告警屏蔽")
				Info.Println("-------- [*]STEP: restore zabbix alert about MySQL master-slave -------- ")
				api.RestoreAlert()
				client.Set(key, "ok", 0)
			} else if val.Val() == "ok" {
				Info.Println("slave status and alert already restore, ignore")
			}
		}
	})
	c.Start()
	select {}
}
