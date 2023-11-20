# go-apidoc

API 文档生成
该项目使用 `apidoc` 工具自动生成 API 文档。通过自动抓取在 Go 语言的 API 接口上下文生成符合 `apidoc`
要求的文档。

## 安装

```shell
go get -u github.com/febelery/go-apidoc
npm install -g apidoc
```

## 使用

- gofiber 框架
```golang
app.Static("/doc", "doc/api")
app.Static("/assets", "doc/api/assets")

app.Use(func(c *fiber.Ctx) error {
  err := c.Next()
  skip := string(c.Response().Header.Peek("Content-Type")) != "application/json" ||
            strings.ToLower(os.Getenv("SKIP_DOC")) == "true"

  go goapidoc.New(goapidoc.ApiDef{
    Method:       c.Method(),
    Path:         c.Route().Path,
    Params:       c.AllParams(),
    Headers:      nil,
    Queries:      c.Queries(),
    StatusCode:   c.Response().StatusCode(),
    ResponseBody: c.Response().Body(),
    Skip:         skip,
  })

  return err
})
```

## 依赖

- Go >= 1.21
- Node.js
- apidoc：https://github.com/apidoc/apidoc

