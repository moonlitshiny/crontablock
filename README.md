# CrontabLock
## 简介：
crontablock是一个非侵入式的crontab任务锁，解决多台设备同时跑同一个crontab任务，底层使用redis原子操作做锁操作
该工具在压测环境上百台设备压测过，现已投入运营环境使用

## 参数：
## Usage of ./crontablock:
   ```bash
      -k string
            lockkey(必选):唯一的锁key，如 test_lockkey
      -a string
            addr:redis地址 (default "127.0.0.1:6379")
      -d int
            db:redis库分区，不使用默认分区 (default 1)
      -e duration
            expiration:lock时间，默认为1个crontab周期之内，如果不同服务的时间不同步建议加大该值（秒） (default 55ns)
      -p string
            password:redis密码 (default "")
      -r int
            reconnect times:当链接断开时会尝试重连几次，如果由于连接断开导致加锁失败，重试次数 (default 3)
      -s duration
            max Sleep Time:当链接断开时会尝试重连几次，如果由于连接断开导致加锁失败，每次重试等待时间（秒） (default 3ns)
   ```

## 使用例子：
### 命令行测试：同一个lockkey默认会锁55秒（一个crontab周期内）
    [root@zy_auto ~]# ./crontablock -k mykey
    2018/01/11 16:16:40 success lock mykey
    [root@zy_auto ~]# ./crontablock -k mykey
    2018/01/11 16:16:42 false lock mykey


### crontab使用：
    格式为：crontablock -k $mykey && 业务自身的任务

    说明：当crontablock拿不到锁时会返回错误，业务自身的任务不会被执行，如:
    * * * * * /root/crontablock -k lockkey >/root/redislock.log 2>&1 && /root/test.sh >/root/test.log 2>&1 &

    机器A（获得锁）的结果：
        [root@zy_auto ~]# cat redislock.log
        2018/01/11 16:10:01 success lock lockkey
        [root@zy_auto ~]# cat test.log
        Thu Jan 11 16:10:01 CST 2018
        Hello world

    机器B（未抢到锁）的结果：
        [root@node1 ~]# cat redislock.log
        2018/01/11 16:10:01 false lock lockkey
        [root@node1 ~]# cat test.log
        [root@node1 ~]#
