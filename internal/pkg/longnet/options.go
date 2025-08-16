package longnet

import (
	"crypto/tls"
	"time"
)

type WSSConfig struct {
	Addr      string // WSS 监听地址
	Path      string // WSS path 路径
	TLSEnable bool   // TLS 是否启用
}

func (w WSSConfig) getPath() string {
	if w.Path == "" {
		return "/"
	}

	return w.Path
}

type TCPConfig struct {
	Addr      string // TCP 监听地址
	TLSEnable bool   // TLS 是否启用
}

type Options struct {
	PingInterval time.Duration // 心跳间隔
	PingTimeout  time.Duration // 心跳超时
	ReadTimeout  time.Duration // 读超时时间
	WriteTimeout time.Duration // 写超时时间

	MaxOpenConns  int // 最大连接数量 -1:不限制
	MaxPacketSize int // 最大数据包大小

	WSSConfig *WSSConfig  // WSS 配置
	TCPConfig *TCPConfig  // TCP 配置
	TLSConfig *tls.Config //
}

func (o Options) init() Options {
	if o.ReadTimeout <= 0 {
		o.ReadTimeout = 3 * time.Minute
	}

	if o.WriteTimeout <= 0 {
		o.WriteTimeout = 3 * time.Second
	}

	if o.PingInterval <= 0 {
		o.PingInterval = 30 * time.Second
	}

	if o.PingTimeout <= 0 {
		o.PingTimeout = 3 * o.PingInterval
	}

	if o.MaxOpenConns <= 0 {
		o.MaxOpenConns = -1
	}

	if o.MaxPacketSize <= 0 {
		o.MaxPacketSize = 1 << 20 // 1M
	}

	// WSS 配置
	if o.WSSConfig == nil {
		o.WSSConfig = &WSSConfig{
			Addr: ":9501",
			Path: "/wss/chat.io",
		}
	}

	if o.TCPConfig != nil && o.TCPConfig.Addr == "" {
		o.TCPConfig.Addr = ":9502"
	}

	return o
}

type ClientOptions struct {
	DialTimeout  time.Duration // 连接超时
	PingInterval time.Duration // 心跳间隔
	PingTimeout  time.Duration // 心跳超时
	ReadTimeout  time.Duration // 读超时时间
	WriteTimeout time.Duration // 写超时时间
}
