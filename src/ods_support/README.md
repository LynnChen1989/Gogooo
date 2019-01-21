# ODS_SUPPORT


## 环境变量

| 变量名 | 变量值 | 说明|
| ------ | ------ | -----|
| REDIS_URI | 192.168.0.53:6379 | REDIS主机端口 
| REDIS_PASSWORD | 123123 | REDIS 密码
| REDIS_DB | 15 | REDIS 数据库
| MESSAGE_URL | http://122.152.209.199:2046/api/v1/message/ | 消息系统URL 
| MESSAGE_CID | 40001 | 消息系统CID





## 系统设计

```shell

1. 于每日23:59:00秒停服
2. 检查日切和跑批的状态
3. 断开主从（屏蔽主从告警）
4. 大数据平台从SLAVE同步数据
5. 开启服务入口 (解除主从告警屏蔽)
6. 等待ODS进行抽数
7. ODS抽数完成后恢复主从
```

