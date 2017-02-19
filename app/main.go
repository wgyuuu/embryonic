package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"github.com/wgyuuu/embryonic/app/handler"
	"github.com/wgyuuu/embryonic/app/middle"
	"github.com/wgyuuu/embryonic/common/register"
	"github.com/wgyuuu/embryonic/helper/signal"
)

func main() {
	signal.OpenMonitor(exitHandler)

	ng := negroni.New()

	ng.Use(middle.NewContext())
	ng.Use(middle.NewAccessLogger())

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(handler.NotFound)
	router.MethodNotAllowed = http.HandlerFunc(handler.MethodNotAllowed)

	for pattern, handler := range handler.GetHandleMap() {
		for _, method := range handler.Methods {
			insertHandle(router, pattern, method, handler.Handle)
		}

	}
	ng.UseHandler(router)

	err := http.ListenAndServe(":8080", ng)
	if err != nil {
		panic(fmt.Sprintf("ListenAndServe error: %s", err.Error()))
	}
}

// 进程结束时，被调用
func exitHandler() {
}

func insertHandle(router *httprouter.Router, pattern string, method register.Method, handle httprouter.Handle) {
	switch method {
	case register.GET:
		router.GET(pattern, handle)
	case register.POST:
		router.POST(pattern, handle)
	case register.PUT:
		router.PUT(pattern, handle)
	case register.PATCH:
		router.PATCH(pattern, handle)
	case register.DELETE:
		router.DELETE(pattern, handle)
	}
}
