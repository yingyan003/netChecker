package utils

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"common/constant"
)

type RedisClient struct {
	pool *redis.Pool
	//subConn *redis.PubSubConn
	//pubConn *redis.Conn
}

var Redis *RedisClient
var once sync.Once

//使用redis前需要实例化redisClient
func NewRedis(maxIdle int) {
	once.Do(func() {
		log.Infof("enter once.do newRedis")
		Redis = new(RedisClient)
		Redis.pool = newPool(maxIdle)
	})
}

func newPool(maxIdle int) *redis.Pool {
	host := LoadEnvVar(constant.ENV_REDISHOST, constant.REDISHOST)
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   constant.MAX_ACTIVE,
		IdleTimeout: constant.IDLE_TIMEOUT,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Errorf("RedisClient Dial host failed: host=%s, err=%v", host, err)
				return nil, err
			}
			return c, nil
		},
	}
	return pool
}

func (r *RedisClient) Get(key string) string {
	//从连接池pool中获取一个可用的空闲连接，实际执行的是redis.Pool{Dial: func}中的func函数
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient Get failed: the connection get from pool is nil.")
		return ""
	}
	defer conn.Close()

	ok, err := r.Exists(key)
	if err != nil {
		log.Errorf("RedisClient Get failed: key=%s Exists error. err=%s", key, err)
		return ""
	}

	if !ok {
		log.Warnf("RedisClient Get: key=%s Not exist", key)
		return ""
	}

	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Errorf("RedisClient Get failed: Do error. err=%s", err)
		return ""
	}

	return result
}

func (r *RedisClient) Set(key, value string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient Set failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisClient Set failed: Do error. err=%s", err)
		return false, err
	}

	return true, nil
}

func (r *RedisClient) SetWithExpire(key, value, expire string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient SetWithExpire failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisClient SetWithExpire failed: Do Set error. error=%s", err)
		return false, err
	}

	_, err = conn.Do("EXPIRE", key, expire)
	if err != nil {
		log.Errorf("RedisClient SetWithExpire failed: Do EXPIRE error. error=%s", err)
		return false, err
	}

	return true, nil
}

func (r *RedisClient) Delete(key string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient Delete failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Errorf("RedisClient Delete failed: Do error. err=%s", err)
		return false, err
	}

	return true, nil
}

func (r *RedisClient) Exists(key string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient Exist failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	result, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		log.Errorf("RedisClient Exist failed: Do error. err=%s", err)
		return false, err
	}

	if result == 0 {
		return false, nil
	}

	return true, nil
}

func (r *RedisClient) Publish(channel string, message []byte) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Errorf("RedisClient Publish failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, message)
	if err != nil {
		log.Errorf("RedisClient Publish failed: Do error. err=%s", err)
		return false, err
	}

	return true, nil
}

func (r *RedisClient) getSubConn() *redis.PubSubConn {
	return &redis.PubSubConn{Conn: Redis.pool.Get()}
}

//连接池的链接数需要修改
func (r *RedisClient) GetSubConn(channels ...string) *redis.PubSubConn {
	//获取可用的redis连接
	subConn := r.getSubConn()

	//订阅频道，订阅成功与否与频道上是否有数据传输无关。当订阅的频道无数据传输时，subConn.Receibe()阻塞
	if err := subConn.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		log.Errorf("GetSubConn: Subscribe channels error. channels=%s, subConn=%v, err=%s", channels, subConn, err)
		subConn.Close()
		return nil
	}
	//todo 成功时subConn不关闭
	log.Infof("GetSubConn: Subscriben channels success. channels=%s", channels)
	return subConn
}

//如果获取subConn失败，重试直到成功.
//注意，重试会尝试开启新的链接，而旧的连接并未关闭，导致连接失败
func (r *RedisClient) RetrySubConn(channel string) *redis.PubSubConn {
	retry := 0
	for {
		retry++
		subConn := r.GetSubConn(channel)
		if subConn != nil {
			log.Infof("CheckSubConn：retry getSubConn success. retry time=%d", retry)
			return subConn
		}
	}
}

func (r *RedisClient) ReceiveSubMessage(subConn *redis.PubSubConn) []byte {
	receive := subConn.Receive()
	switch result := receive.(type) {
	case error:
		log.Errorf("receiveSubMessage: type is error. err=%s", result)
	case redis.Subscription:
		log.Infof("receiveSubMessage: type is Subscription. count=%d", result.Count)
	case redis.Message:
		log.Infof("receiveSubMessage: type is message")
		return result.Data
	}
	return nil
}
