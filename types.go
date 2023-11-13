package goapidoc

type routeDef struct {
	Method string // get,post,put,delete
	Path   string
	Title  string
}

// eg: {String{2-5}="small","big"}
type typeDef struct {
	Type   string   // 类型： {Boolean}, {Number}, {String}, {Object}, {String[]}
	Size   string   // 取值范围
	Allows []string // 允许值
}

// eg: [address.street=ZZZ]
type fieldDef struct {
	Name     string
	Required bool
	Default  string
}

// eg: {String{2-5}="XXX","YYY"} [address.street=ZZZ] 街道
type paramDef struct {
	Type        typeDef
	Field       fieldDef
	Description string
}

type respExampleDef struct {
	Type    string
	Title   string
	Http    string
	Example string
}

type docDef struct {
	Api               routeDef         `name:"@api"`                      // @api {put} /user/:id 修改用户
	ApiVersion        string           `name:"@apiVersion"`               // @apiVersion 0.0.1
	ApiName           string           `name:"@apiName"`                  // @apiName UpdateUser
	ApiGroup          string           `name:"@apiGroup"`                 // @apiGroup User
	ApiPermission     string           `name:"@apiPermission,omitempty"`  // @apiPermission admin
	ApiPrivate        bool             `name:"@apiPrivate,omitempty"`     // @apiPrivate true
	ApiDescription    string           `name:"@apiDescription,omitempty"` // @apiDescription
	ApiUse            []string         `name:"@apiUse,omitempty"`
	ApiDefine         []string         `name:"@apiDefine,omitempty"`  // @apiDefine Common(marshal忽略此字段)
	ApiParam          []paramDef       `name:"@apiParam,omitempty"`   // @apiParam {Number} id 用户ID
	ApiHeader         []paramDef       `name:"@apiHeader,omitempty"`  // @apiHeader {String} Authorization="Bearer xxx" 授权token
	ApiQuery          []paramDef       `name:"@apiQuery,omitempty"`   // @apiQuery {Number{1..}} [page=1] 页码
	ApiBody           []paramDef       `name:"@apiBody,omitempty"`    // @apiBody {String{2-5}="XXX","YYY"} [address.street=ZZZ] 街道
	ApiSuccess        []paramDef       `name:"@apiSuccess,omitempty"` // @apiSuccess {Object[]} profiles 列表
	ApiError          []paramDef       `name:"@apiError,omitempty"`   // @apiError UserNotFound 用户未找到
	ApiSuccessExample []respExampleDef `name:"@apiSuccessExample,omitempty"`
	ApiErrorExample   []respExampleDef `name:"@apiErrorExample,omitempty"`
}

type docsDef []*docDef

type ApiDef struct {
	Method       string
	Path         string
	Params       map[string]string
	Headers      map[string]string
	Queries      map[string]string
	StatusCode   int
	ResponseBody []byte
	Skip         func() bool
	group        string
	version      string
}
