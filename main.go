package main

import (
	"./src/RedisLock"
	"flag"
	"log"
)


func main()  {
	lockkey := flag.String("k", "", "lockkey(必选):唯一的锁key，如 test_lockkey")
	addr := flag.String("a", "127.0.0.1:6379", "addr:redis地址")
	password := flag.String("p", "", "password:redis密码")
	db := flag.Int("d", 1, "db:redis库分区，不使用默认分区")
	expiration := flag.Duration("e", 55, "expiration:lock时间，默认为1个crontab周期之内，如果不同服务的时间不同步建议加大该值（秒）")
	reconnectTimes := flag.Int("r", 3, "reconnect times:如果由于连接断开导致加锁失败，重试次数")
	maxSleepTime := flag.Duration("s", 3, "max Sleep Time:如果由于连接断开导致加锁失败，每次重试等待时间（秒）")

	flag.Parse() //解析输入的参数
	if *lockkey=="" {
		log.Fatalln("参数错误，需要 lockkey 参数，格式为：-k $lockkey ；详情见 --help")
	}
	RedisLock.Lock(lockkey,addr,password,db,expiration,reconnectTimes,maxSleepTime)
}
