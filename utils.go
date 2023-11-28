package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

// GenerateRandomID 设置随机种子
func GenerateRandomID(length int) string {

	// 定义可用的字符集
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// 生成随机字符序列
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b[:])
}

var localIP = ""

func GetClientIP(context *gin.Context) string {
	if context.ClientIP() == "::1" {
		if localIP != "" {
			return localIP
		}
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			fmt.Println("无法获取IP地址:", err)
		}
		var ip net.IP
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ip4 := ipnet.IP.To4(); ip4 != nil && strings.HasPrefix(ip4.String(), "192") {
					ip = ipnet.IP
					break
				}
			}
		}
		localIP = ip.To4().String()
		return localIP
	} else {
		return context.ClientIP()
	}
}

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func OK() (int, Response) {
	return http.StatusOK, Response{
		Code: 0,
		Data: "",
		Msg:  "请求成功",
	}
}

func OKWithData(data any) (int, Response) {
	return http.StatusOK, Response{
		Code: 0,
		Data: data,
		Msg:  "请求成功",
	}
}

func OKWithMsg(msg string) (int, Response) {
	return http.StatusOK, Response{
		Code: 0,
		Data: "",
		Msg:  msg,
	}
}

func FailWithMsg(msg string) (int, Response) {
	return http.StatusOK, Response{
		Code: -1,
		Data: "",
		Msg:  msg,
	}
}
