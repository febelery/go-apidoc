package goapidoc

import (
	"fmt"
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

	docs := NewDoc(api)
	for _, doc := range docs {
		fmt.Println(marshal(*doc))
	}
}
