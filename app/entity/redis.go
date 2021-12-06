package entity

import "fmt"

const (
	SubscribeWsGatewayAll     = "ws:gateway:all"
	SubscribeWsGatewayPrivate = "ws:gateway:%s"
)

func GetSubscribeWsGatewayPrivate(sid string) string {
	return fmt.Sprintf(SubscribeWsGatewayPrivate, sid)
}
