package entity

import "fmt"

const (
	IMGatewayAll     = "ws:gateway:all"
	IMGatewayPrivate = "ws:gateway:%s"
)

func GetIMGatewayPrivate(sid string) string {
	return fmt.Sprintf(IMGatewayPrivate, sid)
}
