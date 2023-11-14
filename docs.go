package goapidoc

import (
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// parseDocs 根据group前缀挑选group.js文件，并解析成docs格式的语法
func parseDocs(file string) docsDef {
	comment, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	var docs docsDef
	matches := regexp.MustCompile(`(?s)/\*\*(.*?)\*/`).FindAllStringSubmatch(string(comment), -1)
	for _, match := range matches {
		contents := parserContent(match[1])
		if len(contents) < 1 {
			continue
		}

		decContents := unmarshal(contents)
		docs = append(docs, &decContents)
	}

	return docs
}

func parserContent(content string) []string {
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

func (docs docsDef) match(api ApiDef) *docDef {
	doc := new(docDef)

	for i := range docs {
		if docs[i].Api.Path != api.Path || strings.ToLower(docs[i].Api.Method) != strings.ToLower(api.Method) {
			continue
		}
		if doc == nil {
			doc = docs[i]
		} else if docs[i].ApiVersion == api.version {
			doc = docs[i]
		}
	}

	return doc
}

// saveAs 最终结果保存为文件
func (docs docsDef) saveAs(file string) docsDef {
	var docStr string
	for _, d := range docs {
		s, err := marshal(*d)
		if err == nil {
			docStr += s
		}
	}

	if len(docStr) < 1 {
		return nil
	}

	os.WriteFile(file, []byte(docStr), 0644)

	return docs
}

func (docs docsDef) toHtmlDoc(configPath, docPath string) {
	cmd := exec.Command("apidoc", "-c", configPath, "-o", docPath)
	_, err := cmd.Output()
	if err != nil {
		slog.Error("使用apidoc生成文档错误", err)
		return
	}

	slog.Info("使用apidoc生成文档成功")
}
