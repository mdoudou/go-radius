package g

import (
	"log"
	"os"
	"net"
	"errors"
	"encoding/binary"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}

func array_include(array_a []string, array_b []string) bool { //b include a
	for _, v := range array_a {
		if in_array(v, array_b) {
			continue
		} else {
			return false
		}
	}
	return true
}

func in_array(a string, array []string) bool {
	for _, v := range array {
		if a == v {
			return true
		}
	}
	return false
}

func Ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func Long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
