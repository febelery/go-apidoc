package goapidoc

import (
	"fmt"
	"net/http"
	"strings"
)

// generateWithOldDoc todo doc.ApiPermission | doc.ApiPrivate | doc.ApiUse
func (doc *docDef) generateWithOldDoc(api ApiDef, olderDoc docDef) {
	doc.ApiGroup = api.group
	doc.ApiVersion = api.version
	doc.Api.Method = strings.ToLower(api.Method)
	doc.Api.Path = api.Path

	if len(olderDoc.Api.Title) > 0 {
		doc.Api.Title = olderDoc.Api.Title
	} else {
		doc.Api.Title = fmt.Sprintf("%s %s", strings.ToUpper(api.Method), api.Path)
	}

	if len(olderDoc.ApiDescription) > 0 {
		doc.ApiDescription = olderDoc.ApiDescription
	} else {
		doc.ApiDescription = fmt.Sprintf("%s %s ", strings.ToUpper(api.Method), api.Path)
	}

	if len(olderDoc.ApiName) > 0 {
		doc.ApiName = olderDoc.ApiName
	} else {
		doc.ApiName = fmt.Sprintf("%s-%s", api.Method, api.Path)
	}

	for k, v := range api.Params {
		doc.appendParam(&doc.ApiParam, generateParamDef(k, v, ""), olderDoc.ApiParam)
	}

	for k, v := range api.Headers {
		doc.appendParam(&doc.ApiHeader, generateParamDef(k, strings.Join(v, " "), ""), olderDoc.ApiHeader)
	}

	for k, v := range api.Queries {
		doc.appendParam(&doc.ApiQuery, generateParamDef(k, v, ""), olderDoc.ApiQuery)
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
			doc.appendParam(&doc.ApiSuccess, param, olderDoc.ApiSuccess)
		} else {
			doc.appendParam(&doc.ApiError, param, olderDoc.ApiError)
		}
	}

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
