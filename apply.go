// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 申请父路由
// Authors: jdh99 <jdh821@163.com>

package tzaccess

import (
	"github.com/jdhxyy/knock"
	"github.com/jdhxyy/lagan"
	"github.com/jdhxyy/utz"
	"net"
	"time"
)

func init() {
	knock.Register(utz.HeaderCmp, utz.CmpMsgTypeAssignSlaveRouter, dealAssignSlaveRouter)
	go applyThread()
}

// dealAssignSlaveRouter 处理分配从机帧
// 返回值是应答数据和应答标志.应答标志为false表示不需要应答
func dealAssignSlaveRouter(req []uint8, params ...interface{}) []uint8 {
	if len(params) != 1 {
		lagan.Warn(tag, "deal apply failed.params len is wrong:%d", len(params))
		return nil
	}

	addr := params[0].(*net.UDPAddr)
	if addr.IP.Equal(coreAddr.IP) == false || addr.Port != coreAddr.Port {
		lagan.Warn(tag, "deal apply failed.ip is not match.ip:%v core ip:%v", addr, coreAddr)
		return nil
	}

	if len(req) == 0 {
		lagan.Warn(tag, "deal apply failed.payload len is wrong:%d", len(req))
		return nil
	}

	j := 0
	if req[j] != 0 {
		lagan.Warn(tag, "deal apply failed.error code:%d", req[j])
		return nil
	}
	j++

	if len(req) != 16 {
		lagan.Warn(tag, "deal apply failed.payload len is wrong:%d", len(req))
		return nil
	}

	parent.ia = utz.BytesToIA(req[j : j+utz.IALen])
	j += utz.IALen

	ip := make([]uint8, 4)
	copy(ip, req[j:j+4])
	j += 4
	port := (int(req[j]) << 8) + int(req[j+1])
	j += 2
	parent.addr = net.UDPAddr{IP: net.IPv4(ip[0], ip[1], ip[2], ip[3]), Port: port}

	lagan.Info(tag, "apply success.parent ia:0x%x addr:%v cost:%d", parent.ia, parent.addr, req[j])
	return nil
}

func applyThread() {
	for {
		if standardLayerIsAllowSend() == false || parent.ia != utz.IAInvalid {
			time.Sleep(time.Second)
			continue
		}

		lagan.Info(tag, "send apply frame")
		sendApplyFrame()

		if parent.ia == utz.IAInvalid {
			time.Sleep(3 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func sendApplyFrame() {
	var securityHeader utz.SimpleSecurityHeader
	securityHeader.NextHead = utz.HeaderCmp
	securityHeader.Pwd = localPwd
	payload := utz.SimpleSecurityHeaderToBytes(&securityHeader)

	var body []uint8
	body = append(body, utz.CmpMsgTypeRequestSlaveRouter)
	body = append(body, utz.IAToBytes(parent.ia)...)
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

	standardLayerSend(payload, &header, &coreAddr)
}
