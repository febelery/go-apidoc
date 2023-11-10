package goapidoc

import (
	"regexp"
	"strings"
)

func ParseDocs(group string) DocsDef {
	comment, err := readSrcFileContent(group + ".js")
	if err != nil {
		return nil
	}

	var docs DocsDef
	matches := regexp.MustCompile(`(?s)/\*\*(.*?)\*/`).FindAllStringSubmatch(string(comment), -1)
	for _, match := range matches {
		contents := parseContent(match[1])
		if len(contents) < 1 {
			continue
		}

		decContents := Unmarshal(contents)
		docs = append(docs, &decContents)
	}

	return docs
}

func parseContent(content string) []string {
	apiAnnotePrefix := "@api"
	var resultLines []string

	// 移除没有文档结构的块
	index := strings.Index(content, apiAnnotePrefix)
	if index == -1 {
		return resultLines
	}
	content = content[index:]

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// 判断第一个字符串是否为 *
		if strings.Index(strings.TrimSpace(line), "*") == 0 {
			line = line[strings.Index(line, "*")+1:]
		}

		// 移除块内容中的空行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 组装
		if strings.Contains(line, apiAnnotePrefix) {
			resultLines = append(resultLines, strings.TrimSpace(line))
		} else {
			linesLen := len(resultLines)
			if linesLen < 1 {
				continue
			}

			resultLines[linesLen-1] = strings.Join([]string{resultLines[linesLen-1], line}, "\n")
		}
	}

	return resultLines
}
