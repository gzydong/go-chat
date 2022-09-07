package entity

// IM 渠道分组(用于业务划分，业务间相互隔离)
const (
	// ImChannelDefault 默认分组
	ImChannelDefault = "default" // im.Sessions.Default.Name()
	ImChannelExample = "example" // im.Sessions.Example.Name()
)

const (
	// ImTopicDefault 默认渠道消息订阅
	ImTopicDefault        = "im:message:default:all"
	ImTopicDefaultPrivate = "im:message:default:%s"

	// ImTopicExample Example渠道消息订阅
	ImTopicExample        = "im:message:example:all"
	ImTopicExamplePrivate = "im:message:example:%s"
)
