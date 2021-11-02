package main

import "go-chat/config"

func main() {
	// 第一步：初始化配置信息
	conf := config.Init("./../config.yaml")

	Initialize(conf)
}
