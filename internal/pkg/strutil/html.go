package strutil

import (
	"regexp"
)

var matchMdImageReg = regexp.MustCompile(`\!\[(.*?)\]\((.*?)\)`)

func ParseMarkdownImages(content string) []string {
	matches := matchMdImageReg.FindAllStringSubmatch(content, -1)

	items := make([]string, 0)
	for _, match := range matches {
		items = append(items, match[2])
	}

	return items
}
