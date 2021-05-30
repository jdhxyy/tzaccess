package tzaccess

import (
	"github.com/jdhxyy/lagan"
	"net"
	"testing"
)

func TestCase1(t *testing.T) {
	_ = lagan.Load(0)
	lagan.SetFilterLevel(lagan.LevelDebug)
	lagan.EnableColor(true)

	addr := net.UDPAddr{IP: net.ParseIP("192.168.1.119"), Port: 12021}
	listener, err := net.ListenUDP("udp", &addr)
	if err != nil {
		lagan.Error(tag, "bind pipe net failed:%v", err)
		return
	}

	Load(
		0x2141000000000401,
		"abc123",
		func(data []uint8, addr *net.UDPAddr) {
			_, err := listener.WriteToUDP(data, addr)
			if err != nil {
				lagan.Error(tag, "udp send error:%v addr:%v", err, addr)
				return
			}
			lagan.Info(tag, "udp send:addr:%v len:%d", addr, len(data))
			lagan.PrintHex(tag, lagan.LevelDebug, data)
		},
		func() bool {
			return true
		})

	go func() {
		data := make([]uint8, frameMaxLen)
		for {
			num, addr, err := listener.ReadFromUDP(data)
			if err != nil {
				lagan.Error(tag, "listen pipe net failed:%v", err)
				continue
			}
			if num <= 0 {
				continue
			}
			lagan.Info(tag, "udp rx:%v len:%d", addr, num)
			lagan.PrintHex(tag, lagan.LevelDebug, data[:num])

			Receive(data[:num], addr)
		}
	}()

	select {}
}
