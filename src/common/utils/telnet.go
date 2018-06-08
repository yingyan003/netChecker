package utils

import (
	"strconv"
	"net"
	"time"
)

func Telnet(ip string, port int32, timeout int32) bool {
	portnum := strconv.Itoa(int(port))
	address := ip + ":" + portnum
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Errorf("Telnet Error: err=%s, ip=%s, port=%d, timeout=%d\n", err, ip, port, timeout)
		return false
	}
	defer conn.Close()
	return true
}
