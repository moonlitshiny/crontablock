package RedisLock

import (
	"gopkg.in/redis.v4"
	"time"
)

type RedisLock struct {
	client *redis.Client
	expiration *time.Duration
}

func NewRedisLock(addr *string,password *string,db *int,expiration *time.Duration) *RedisLock {
	client := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: *password, // no password set
		DB:       *db,  // use default DB
	})
	redislock:=&RedisLock{
		client,
		expiration,
	}
	return redislock
}

func (this *RedisLock) Close() error {
	return this.client.Close()
}

func (this *RedisLock) Unlock(lockkey *string) bool {
	delres:=this.client.Del(*lockkey)
	if delres.Val()!=1 {
		return false
	}
	return true
}

func (this *RedisLock) Lock(lockkey *string) bool {
	res:=this.client.SetNX(*lockkey,1,*this.expiration*time.Second)
	return res.Val()
}

func (this *RedisLock) IsConnected() bool {
	res:=this.client.Ping()
	if res.Val()!="PONG" {
		return false
	}
	return true
}
