package goapidoc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func New(api ApiDef) {
	if api.Skip {
		return
	}

	api.group = api.extractGroupFromRoute()
	api.version = getVersion()
	file := getFilePath(fmt.Sprintf("src/%s.js", api.group))

	docs := parseDocs(file)
	doc := docs.match(api)
	olderDoc := *doc
	if doc.ApiVersion != api.version {
		doc = new(docDef)
		docs = append(docs, doc)
	}
	doc.generateWithOldDoc(api, olderDoc)

	docs.saveAs(file).toHtmlDoc(getFilePath("src/apidoc.json"), getFilePath("api"))
}

func (api ApiDef) extractFieldTypeFromResponse() map[string]string {
	var data map[string]any
	fieldMap := make(map[string]string)

	err := json.Unmarshal(api.ResponseBody, &data)
	if err == nil {
		extractFieldTypes("", data, fieldMap)
	}

	return fieldMap
}

func (api ApiDef) extractGroupFromRoute() string {
	// eg: /api/v1/user
	matches := regexp.MustCompile(`/?(?:[\w-_]+/)*v\d+/([\w-_]+)/?`).FindStringSubmatch(api.Path)
	if len(matches) > 1 {
		return matches[1]
	}

	// eg: /api/user
	matches = regexp.MustCompile(`/?(?:[\w-_]+/)?api/([\w-_]+)/?`).FindStringSubmatch(api.Path)
	if len(matches) > 1 {
		return matches[1]
	}

	// 使用字符串分割提取路径中的分组
	pathParts := strings.Split(strings.Trim(api.Path, "/"), "/")
	if len(pathParts) > 0 {
		return pathParts[0]
	}

	return ""
}

func getVersion() string {
	apidoc, err := os.ReadFile(getFilePath("src/apidoc.json"))
	if err != nil {
		slog.Warn("no version found, use default[0.0.1]. err: ", err)
		return "0.0.1"
	}

	var ver struct {
		Version string
	}

	err = json.Unmarshal(apidoc, &ver)
	if err != nil {
		slog.Warn("no version found, use default[0.0.1]. err: ", err)
		return "0.0.1"
	}

	return ver.Version
}

func getFilePath(file string) string {
	dirPath := []string{"doc/", "../doc/", "../../doc/", "template/"}
	fullPath := ""

	for _, path := range dirPath {
		if _, err := os.Stat(path); err == nil {
			fullPath = path
			break
		}
	}
	fullPath += file

	return fullPath

}

func extractFieldTypes(prefix string, data any, fieldMap map[string]string) {
	switch value := data.(type) {
	case map[string]any:
		for key, val := range value {
			// {"data":xxxx} 这种格式默认跳过[data.]前缀
			if key == "data" {
				extractFieldTypes("", val, fieldMap)
				continue
			}

			fieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, key), ".")
			if _, ok := fieldMap[fieldName]; ok {
				continue
			}

			var fieldType string
			switch val.(type) {
			case map[string]any:
				fieldType = "Object"
			case []any:
				fieldType = ""
			case string:
				fieldType = "String"
			case bool:
				fieldType = "Boolean"
			default:
				fieldType = "Number"
			}
			fieldMap[fieldName] = fieldType

			extractFieldTypes(fieldName, val, fieldMap)
		}
	case []any:
		if len(value) == 0 && fieldMap[prefix] == "" {
			fieldMap[prefix] = "[]"
		}

		for _, val := range value {
			valType := reflect.TypeOf(val)

			switch valType.Kind() {
			case reflect.Map, reflect.Slice:
				if fieldMap[prefix] == "" {
					fieldMap[prefix] = "Object[]"
				}
				extractFieldTypes(prefix, val, fieldMap)
			default:
				if fieldMap[prefix] == "" {
					switch val.(type) {
					case string:
						fieldMap[prefix] = "String[]"
					default:
						fieldMap[prefix] = "Number[]"
					}
				}
			}
		}
	}
}

func generateParamDef(key, val, typ string) paramDef {
	if typ == "" {
		typ = "String"
		if _, err := strconv.Atoi(val); err != nil {
			typ = "Number"
		}
	}

	param := paramDef{
		Type: typeDef{
			Type:   typ,
			Size:   "",
			Allows: nil,
		},
		Field: fieldDef{
			Name:     key,
			Required: false,
			Default:  val,
		},
		Description: "",
	}

	return param
}
