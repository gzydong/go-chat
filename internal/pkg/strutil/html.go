package strutil

import "regexp"

// ParseHtmlImage 解析 Md 文本中的图片信息
func ParseHtmlImage(text string) string {
	reg, _ := regexp.Compile(`<img src=[\'|\"](.*?)[\'|\"].*?[\/]?>`)

	items := reg.FindAllStringSubmatch(text, 1)
	for _, item := range items {
		return item[1]
	}

	return ""
}

// ParseHtmlImageAll 解析 Md 文本中的所有图片信息
func ParseHtmlImageAll(text string) []string {
	reg, _ := regexp.Compile(`<img src=[\'|\"](.*?)[\'|\"].*?[\/]?>`)

	list := make([]string, 0)

	for _, item := range reg.FindAllStringSubmatch(text, -1) {
		list = append(list, item[1])
	}

	return list
}
