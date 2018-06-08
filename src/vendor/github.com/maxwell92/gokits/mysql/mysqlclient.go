package mysql

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	mylog "github.com/maxwell92/gokits/log"
)

var log = mylog.Log

type MysqlClient struct {
	DB *sql.DB
}

type MySQL interface {
	Open(cfg *DefaultConfig)
	Close()
	Conn() *sql.DB
	Ping()
}

var instance MySQL
var once sync.Once

func MysqlInstance() MySQL {
	once.Do(func() {
		instance = new(MysqlClient)
	})
	return instance
}

func FakeMysqlInstance() MySQL {
	once.Do(func() {
		instance = new(MysqlClient)
	})
	return instance
}

type DefaultConfig struct {
	Host          string
	Port          string
	Driver        string
	User          string
	Pass          string
	Database      string
	MaxActiveConn string
	MaxIdleConn   string
}

const (
	DB_CONNECTION_SUFFIX = "?parseTime=true"
	DELAY_MILLISECONDS   = 5000
)

func (c *MysqlClient) Open(cfg *DefaultConfig) {
	// endpoint := config.Instance().GetDbEndpoint()
	host := cfg.Host + ":" + cfg.Port
	endpoint := cfg.User + ":" + cfg.Pass + "@tcp(" + host + ")/" + cfg.Database + DB_CONNECTION_SUFFIX

	// db, err := sql.Open(config.DATABASE_DRIVER, endpoint)
	db, err := sql.Open(cfg.Driver, endpoint)

	if err != nil {
		log.Fatalf("MysqlClient Open Error: err=%s", err)
		return
	}

	// Set Connection Pool
	maxActive, _ := strconv.Atoi(cfg.MaxActiveConn)
	// db.SetMaxOpenConns(config.Instance().RedisMaxActiveConn)
	db.SetMaxOpenConns(maxActive)
	maxIdle, _ := strconv.Atoi(cfg.MaxIdleConn)
	// db.SetMaxIdleConns(config.Instance().RedisMaxIdleConn)
	db.SetMaxIdleConns(maxIdle)

	c.DB = db
	// log.Infof("MysqlClient Open Success: host=%s", config.DB_HOST)
	log.Tracef("MysqlClient Open Success: host=%s", cfg.Host)
}

func (c *MysqlClient) Close() {
	c.DB.Close()
}

func (c *MysqlClient) Conn() *sql.DB {
	return c.DB
}

// Ping the connection, keep connection alive
func (c *MysqlClient) Ping() {
	select {
	case <-time.After(time.Millisecond * time.Duration(DELAY_MILLISECONDS)):
		err := c.DB.Ping()
		if err != nil {
			log.Fatalf("MysqlClient Ping Error: err=%s", err)
			c.Open(nil)
		}
	}
}
