package strutil

import (
	"testing"
)

func TestParseImage(t *testing.T) {

	// r := regexp.MustCompile(`{[^}]*}`)
	// matches := r.FindAllString("{city}, {state} {zip}", -1)
	//
	// fmt.Println(matches)
	// <img src=\"http://im-img.gzydong.club/media/images/notes/20210816/611a59b851992zHyLjgAOzJ11n3B4_1140x760.jpg\" alt=\"4437e6581b221c48d1781f.jpg\" />

	ParseImage(`<h1><a i<p><img src=\"http://im-img.gzydong.club/media/images/notes/20210816/611a59b851992zHyLjgAOzJ11n3B4_1140x760.jpg" alt="4437e6581b221c48d1781f.jpg" /><br />\nLumen IM 是一个网页版在线即时聊天项目，前端使用 Element-ui + Vue，后端采用了基于 Swoole 开发的 Hyp<img src="http://im-img.gzydong.club/media/images/notes/20210816/611a59b851992zHyLjgAOzJ11n3B4_1140x760.jpg" alt="4437e6581b221c48d1781f.jpg" />erf 协程框架进行接口开发，并使用 WebSocket 服务进行消息实时推送。目前后端 WebSocket 已支持分布式集群部署。</p>\n<p>目前该项目是在 旧版本 项目的基础上进行了后端重构，且前后端都有较大的改动。</p>\n<h2><a id=\"_6\"></a>功能模块</h2>`)
}
