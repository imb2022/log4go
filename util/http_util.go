/**@description	HTTP相关的方法集合
@author imb2022
**/
package util

import (
	"math/rand"
	"net"
	"strings"
)

// GetLocalIpByTcp 有外网的情况下, 通过tcp访问获得本机ip地址
func GetLocalIpByTcp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		if conn != nil {
			_ = conn.Close()
		}
		conn, err = net.Dial("tcp", "www.baidu.com:80")
		if err != nil {
			return ""
		}
	}

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// RandomInt 随机数
func RandomInt(num int) int {
	return rand.Intn(65536) % num
}

// MaxInt return max int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt return min int
func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}
