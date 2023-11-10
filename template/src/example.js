/**

@api {post} /api/example/:id 示例接口
@apiVersion 0.0.1
@apiName PostExample
@apiGroup Example
@apiPermission admin
@apiDescription 这是一个测试接口，用来测试文档生成
需要注意有换行的描述

@apiUse GlobalHeader

@apiParam {Number} id ID
@apiQuery {Number{1..}} [page=1] 页码

@apiBody {string} name 姓名
@apiBody {string{1..8}} nickname 用户名
@apiBody {number} [age] 年龄
@apiBody {string{2-5}="small","big"} statue 状态
@apiBody {String[]} area 地区
@apiBody {Object} [address] 地址
@apiBody {String{2-5}="XX","YY"} [address.street=ZZ] 街道
@apiBody {String} [address.zip] 邮编

@apiError {string} name 姓名
@apiErrorExample {json} NotFound:
HTTP/1.1 404 Not Found
{
    "error": "UserNotFound"
}
@apiErrorExample {json} Unauthorized:
HTTP/1.1 401 Unauthorized
{
    "error": "未授权"
}

@apiSuccess {Object}  profile  用户信息
@apiSuccess {Number} profile.age 年龄
@apiSuccess {String}  profile.image 头像
@apiSuccess {Object[]} profiles 列表
@apiSuccess {Number}   profiles.age 年龄
@apiSuccess {String}   profiles.image 头像
@apiSuccess {String}   url 链接
@apiSuccessExample {json} Success-Response:
HTTP/1.1 200 OK
{
    "code": 0,
    "data": {
        "avatar": "https://cdn.the.cn",
        "nickname": "HOYOO"
    },
    "meta": {
        "token": "yJhbGciOinR5cCI"
    }
}

*/

