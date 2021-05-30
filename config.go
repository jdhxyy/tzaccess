// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 配置文件
// Authors: jdh99 <jdh821@163.com>

package tzaccess

import (
	"net"
)

const (
	tag = "tzaccess"

	// 最大帧字节数
	frameMaxLen = 4096

	protocolNum = 0

	// 连接间隔.单位:s
	connInterval = 30
	// 连接超时时间.单位:s
	connTimeoutMax = 120
)

// SendFunc 发送函数.addr:目标地址
type SendFunc func(data []uint8, addr *net.UDPAddr)

// IsAllowSendFunc 是否允许发送
type IsAllowSendFunc func() bool

// tParentInfo 父路由信息
type tParentInfo struct {
	ia        uint64
	addr      net.UDPAddr
	cost      uint8
	isConn    bool
	timestamp int64
}

var parent tParentInfo

// 本机单播地址
var localIA uint64
var localPwd string

// 核心网参数
var coreIA uint64 = 0x2141000000000002
var coreIP = "115.28.86.171"
var corePort = 12914
var coreAddr net.UDPAddr

// 发送函数
var sendFunc SendFunc = nil
var isAllowSendFunc IsAllowSendFunc = nil

func init() {
	coreAddr = net.UDPAddr{IP: net.ParseIP(coreIP), Port: corePort}
}

// ConfigCoreParam 配置核心网参数
func ConfigCoreParam(ia uint64, ip string, port int) {
	coreIA = ia
	coreIP = ip
	corePort = port
	coreAddr = net.UDPAddr{IP: net.ParseIP(coreIP), Port: corePort}
}
