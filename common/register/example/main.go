package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	"github.com/wgyuuu/embryonic/app/middle"
	"github.com/wgyuuu/embryonic/common/register"
	"github.com/wgyuuu/embryonic/common/register/example/epmiddle"
	"github.com/wgyuuu/embryonic/common/response"
	"github.com/wgyuuu/embryonic/helper/dlog"
	"github.com/wgyuuu/embryonic/helper/signal"
)

func main() {
	signal.OpenMonitor(func() {
		dlog.Strace(epmiddle.AvgTimeInfo())
	})

	hMap := register.NewHandlerMap()
	hMap.Register("/test", TestHandler, TestBefore)
	hMap.Register("/newdata", NewDataHandler, nil)

	ng := negroni.New()
	ng.Use(middle.NewContext())
	ng.Use(epmiddle.NewAccessLogger())

	router := httprouter.New()
	for pattern, handler := range hMap.HandleMap() {
		for _, method := range handler.Methods {
			switch method {
			case register.GET:
				router.GET(pattern, handler.Handle)
			case register.POST:
				router.POST(pattern, handler.Handle)
			case register.PUT:
				router.PUT(pattern, handler.Handle)
			case register.PATCH:
				router.PATCH(pattern, handler.Handle)
			case register.DELETE:
				router.DELETE(pattern, handler.Handle)
			}
		}

	}
	ng.UseHandler(router)

	err := http.ListenAndServe(":8888", ng)
	if err != nil {
		panic(fmt.Sprintf("ListenAndServe error: %s", err.Error()))
	}
}

type Test struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type NewData struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func NewDataHandler(d *NewData) response.StdResp {
	d.A += 10
	return response.NewStdResponse(d)
}

func TestHandler(t *Test) response.StdResp {
	return response.NewStdResponse(t)
}

//func TestHandler(t *Test) error {
//return errors.New("test")
//}

func TestBefore(r *http.Request, ps httprouter.Params) *Test {
	return &Test{1, "test"}
}
