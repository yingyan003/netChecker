package cache

import (
	"fmt"
	"strings"
	"sync"

	"github.com/garyburd/redigo/redis"
	mylog "github.com/maxwell92/gokits/log"
)

var log = mylog.Log

type RedisCache struct {
	pool *redis.Pool
}

var instance *RedisCache
var once sync.Once

func RedisCacheInstance() *RedisCache {
	return instance
}

func NewRedisCache(p *redis.Pool) *RedisCache {
	once.Do(func() {
		instance = &RedisCache{
			pool: p,
		}
	})
	log.Tracef("RedisCache Open Success")
	return instance
}

func (rc *RedisCache) Get(key string) string {
	conn := rc.pool.Get()
	// 用于压测Redis链接数
	// log.Errorf("RedisCache: active=%d", rc.pool.ActiveCount())
	if conn == nil {
		log.Errorf("RedisCache Get Connection Error: Nil")
		return ""
	}
	defer conn.Close()

	ok, err := rc.Exist(key)
	if err != nil {
		log.Errorf("RedisCache Get Exist key %s Error: error=%s", key, err)
		return ""
	}
	if !ok {
		log.Warnf("RedisCache Get key %s Not Exist", key)
		return ""
	}

	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Errorf("RedisCache Get Key Error: error=%s", err)
		return ""
	}
	return result
}

func (rc *RedisCache) MExist(keys []string) (bool, error) {
	var ok bool
	var err error
	for _, key := range keys {
		ok, err = rc.Exist(key)
		if err != nil {
			log.Errorf("RedisCache MExist Error: error=%s")
			return ok, err
		}
	}
	return ok, nil
}

func (rc *RedisCache) MGet(keys []string) []string {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Get Connection Error: Nil")
		return nil
	}
	defer conn.Close()

	_, err := rc.MExist(keys)
	if err != nil {
		log.Errorf("RedisCache Get MExist Error: error=%s", err)
		return nil
	}

	results, err := redis.Strings(conn.Do("MGET", keys))
	if err != nil {
		log.Errorf("RedisCache Get Key Error: error=%s", err)
		return nil
	}
	return results
}

func (rc *RedisCache) Set(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Set Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	result, err := conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisCache Set Key Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache Set Key Success: result=%v", result)
	return true, nil
}

func (rc *RedisCache) SetWithExpire(key, value, expire string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache SetWithExpire Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	result1, err1 := conn.Do("SET", key, value)
	if err1 != nil {
		log.Errorf("RedisCache SetWithExpire Set Key Error: error=%s", err1)
		return false, err1
	}
	result2, err2 := conn.Do("EXPIRE", key, expire)
	if err2 != nil {
		log.Errorf("RedisCache SetWithExpire Expire Key Error: error=%s", err2)
		return false, err2
	}

	log.Tracef("RedisCache SetWithExpire Key Success: set result=%v, expire result=%v", result1, result2)
	return true, nil
}

func (rc *RedisCache) Mget(keys []string) interface{} {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache MGet Connection Error: Nil")
		return nil
	}
	defer conn.Close()

	for _, k := range keys {
		if exists, err := rc.Exist(k); !exists || err != nil {
			log.Errorf("RedisCache MGet Non-exist key Error")
			return nil
		}
	}

	// TODO: it should be keys...
	result, err := redis.Strings(conn.Do("MGET", keys[0]))
	if err != nil {
		log.Errorf("RedisCache MGet Key Error: error=%s", err)
		return nil
	}
	return result
}

func (rc *RedisCache) Delete(key string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Delete Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	result, err := conn.Do("DEL", key)
	if err != nil {
		log.Errorf("RedisCache Delete Key Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache Delete Key Success: result=%v", result)
	return true, nil
}

func (rc *RedisCache) Lpush(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LPush Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		log.Errorf("RedisCache LPush Key Error: error=%s", err)
		return false, err
	}
	return true, nil
}

func (rc *RedisCache) Lrange(key, start, end string) ([]string, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LRange Connection Error: Nil")
		return nil, nil
	}
	defer conn.Close()

	result, err := redis.Strings(conn.Do("LRANGE", key, start, end))
	if err != nil {
		log.Errorf("RedisCache LRange Key Error: error=%s", err)
		return nil, err
	}
	return result, nil
}

func (rc *RedisCache) Llen(key string) int32 {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LLen Connection Error: Nil")
		return 0
	}
	defer conn.Close()

	result, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		log.Errorf("RedisCache LLen Key Error: error=%s", err)
		return 0
	}
	return int32(result)
}

func (rc *RedisCache) Lrem(key, count, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LLen Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("LREM", key, count, value)
	if err != nil {
		log.Errorf("RedisCache LLen Key Error: error=%s", err)
		return false, err
	}
	return true, nil
}

func (rc *RedisCache) Exist(key string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Exist Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		log.Errorf("RedisCache Exist Key Error: error=%s", err)
		return false, err
	}

	if results == 0 {
		return false, nil
	}
	return true, nil
}

func (rc *RedisCache) Sadd(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Sadd Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("SADD", key, value))
	if err != nil {
		log.Errorf("RedisCache Sadd Key Error: error=%s", err)
		return false, err
	}

	if results == 0 {
		log.Warnf("RedisCache Sadd Key Warning: zero reply")
		return false, nil
	}
	return true, nil
}

func (rc *RedisCache) Smember(key string) ([]string, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Smember Connection Error: Nil")
		return nil, nil
	}
	defer conn.Close()

	results, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Errorf("RedisCache Smember Key Error: error=%s", err)
		return nil, err
	}

	return results, nil
}

func (rc *RedisCache) Scard(key string) (int32, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Scard Connection Error: Nil")
		return 0, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("SCARD", key))
	if err != nil {
		log.Errorf("RedisCache Scard Key Error: error=%s", err)
		return 0, err
	}

	return int32(results), nil
}

func (rc *RedisCache) Srem(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Srem Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := conn.Do("SREM", key, value)
	if err != nil {
		log.Errorf("RedisCache Srem Key Error: error=%s", err)
		return false, err
	}

	if results == 0 {
		return false, nil
	}
	return true, nil
}

func (rc *RedisCache) LockThenGet(key string) string {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LockThenGet Conneciton Error: Nil")
		return ""
	}
	defer conn.Close()

	lock := key + "-lock"
	results, err := conn.Do("WATCH", lock)
	if err != nil {
		log.Errorf("RedisCache LockThenGet Watch Error: error=%s", err)
		return ""
	}

	results, err = conn.Do("MULTI")
	if err != nil {
		log.Errorf("RedisCache LockThenGet Error: error=%s", err)
		return ""
	}
	log.Tracef("RedisCache LockThenGet Multi Results: results=%s", results)

	results, err = conn.Do("SET", lock, "1")
	if err != nil {
		log.Errorf("RedisCache LockThenGet SET Error: error=%s", err)
		return ""
	}
	log.Tracef("RedisCache LockThenGet Set Results: results=%s", results)

	results, err = conn.Do("GET", key)
	if err != nil {
		log.Errorf("RedisCache LockThenGet GET Error: error=%s, key=%s", err, key)
		return ""
	}
	log.Tracef("RedisCache LockThenGet Get Results: results=%s", results)

	results, err = conn.Do("EXEC")
	if err != nil {
		log.Errorf("RedisCache LogThenGet Exec Error: error=%s", err)
		return ""
	}
	log.Tracef("RedisCache LockThenGet Exec Results: results=%s", results)
	raw := results.([]interface{})
	log.Tracef("RedisCache LockThenGet raw: raw[0]=%s, raw[1]=%s", raw[0], raw[1])
	var r0, r1 string
	if raw[0] != nil {
		r0 = fmt.Sprintf("%s", raw[0])
		if strings.EqualFold(r0, "OK") && strings.EqualFold(string(r1), "") {
			r1 = fmt.Sprintf("%s", raw[1])
			log.Tracef("RedisCache LockThenGet Return: return=%s", r1)
			return r1
		}
	}
	/*
		if raw[1] != nil {
			r1 = fmt.Sprintf("%s", raw[1])
			if strings.EqualFold(r1, "OK") {
				log.Tracef("RedisCache LockThenGet Return: return=%s", r0)
				return true, r0
			}
		}
	*/

	return ""
}

func (rc *RedisCache) SetThenUnlock(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache SetThenUnlock Conneciton Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := conn.Do("MULTI")
	if err != nil {
		log.Errorf("RedisCache SetThenUnlock Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache SetThenUnlock Multi Results: results=%s", results)

	results, err = conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisCache SetThenUnlock Set Error: error=%s, key=%s", err, key)
		return false, err
	}
	log.Tracef("RedisCache SetThenUnlock Set Results: results=%s", results)

	lock := key + "-lock"
	// results, err = conn.Do("SET", lock, "0")
	results, err = conn.Do("DEL", lock)
	if err != nil {
		log.Errorf("RedisCache SetThenUnlock SET Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache SetThenUnlock Set Results: results=%s", results)

	results, err = conn.Do("EXEC")
	if err != nil {
		log.Errorf("RedisCache SetThenUnlock Exec Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache SetThenUnlock Exec Results: results=%s", results)
	return true, nil
}

func (rc *RedisCache) Txn(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Txn Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := conn.Do("MULTI")
	if err != nil {
		log.Errorf("RedisCache Txn Error: error=%s", err)
		return false, err
	}
	log.Infof("RedisCache Txn Results: results=%s", results)

	results, err = conn.Do("GET", key)
	if err != nil {
		log.Errorf("RedisCache Txn Error: error=%s", err)
		return false, err
	}
	log.Infof("RedisCache Txn Results: results=%s", results)

	results, err = conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisCache Txn Error: error=%s", err)
		return false, err
	}
	log.Infof("RedisCache Txn Results: results=%s", results)

	results, err = conn.Do("EXEC")
	if err != nil {
		log.Errorf("RedisCache Txn Error: error=%s", err)
		return false, err
	}

	log.Infof("RedisCache Txn Results: results=%s", results)
	return true, nil
}

/*
func (rc *RedisCache) Transaction(key, value string, fns func(string, string))(bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Transaction Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()
	sync.Mutex.Lock()
	_, err := conn.Do("MULTI")
	if err != nil {
		log.Errorf("RedisCache Transaction Key Error: error=%s", err)
		return false, err
	}

	for _, f := range fns {
		f(key, value)
	}

	_, err1 := conn.Do("EXEC")
	if err1 != nil {
		log.Errorf("RedisCache Transaction Key Error: error=%s", err)
		return false, err1
	}
	sync.Mutex.Unlock()
	return true, nil
}


func (rc *RedisCache) Search(keyword string) (*map[string]string, error) {return nil, nil}
func (rc *RedisCache) Watch(name string) error {return nil}
func (rc *RedisCache) Create(db string) error {return nil}
func (rc *RedisCache) Index(columns []string) error {return nil}
*/
