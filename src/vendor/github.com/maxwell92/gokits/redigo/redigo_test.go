package redigo

import (
	"fmt"
	redis "github.com/garyburd/redigo/redis"
	"testing"
)

func Test_NewRedisClient(*testing.T) {
	pool := NewRedisClient()

	conn := pool.Get()

	if _, err := conn.Do("SET", "hello", "world"); err != nil {
		fmt.Println("OK")
	}

	hello, err := redis.String(conn.Do("GET", "hello"))

	if err != nil {
		fmt.Println("Redis get failed: ", err)
	}

	fmt.Printf("Got username: %s\n", hello)
}
