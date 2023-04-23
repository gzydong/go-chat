package strutil

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

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

var amReg = regexp.MustCompile(`<a href="([^"]*)" alt="link"[^>]*>(.*?)</a>|<img src="([^"]*)" alt="img"[^>]*/>`)

func EscapeHtml(value string) string {
	items := make(map[string]string)
	for index, v := range amReg.FindAllString(value, -1) {
		val := fmt.Sprintf("{#%d#}", index)
		items[val] = v
		value = strings.Replace(value, v, val, -1)
	}

	value = html.EscapeString(value)
	if len(items) == 0 {
		return value
	}

	for k, v := range items {
		value = strings.Replace(value, k, v, -1)
	}

	return value
}

var imgReg = regexp.MustCompile(`<img .*?>`)

func ReplaceImgAll(value string) string {
	return strings.TrimSpace(string(imgReg.ReplaceAll([]byte(value), []byte(""))))
}
