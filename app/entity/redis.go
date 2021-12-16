package entity

import "fmt"

const (
	IMGatewayAll     = "im:gateway:all"
	IMGatewayPrivate = "im:gateway:%s"
)

func GetIMGatewayPrivate(sid string) string {
	return fmt.Sprintf(IMGatewayPrivate, sid)
}
