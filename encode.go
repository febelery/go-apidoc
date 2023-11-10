package goapidoc

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func Marshal(doc DocDef) (string, error) {
	var resultLines []string

	if len(doc.Api.Method) < 1 || len(doc.Api.Path) < 1 || len(doc.Api.Title) < 1 {
		return "", errors.New("@api format error: " + doc.Api.String())
	}
	if !regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`).MatchString(doc.ApiVersion) {
		return "", errors.New("@apiVersion must like be [x.y.z], but got: " + doc.ApiVersion)
	}
	if len(doc.ApiName) < 1 {
		return "", errors.New("@apiName format error: " + doc.ApiName)
	}
	if len(doc.ApiGroup) < 1 {
		return "", errors.New("@apiGroup format error: " + doc.ApiGroup)
	}

	resultLines = append(resultLines, fmt.Sprintf("@api %s", doc.Api))
	resultLines = append(resultLines, fmt.Sprintf("@apiVersion %s", doc.ApiVersion))
	resultLines = append(resultLines, fmt.Sprintf("@apiName %s", doc.ApiName))
	resultLines = append(resultLines, fmt.Sprintf("@apiGroup %s", doc.ApiGroup))

	if doc.ApiPrivate {
		resultLines = append(resultLines, "@apiPrivate true")
	}
	if len(doc.ApiPermission) > 0 {
		resultLines = append(resultLines, fmt.Sprintf("@apiPermission %s", doc.ApiPermission))
	}
	if len(doc.ApiDescription) > 0 {
		resultLines = append(resultLines, fmt.Sprintf("@apiDescription %s", doc.ApiDescription))
	}
	for _, def := range doc.ApiUse {
		resultLines = append(resultLines, fmt.Sprintf("@apiUse %s", def))
	}
	for _, def := range doc.ApiParam {
		resultLines = append(resultLines, fmt.Sprintf("@apiParam %s", def))
	}
	for _, def := range doc.ApiHeader {
		resultLines = append(resultLines, fmt.Sprintf("@apiHeader %s", def))
	}
	for _, def := range doc.ApiQuery {
		resultLines = append(resultLines, fmt.Sprintf("@apiQuery %s", def))
	}
	for _, def := range doc.ApiBody {
		resultLines = append(resultLines, fmt.Sprintf("@apiBody %s", def))
	}
	for _, def := range doc.ApiSuccess {
		resultLines = append(resultLines, fmt.Sprintf("@apiSuccess %s", def))
	}
	for _, def := range doc.ApiSuccessExample {
		resultLines = append(resultLines, fmt.Sprintf("@apiSuccessExample %s", def))
	}
	for _, def := range doc.ApiError {
		resultLines = append(resultLines, fmt.Sprintf("@apiError %s", def))
	}
	for _, def := range doc.ApiErrorExample {
		resultLines = append(resultLines, fmt.Sprintf("@respExampleDef %s", def))
	}
	result := fmt.Sprintf("/**\n\n%s\n\n*/", strings.Join(resultLines, "\n"))

	return result, nil
}

func (a routeDef) String() string {
	return fmt.Sprintf("{%s} %s %s", a.Method, a.Path, a.Title)
}

func (p paramDef) String() string {
	typ := "{" + p.Type.Type
	if len(p.Type.Size) > 0 {
		typ += "{" + p.Type.Size + "}"
	}
	if len(p.Type.Allows) > 0 {
		typ += "=" + strings.Join(p.Type.Allows, ",")
	}
	typ += "}"

	field := p.Field.Name
	if len(p.Field.Default) > 0 {
		field += "=" + p.Field.Default
	}
	if !p.Field.Required {
		field = fmt.Sprintf("[%s]", field)
	}

	return fmt.Sprintf("%s %s %s", typ, field, p.Description)
}

func (r respExampleDef) String() string {
	var formattedString string
	if len(r.Example) > 0 {
		var jsonData any
		err := json.Unmarshal([]byte(r.Example), &jsonData)
		if err == nil {
			if indentedJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				formattedString = string(indentedJSON)
			}
		}
	}

	return fmt.Sprintf("{%s} %s\n%s\n%s", r.Type, r.Title, r.Http, formattedString)
}
