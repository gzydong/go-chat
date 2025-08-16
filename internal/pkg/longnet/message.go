package longnet

import (
	"fmt"
	"time"
)

const PacketMaxSize = 2 << 20 // 2MB

const (
	Ping         = 1001
	Pong         = 1002
	Ack          = 1003
	Connect      = 1004
	Identity     = 1005
	Authorize    = 1006
	Unauthorized = 1007
)

type Packet struct {
	Cmd     int32  // 命令
	Payload []byte // 消息体
	Msgid   int64  // 消息id
	Version uint8  // 版本号
}

func NewCustomizePacket(cmd int32, body []byte) *Packet {
	return &Packet{
		Cmd:     cmd,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

// NewCustomizePacketWithCompress 创建自定义命令（基于已压缩的body）
func NewCustomizePacketWithCompress(cmd int32, body []byte) *Packet {
	return &Packet{
		Cmd:     cmd,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewPingPacket() *Packet {
	return &Packet{
		Cmd:     Ping,
		Payload: nil,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewPondPacket() *Packet {
	return &Packet{
		Cmd:     Pong,
		Payload: nil,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewConnectPacket(body []byte) *Packet {
	return &Packet{
		Cmd:     Connect,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewIdentityPacket(body []byte) *Packet {
	return &Packet{
		Cmd:     Identity,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewAuthorizePacket(body []byte) *Packet {
	return &Packet{
		Cmd:     Authorize,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewUnauthorizedPacket(body []byte) *Packet {
	return &Packet{
		Cmd:     Unauthorized,
		Payload: body,
		Msgid:   time.Now().UnixMilli(),
	}
}

func NewAckPacket(ackid int64) *Packet {
	return &Packet{
		Cmd:     Ack,
		Msgid:   time.Now().UnixMilli(),
		Payload: []byte(fmt.Sprintf(`{"ackid":%d}`, ackid)),
	}
}
