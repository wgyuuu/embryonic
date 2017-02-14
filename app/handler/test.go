package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wgyuuu/embryonic/common/register"
	"github.com/wgyuuu/embryonic/common/response"
)

type Test struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func init() {
	handlerMap.Register("/test", HandleTest, BeforeTest, register.GET, register.POST)
}

//  参数要为指针
func HandleTest(t *Test) response.StdResp { // 返回值 也支持error
	return response.NewStdResponse(t)
}

// 参数为*http.Request、 httprouter.Params两种类型,可以取任意一个或两个,顺序无限制,也支持无参
// 返回值必须和handle方法的参数保持一致，会直接作为handle的参数
// before方法可以不实现，注册传nil
// 系统会自动读取http请求里的body用json方式复制给handle结构体，如果额外需求则在before里实现
func BeforeTest(r *http.Request, ps httprouter.Params) *Test {
	return &Test{
		A: 1,
		B: "test",
	}
}
