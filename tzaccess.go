// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 接入海萤物联网
// Authors: jdh99 <jdh821@163.com>

package tzaccess

func Load(ia uint64, pwd string, send SendFunc, isAllowSend IsAllowSendFunc) {
	localIA = ia
	localPwd = pwd
	sendFunc = send
	isAllowSendFunc = isAllowSend
}
