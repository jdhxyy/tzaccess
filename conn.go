// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 连接父路由
// Authors: jdh99 <jdh821@163.com>

package tzaccess

import (
	"github.com/jdhxyy/knock"
	"github.com/jdhxyy/lagan"
	"github.com/jdhxyy/utz"
	"net"
	"time"
)

// 最大连接次数.超过连接次数这回清除父路由IA地址,重连父路由
const connNumMax = 3

var connNum = 0

func init() {
	knock.Register(utz.HeaderCmp, utz.CmpMsgTypeAckConnectParentRouter, dealAckConnectParentRouter)
	go connThread()
	go connTimeout()
}

// dealAckConnectParentRouter 处理应答连接帧
// 返回值是应答数据和应答标志.应答标志为false表示不需要应答
func dealAckConnectParentRouter(req []uint8, params ...interface{}) []uint8 {
	if len(params) != 1 {
		lagan.Warn(tag, "deal conn failed.params len is wrong:%d", len(params))
		return nil
	}

	addr := params[0].(*net.UDPAddr)
	if addr.IP.Equal(parent.addr.IP) == false || addr.Port != parent.addr.Port {
		lagan.Warn(tag, "deal conn failed.ip is not match.ip:%v core ip:%v", addr, coreAddr)
		return nil
	}

	if len(req) == 0 {
		lagan.Warn(tag, "deal conn failed.payload len is wrong:%d", len(req))
		return nil
	}

	j := 0
	if req[j] != 0 {
		lagan.Warn(tag, "deal conn failed.error code:%d", req[j])
		return nil
	}
	j++

	if len(req) != 2 {
		lagan.Warn(tag, "deal conn failed.payload len is wrong:%d", len(req))
		return nil
	}

	connNum = 0
	parent.isConn = true
	parent.cost = req[j]
	parent.timestamp = time.Now().Unix()
	lagan.Info(tag, "conn success.parent ia:0x%x cost:%d", parent.ia, parent.cost)
	return nil
}

func connThread() {
	for {
		if standardLayerIsAllowSend() == false || parent.ia == utz.IAInvalid {
			time.Sleep(time.Second)
			continue
		}

		connNum += 1
		if connNum > connNumMax {
			connNum = 0
			parent.ia = utz.IAInvalid
			lagan.Warn(tag, "conn num is too many!")
			continue
		}
		lagan.Info(tag, "send conn frame")
		sendConnFrame()

		if parent.ia == utz.IAInvalid {
			time.Sleep(3 * time.Second)
		} else {
			time.Sleep(connInterval * time.Second)
		}
	}
}

func sendConnFrame() {
	var securityHeader utz.SimpleSecurityHeader
	securityHeader.NextHead = utz.HeaderCmp
	securityHeader.Pwd = localPwd
	payload := utz.SimpleSecurityHeaderToBytes(&securityHeader)

	var body []uint8
	body = append(body, utz.CmpMsgTypeConnectParentRouter)
	// 前缀长度
	body = append(body, 64)
	// 子膜从机固定单播地址
	body = append(body, make([]uint8, utz.IALen)...)
	// 开销值
	body = append(body, 0)
	body = utz.BytesToFlpFrame(body, true, 0)

	payload = append(payload, body...)

	var header utz.StandardHeader
	header.Version = utz.ProtocolVersion
	header.FrameIndex = utz.GenerateFrameIndex()
	header.PayloadLen = uint16(len(payload))
	header.NextHead = utz.HeaderSimpleSecurity
	header.HopsLimit = 0xff
	header.SrcIA = localIA
	header.DstIA = coreIA

	standardLayerSend(payload, &header, &parent.addr)
}

func connTimeout() {
	for {
		if parent.ia == utz.IAInvalid || parent.isConn == false {
			time.Sleep(time.Second)
			continue
		}
		if time.Now().Unix()-parent.timestamp > connTimeoutMax {
			parent.ia = utz.IAInvalid
			parent.isConn = false
		}
		time.Sleep(time.Second)
	}
}

// IsConn 是否连接核心网
func IsConn() bool {
	return parent.ia != utz.IAInvalid && parent.isConn
}

// GetParentAddr 读取父节点地址.如果父节点不存在则返回nil
func GetParentAddr() *net.UDPAddr {
	if IsConn() == false {
		return nil
	} else {
		return &parent.addr
	}
}
