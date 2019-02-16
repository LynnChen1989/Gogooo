package main

import (
	"github.com/robfig/cron"
	"os"
	"time"
)

func OdsTimer() {
	// 停服定时器，于每天23:59:50定时发布停服公告
	c := cron.New()
	// 23点50分关闭服务入口
	spec01 := "0 50 23 * * ?"
	//DEBUG
	//spec01 := "0 51 21 * * ?"
	c.AddFunc(spec01, func() {
		Info.Println("[*]STEP: start close CROS service entry point")
		PauseSrv()
		status := PauseSrv()
		if status == "success" {
			SendNotify("general", "ODS抽数:关闭服务入口SUCCESS", "关闭服务入口成功")
		} else if status == "failed" {
			SendNotify("fatal", "ODS抽数:关闭服务入口FAILURE", "关闭服务入口失败")
		}
	})

	// 0点2分检查日切状态，日切成功后断开核心和管理主从关系
	spec02 := "0 2 0 * * ?"
	c.AddFunc(spec02, func() {
		cd := &CutDate{}
		if cd.CheckCutDateStatus() {
			//屏蔽告警
			api := getAPI()
			Info.Println("[*]STEP: prevent zabbix alert about MySQL master-slave")
			api.PauseAlert()
			SendNotify("general", "ODS抽数:开始屏蔽主从告警", "开始屏蔽主从告警")

			//断开核心，管理主从
			Info.Println("[*]STEP: pause loan CAS and CMS MySQL master-slave")
			SendNotify("general", "ODS抽数:开始断开主从", "开始断开信贷核心、信贷管理主从")
			BatchHandleDbSlave("cascms", "stop")

			//打开服务入口
			Info.Println("[*]STEP: pause loan CAS and CMS MySQL master-slave")
			SendNotify("general", "ODS抽数:开始恢复服务入口", "开始恢复服务入口")
			status := RestoreSrv()
			if status == "success" {
				SendNotify("general", "ODS抽数:开启服务入口SUCCESS", "开启服务入口成功")
			} else if status == "failed" {
				SendNotify("fatal", "ODS抽数:开启服务入口FAILURE", "开启服务入口失败")
			}
		}
	})

	// 每隔30秒检查一次日终状态, 日终完成断开核算主从，后通知下游服务, 这个要避免周期内的重复推送
	spec03 := "*/30 3-59 0 * * ?"
	c.AddFunc(spec03, func() {
		cd := &CutDate{}
		if cd.CheckCutEndStatus() {
			//Redis: set "20190214cutend" "ok"
			currentTime := time.Now()
			key := currentTime.Format("20060102") + "cutend"
			client := RedisClient(os.Getenv("CACHE_REDIS_HOST"),
				os.Getenv("CACHE_REDIS_PASSWORD"),
				os.Getenv("CACHE_REDIS_DB"))
			val := client.Get(key)
			if val.Val() == "" {
				// 断开会计核算主从
				Info.Println("[*]STEP: pause loan ACT MySQL master-slave")
				SendNotify("general", "ODS抽数:开始断开主从", "开始断开会计核算主从")
				BatchHandleDbSlave("act", "stop")

				//通知下游系统抽数
				SendNotify("general", "ODS抽数:开始通知下游系统抽数", "开始通知下游系统抽数")
				Info.Println("[*]STEP: notify downstream service[Python,ODS,Bigdata] to start job")
				PushCutBatchMsg()
				client.Set(key, "ok", 0)
			} else if val.Val() == "ok" {
				Info.Println("cutdate and cutend finish aleady push status to other job platform, ignore")
			}
		}
	})

	// 0-1点每分钟检查一次各个系统的状态后恢复主从, 恢复告警，这个要避免周期内的重复操作
	spec04 := "*/60 3-59 0-1 * * ?"
	c.AddFunc(spec04, func() {
		if GetTodayOdsAllStatusFormRedis() {
			currentTime := time.Now()
			key := currentTime.Format("20060102") + "restore"
			client := RedisClient(os.Getenv("CACHE_REDIS_HOST"),
				os.Getenv("CACHE_REDIS_PASSWORD"),
				os.Getenv("CACHE_REDIS_DB"))
			val := client.Get(key)
			if val.Val() == "" {
				// 恢复会计核算主从状态
				SendNotify("general", "ODS抽数:开始恢复主从", "开始恢复会计核算主从")
				Info.Println("[*]STEP: restore loan ACT MySQL master-slave")
				BatchHandleDbSlave("act", "start")
				// 恢复信贷核心、信贷管理主从状态
				SendNotify("general", "ODS抽数:开始恢复主从", "开始恢复信贷核心、信贷管理主从")
				Info.Println("[*]STEP: restore loan CAS and CMS MySQL master-slave")
				BatchHandleDbSlave("cascms", "start")
				// 恢复告警
				api := getAPI()
				SendNotify("general", "ODS抽数:告警恢复", "开始恢复主从异常告警")
				Info.Println("[*]STEP: restore zabbix alert about MySQL master-slave")
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
