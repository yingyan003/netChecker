package cache

import (
	config "app/backend/common/yce/config"
	"fmt"
	"testing"
)

func Benchmark_MGet(b *testing.B) {
	config.Instance().RedisHost = "10.151.160.11:32379"
	cache := NewRedisCache()

	keys := []string{
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
		"216:dev:hess",
	}

	for i := 0; i < b.N; i++ {
		cache.MGet(keys)
	}
}

func Benchmark_Get(b *testing.B) {
	config.Instance().RedisHost = "10.151.160.11:32379"
	cache := NewRedisCache()

	key := "216:dev:hess"

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			cache.Get(key)
		}
	}
}
func TestAll(t *testing.T) {
	config.Instance().RedisHost = "10.151.160.11:32379"
	cache := NewRedisCache()

	key := "1:test:nginx"
	value := "{\"userName\":\"liyao.miao\", \"deploymentName\":\"nginx\"}"
	ok, err := cache.Exist("1:test:nginx")
	if err != nil {
		fmt.Printf("Exists Error: error=%s\n", err)
	}
	if !ok {
		fmt.Printf("key %s Not Exists\n", key)
	} else {
		fmt.Printf("key %s exists\n", key)
	}

	ok, err = cache.Set(key, value)
	if err != nil || !ok {
		fmt.Printf("Set Error: error=%s\n", err)
		return
	}
	fmt.Printf("Key %s set success\n", key)

	results := cache.Get(key)
	if results == "" {
		fmt.Printf("Get key %s Error\n", key)
		return
	}
	fmt.Printf("Key %s Get success value %s\n", key, results)

	tmpKey := "test:1"
	tmpValue := "nginx"

	l, err := cache.Scard(tmpKey)
	if err != nil {
		fmt.Printf("Key %s Scard Error: error=%s\n", tmpKey, err)
		return
	}
	fmt.Printf("Before %d len\n", l)

	ok, err = cache.Sadd(tmpKey, tmpValue)
	if err != nil {
		fmt.Printf("Key %s Sadd Error: error=%s\n", tmpKey, err)
		return
	}

	if !ok {
		fmt.Printf("Key %s Sadd Failed\n", tmpKey)
	} else {
		fmt.Printf("Key %s Sadd Success\n", tmpKey)
	}

	l, err = cache.Scard(tmpKey)
	if err != nil {
		fmt.Printf("Key %s Scard Error: error=%s\n", tmpKey, err)
		return
	}
	fmt.Printf("After %d len\n", l)

	result, err := cache.Smember(tmpKey)
	if err != nil {
		fmt.Printf("Key %s Smember Error: error=%s\n", tmpKey, err)
		return
	}
	fmt.Printf("Key %s Smember results %v\n", tmpKey, result)

	/*
		l := cache.Llen(tmpKey)
		fmt.Printf("Before %d len\n", l)

		ok, err = cache.Lpush(tmpKey, tmpValue)
		if err != nil || !ok {
			fmt.Printf("Lpush Error: error=%s\n", err)
			return
		}
		fmt.Printf("Key %s lpush success value %s\n", tmpKey, tmpValue)

		l = cache.Llen(tmpKey)
		fmt.Printf("After %d len\n", l)

		start := "0"
		end := "1"
		result, err := cache.Lrange(tmpKey, start, end)
		if err != nil {
			fmt.Printf("Lrange Error: error=%s\n", err)
			return
		}
		fmt.Printf("Key %s Lrange success value %v\n", tmpKey, result)
	*/
	ok, err = cache.Delete(key)
	if err != nil || !ok {
		fmt.Printf("Delete Error: error=%s\n", err)
		return
	}
	fmt.Printf("Delete Key %s success\n", key)

	ok, err = cache.Exist("1:test:nginx")
	if err != nil {
		fmt.Printf("Exists Error: error=%s\n", err)
		return
	}
	if !ok {
		fmt.Printf("key %s not exists\n", key)

	} else {
		fmt.Printf("key %s exists\n", key)
	}

	ok, err = cache.Srem(tmpKey, tmpValue)
	if err != nil {
		fmt.Printf("Key %s Srem Error: error=%s\n", tmpKey, err)
		return
	}

	if !ok {
		fmt.Printf("Key %s Srem Failed\n", tmpKey)
	} else {
		fmt.Printf("Key %s Srem Succeed\n", tmpKey)
	}

	l, err = cache.Scard(tmpKey)
	if err != nil {
		fmt.Printf("Key %s Scard Error: error=%s\n", tmpKey, err)
		return
	}
	fmt.Printf("After %d len\n", l)

	result, err = cache.Smember(tmpKey)
	if err != nil {
		fmt.Printf("Key %s Smember Error: error=%s\n", tmpKey, err)
		return
	}
	fmt.Printf("Key %s Smember results %v\n", tmpKey, result)

	cache.Txn("test-key", "test-value-1")
	cache.LockThenGet("test-key")
	cache.SetThenUnlock("test-key", "test-value-3")

	/*
		ok, err = cache.Lrem(tmpKey, "0", "nginx")
		if err != nil || !ok {
			fmt.Printf("Lrem Error: error=%s\n", err)
			return
		}
		fmt.Printf("Lrem key %s success\n", tmpKey)

		result, err = cache.Lrange(tmpKey, start, end)
		if err != nil {
			fmt.Printf("Lrange Error: error=%s\n", err)
			return
		}
		if len(result) == 0{
			fmt.Printf("key %s Lrange no result\n", tmpKey)
		} else {
			fmt.Printf("Key %s Lrange success value %v\n", tmpKey, result)
		}

		l = cache.Llen(tmpKey)
		fmt.Printf("After Delete %d len\n", l)
	*/
}
