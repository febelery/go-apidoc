package goapidoc

import (
	"testing"
)

func TestNewDoc(t *testing.T) {
	api := ApiDef{
		Method:       "post",
		Path:         "/api/example/:id",
		Params:       nil,
		Headers:      nil,
		Queries:      nil,
		StatusCode:   200,
		ResponseBody: []byte(`{ "code": 0, "data": { "url": "https://cdn.the.cn", "nickname": "yyy" }, "meta": { "token": "yJhbGciOinR5cCI" } }`),
	}

	New(api)
}
