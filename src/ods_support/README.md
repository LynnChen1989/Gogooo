# ODS_SUPPORT


## 环境变量

| 变量名 | 变量值 | 说明|
| ------ | ------ | -----|
| YHB_REDIS_HOST | 192.168.0.53:6379 | REDIS主机端口 
| CACHE_REDIS_HOST | 192.168.0.53:6379 | REDIS主机端口 
| YHB_REDIS_PASSWORD | 123123 | REDIS 密码
| CACHE_REDIS_PASSWORD | 123123 | REDIS 密码
| YHB_REDIS_DB | 15 | REDIS 数据库
| CACHE_REDIS_DB | 0 | REDIS 数据库
| MESSAGE_URL | http://122.152.209.199:2046/api/v1/message/ | 消息系统URL 
| MESSAGE_CID | 40001 | 消息系统CID
| ZABBIX_URL | https://monitor.xwfintech.com/api_jsonrpc.php | 监控API地址
| ZABBIX_USER | snakechen | 监控API用户名
| ZABBIX_PASSWORD | 123 | 监控API密码
| DAP_DB_INFO | root:fhl3mjsdwj@tcp(172.16.1.18:5001)/db_dap | 序列服务MySQL  
| CAS_DB_INFO | root:fhl3mjsdwj@tcp(172.16.1.18:5001)/cas | 信贷核心MySQL
| ACT_DB_INFO | root:fhl3mjsdwj@tcp(172.16.1.18:5000)/db_act | 会计核算MySQL
| MQ_URI | amqp://snake:snake@127.0.0.1:5672/snakehost | 消息队列地址
| CAS_CMS_SLAVE_LIST | amqp://snake:snake@127.0.0.1:5672/snakehost##amqp://snake:snake@127.0.0.1:5672/snakehost| 从库的列表，用两个井号隔开
| ACT_SLAVE_LIST | amqp://snake:snake@127.0.0.1:5672/snakehost##amqp://snake:snake@127.0.0.1:5672/snakehost| 从库的列表，用两个井号隔开
| CRON_STOP_SRV | 0 0 17 * * ?| 停服cron表达式
| CRON_CUT_DATE | 0 0 17 * * ?| 日切cron表达式
| CRON_CUT_END | 0 0 17 * * ?| 日终cron表达式
| CRON_RESTORE_SRV | 0 0 17 * * ?| 启服cron表达式


## 系统整体设计逻辑

```shell

1. 于每日23:50:00秒关闭服务入口
2. 00:02:00检查日切状态，日切完成后，屏蔽主从告警，断开主从
3. 00:03:00~00:59:59轮询日终状态，日终完成后，通知下游服务（python,ods,bigdata）开始作业
4. 00:03:00~01:59:59获取下游服务作业状态，下游服务作业完成后，恢复主从，恢复告警
```


## 调试API

#### 0、空接口

```
/ 

```
#### 1、关闭服务入口

```
GET /stopsrv
```

#### 2、开启服务入口

```
GET /startsrv
```

#### 3、通知下游服务

```
GET /push
```


#### 4、获取日切、日终状态

```
GET /cut
```

####  5、主从批量操作

```
GET /slave?system=act&opt=stop

*system可选参数*
 + act
 + cascms
 
*opt可选参数*
 + start
 + stop 

```

#### 6、告警操作

```
GET /alert?opt=stop
 
*opt可选参数*
 + start
 + stop 

```
