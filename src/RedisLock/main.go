package RedisLock

import (
	"log"
	"time"
	"os"
)


/*
 *	重连等待时间，最多等待maxSleepTime秒
 */
func reconnectSleep(sleeptime *time.Duration,maxSleepTime *time.Duration) {
	crontype:=time.Second
	//避免0，如果是0使用毫秒计算；最长不超过MAX_SLEEP_TIME
	if *sleeptime==0 {
		*sleeptime=300
		crontype=time.Millisecond
	}else if *sleeptime>*maxSleepTime {
		*sleeptime=*maxSleepTime
	}

	log.Println("waiting for seconds......")
	time.Sleep(*sleeptime*crontype)
}

func Lock(lockkey *string,addr *string,password *string,db *int,expiration *time.Duration,reconnectTimes *int,maxSleepTime *time.Duration)  {
	redislock:=NewRedisLock(addr,password,db,expiration)
	defer redislock.Close()
	if redislock.Lock(lockkey)==true {
		log.Println("success lock "+*lockkey)
		os.Exit(0)
	}else{
		//如果是因为连接断开原因导致的加锁失败，重试N次
		if redislock.IsConnected()==false {
			sleeptime:=*expiration/time.Duration(*reconnectTimes)
			for i:=0;i<*reconnectTimes;i++ {
				log.Println("redis connect error,try to reconnect")
				reconnectSleep(&sleeptime,maxSleepTime)

				//如果恢复连接，再尝试加一次锁
				if redislock.IsConnected()==true {
					if redislock.Lock(lockkey)==true {
						log.Println("reconnect success,success lock " + *lockkey)
						os.Exit(0)
					}else{
						log.Fatalln("reconnect success,false lock " + *lockkey)
					}
				}
			}
		}
		//异常并退出,code 1
		log.Fatalln("false lock "+*lockkey)
	}
}