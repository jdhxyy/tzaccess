// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 标准层处理模块
// Authors: jdh99 <jdh821@163.com>

package tzaccess

import (
	"github.com/jdhxyy/lagan"
	"github.com/jdhxyy/utz"
	"net"
)

// standardLayerRxCallback 接收回调函数
type standardLayerRxCallback func(data []uint8, standardHeader *utz.StandardHeader, addr *net.UDPAddr)

var standardLayerObservers []standardLayerRxCallback

// standardLayerRx 标准层接收
func standardLayerRx(data []uint8, addr *net.UDPAddr) {
	header := getStandardHeader(data)
	if header == nil {
		return
	}
	notifyStandardLayerObservers(data[utz.NLv1HeadLen:], header, addr)
}

func getStandardHeader(data []uint8) *utz.StandardHeader {
	header, offset := utz.BytesToStandardHeader(data)
	if header == nil || offset == 0 {
		lagan.Debug(tag, "get standard header failed:bytes to standard header failed")
		return nil
	}
	if header.Version != utz.ProtocolVersion {
		lagan.Debug(tag, "get standard header failed:protocol version is not match:%d", header.Version)
		return nil
	}
	if int(header.PayloadLen)+offset != len(data) {
		lagan.Debug(tag, "get standard header failed:payload len is not match:%d", header.PayloadLen)
		return nil
	}

	return header
}

func notifyStandardLayerObservers(data []uint8, standardHeader *utz.StandardHeader, addr *net.UDPAddr) {
	n := len(standardLayerObservers)
	for i := 0; i < n; i++ {
		standardLayerObservers[i](data, standardHeader, addr)
	}
}

// standardLayerRegisterRxObserver 注册接收观察者
func standardLayerRegisterRxObserver(callback standardLayerRxCallback) {
	standardLayerObservers = append(standardLayerObservers, callback)
}

// standardLayerSend 基于标准头部发送
func standardLayerSend(data []uint8, standardHeader *utz.StandardHeader, addr *net.UDPAddr) {
	dataLen := len(data)
	if dataLen > frameMaxLen {
		lagan.Error(tag, "standard layer send failed!data len is too long:%d src ia:0x%x dst ia:0x%x", dataLen,
			standardHeader.SrcIA, standardHeader.DstIA)
		return
	}
	if standardHeader.PayloadLen != uint16(dataLen) {
		standardHeader.PayloadLen = uint16(dataLen)
	}
	frame := utz.StandardHeaderToBytes(standardHeader)
	frame = append(frame, data...)
	sendFunc(frame, addr)
}

// standardLayerIsAllowSend 是否允许发送
func standardLayerIsAllowSend() bool {
	if isAllowSendFunc == nil {
		return false
	}
	return isAllowSendFunc()
}

// Receive 用户接收到数据时需回调本函数
func Receive(data []uint8, addr *net.UDPAddr) {
	standardLayerRx(data, addr)
}
