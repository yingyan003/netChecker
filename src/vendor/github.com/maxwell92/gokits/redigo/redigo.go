package redigo

/*
import (
	config "app/backend/common/yce/config"
	redis "github.com/garyburd/redigo/redis"
	mylog "github.com/maxwell92/gokits/log"
	"sync"
)

var once sync.Once
var instance *redis.Pool
var log = mylog.Log

func NewRedisClient() *redis.Pool {
	once.Do(func() {
		instance = &redis.Pool{
			MaxIdle:     config.MAX_IDLE,
			MaxActive:   config.MAX_ACTIVE,
			IdleTimeout: config.IDEL_TIMEOUT,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", config.Instance().GetRedisEndpoint())
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		} // instance
	})

	log.Infof("RedisClient Open Success: host=%s", config.REDIS_HOST)
	return instance
}
*/
