package entity

import "fmt"

// IM 渠道分组(用于业务划分，业务间相互隔离)
const (
	// ImChannelDefault 默认分组
	ImChannelDefault = "default" // im.Sessions.Default.Name()
	ImChannelExample = "example" // im.Sessions.Example.Name()
)

const (
	IMGatewayAll     = "im:gateway:all"
	IMGatewayPrivate = "im:gateway:%s"
)

func GetIMGatewayPrivate(sid string) string {
	return fmt.Sprintf(IMGatewayPrivate, sid)
}
