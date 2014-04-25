package shardredis

import (
	"github.com/latermoon/redigo/redis"
	"sync"
	"time"
)

/*
shardredis.Load(fd)
cluster := shardredis.Get("redis-profile-b")
rd := cluster.Get("100422")
reply, err := rd.Do("SET", "name", "latermoon")
rd.Close()
*/

var mu sync.Mutex
var pools = map[string]*redis.Pool{}

func RedisPool(host string) (pool *redis.Pool) {
	var ok bool
	if pool, ok = pools[host]; !ok {
		mu.Lock()
		defer mu.Unlock()
		if pool, ok = pools[host]; !ok {
			pool = &redis.Pool{
				MaxIdle:     100,
				IdleTimeout: 240 * time.Second,
				Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host) },
			}
			pools[host] = pool
		}
	}
	return
}
