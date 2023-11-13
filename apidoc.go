package goapidoc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// NewDoc todo doc.ApiPermission | doc.ApiPrivate | doc.ApiUse
func NewDoc(api ApiDef) docsDef {
	group := api.extractGroupFromRoute()
	version := getVersion()

	docs := ParseDocs(group)
	doc := docs.match(api)
	oldDoc := new(docDef)

	if doc == nil || doc.ApiVersion != version {
		if doc != nil {
			oldDoc = doc
		}

		doc = new(docDef)
		docs = append(docs, doc)

		doc.Api.Title = fmt.Sprintf("%s[%s]", strings.ToUpper(api.Method), api.Path)
		doc.Api.Method = strings.ToLower(api.Method)
		doc.Api.Path = api.Path

		doc.ApiGroup = group
		doc.ApiVersion = version
		doc.ApiDescription = fmt.Sprintf("%s %s ", strings.ToUpper(api.Method), api.Path)
		doc.ApiName = fmt.Sprintf("%s-%s", api.Method, api.Path)
	}

	for k, v := range api.Params {
		doc.appendParam(&doc.ApiParam, generateParamDef(k, v, ""), oldDoc.ApiParam)
	}

	for k, v := range api.Headers {
		doc.appendParam(&doc.ApiHeader, generateParamDef(k, v, ""), oldDoc.ApiHeader)
	}

	for k, v := range api.Queries {
		doc.appendParam(&doc.ApiQuery, generateParamDef(k, v, ""), oldDoc.ApiQuery)
	}

	respExample := respExampleDef{
		Type:    "json",
		Title:   http.StatusText(api.StatusCode),
		Http:    fmt.Sprintf("HTTP/1.1 %d %s", api.StatusCode, http.StatusText(api.StatusCode)),
		Example: string(api.ResponseBody),
	}

	if api.StatusCode < 400 {
		doc.appendExample(&doc.ApiSuccessExample, respExample)
	} else {
		doc.appendExample(&doc.ApiErrorExample, respExample)
	}

	fieldTypes := api.extractFieldTypeFromResponse()
	for field, typ := range fieldTypes {
		param := generateParamDef(field, "", typ)
		if api.StatusCode < 400 {
			doc.appendParam(&doc.ApiSuccess, param, oldDoc.ApiSuccess)
		} else {
			doc.appendParam(&doc.ApiError, param, oldDoc.ApiError)
		}
	}

	return docs
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

func (docs docsDef) match(api ApiDef) *docDef {
	for i := range docs {
		if docs[i].Api.Path == api.Path && docs[i].Api.Method == api.Method {
			return docs[i]
		}
	}

	return nil
}

func (doc *docDef) appendParam(ps *[]paramDef, p paramDef, olds []paramDef) {
	for i := range *ps {
		if (*ps)[i].Type.Type == p.Type.Type && (*ps)[i].Field.Name == p.Field.Name {
			(*ps)[i].Field.Required = p.Field.Required
			(*ps)[i].Field.Default = p.Field.Default
			return
		}
	}

	for _, old := range olds {
		if old.Field.Name != p.Field.Name {
			continue
		}

		if len(old.Description) > 0 && len(p.Description) < 1 {
			p.Description = old.Description
		}

		if len(old.Type.Size) > 0 && len(p.Type.Size) < 1 {
			p.Type.Size = old.Type.Size
		}

		if old.Type.Allows != nil && p.Type.Allows == nil {
			p.Type.Allows = old.Type.Allows
		}

		if old.Field.Required && !p.Field.Required {
			p.Field.Required = old.Field.Required
		}

		if len(old.Field.Default) > 0 && len(p.Field.Default) < 1 {
			p.Field.Default = old.Field.Default
		}

	}

	*ps = append(*ps, p)
}

func (doc *docDef) appendExample(rs *[]respExampleDef, r respExampleDef) {
	for i := range *rs {
		if (*rs)[i].Http == r.Http {
			(*rs)[i].Example = r.Example
			return
		}
	}

	*rs = append(*rs, r)
}

func getVersion() string {
	apidoc, err := readSrcFileContent("apidoc.json")
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

func readSrcFileContent(file string) ([]byte, error) {
	filePaths := []string{"doc/src/", "../doc/src/", "../../doc/src/", "template/src/"}

	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath + file)
		if err == nil && len(content) > 0 {
			return content, nil
		}
	}

	return nil, fmt.Errorf("read file [%v] failed, path: [%s]", file, filePaths)
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
