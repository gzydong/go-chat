package strutil

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

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

var matchMdImageReg = regexp.MustCompile(`\!\[(.*?)\]\((.*?)\)`)

func ParseMarkdownImages(content string) []string {
	matches := matchMdImageReg.FindAllStringSubmatch(content, -1)

	items := make([]string, 0)
	for _, match := range matches {
		items = append(items, match[2])
	}

	return items
}
