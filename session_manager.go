package main

import (
	"sync"
	"time"
)

const DefaultAliveTime = 300

type UserInfo struct {
	Host     string `json:"host"`
	NickName string `json:"nickName"`
	//存在的秒数，小于等于0就删除
	time int
}

type SessionManager struct {
	cache     []UserInfo
	isRunning bool
	mux       sync.Mutex
	ticker    *time.Ticker
}

// IsExist /** 判断该host是否已存在
func (manager *SessionManager) IsExist(host string) (*UserInfo, bool) {
	manager.mux.Lock()
	defer manager.mux.Unlock()
	for i := 0; i < len(manager.cache); i++ {
		if manager.cache[i].Host == host {
			return &manager.cache[i], true
		}
	}
	return nil, false
}

func (manager *SessionManager) GetNickNameByHost(host string) string {
	info, _ := manager.IsExist(host)
	return info.NickName
}

func (manager *SessionManager) AddSession(host string) {
	info, exist := manager.IsExist(host)
	manager.mux.Lock()
	defer manager.mux.Unlock()
	if !exist {
		//新加入
		manager.cache = append(manager.cache, UserInfo{
			Host:     host,
			NickName: host,
			time:     DefaultAliveTime,
		})
	} else {
		info.time = DefaultAliveTime
	}
}

// Start 开始运行
func (manager *SessionManager) Start() {
	if manager.isRunning {
		return
	}
	if manager.ticker == nil {
		manager.ticker = time.NewTicker(time.Second)
	}
	manager.isRunning = true
	go func() {
		for manager.isRunning {
			// 等待触发器触发事件
			<-manager.ticker.C
			manager.mux.Lock()
			// 遍历主机，将每秒存活时间减一，为0时移除
			for i := 0; i < len(manager.cache); i++ {
				if manager.cache[i].time-1 <= 0 {
					manager.cache = append(manager.cache[:i])
				} else {
					manager.cache[i].time--
				}
			}
			manager.mux.Unlock()
		}
	}()
}

// Rename 设置nickname
func (manager *SessionManager) Rename(host string, nickName string) {
	info, exist := manager.IsExist(host)
	if !exist {
		return
	}
	info.NickName = nickName
}

// Stop 停止
func (manager *SessionManager) Stop() {
	manager.isRunning = false
}
