package register

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/wgyuuu/embryonic/common/response"
	"github.com/wgyuuu/embryonic/helper/util"

	"github.com/julienschmidt/httprouter"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	PATCH
	DELETE
)

type HandlerData struct {
	methods []Method
	proc    *FuncData
	before  *FuncData
}

type Handler struct {
	Methods []Method
	Handle  httprouter.Handle
}

type HandlerMap struct {
	handlerMap map[string]*HandlerData
}

func NewHandlerMap() *HandlerMap {
	return &HandlerMap{
		handlerMap: make(map[string]*HandlerData),
	}
}

func (h *HandlerMap) Register(pattern string, handlerFunc, beforeFunc interface{}, methods ...Method) {
	if !checkRets(handlerFunc) {
		panic(fmt.Sprintf("%s return error", pattern))
	}

	if len(methods) == 0 {
		methods = append(methods, POST)
	}
	handlerData := HandlerData{methods: methods}

	args, resp := funcAnalysis(handlerFunc)
	handlerData.proc = newFuncData(handlerFunc, args, resp)

	if beforeFunc != nil {
		befArgs, befResp := funcAnalysis(beforeFunc)
		if len(befResp) != len(args) {
			panic(fmt.Sprintf("%s handler args != before return", pattern))
		}
		for idx := range args {
			if !args[idx].Equal(befResp[idx]) {
				panic(fmt.Sprintf("%s handler args != before return", pattern))
			}
		}
		handlerData.before = newFuncData(beforeFunc, befArgs, befResp)
	}

	h.handlerMap[pattern] = &handlerData
}

func (h *HandlerMap) HandleMap() map[string]Handler {
	hMap := make(map[string]Handler, len(h.handlerMap))
	for pattern, handler := range h.handlerMap {
		hMap[pattern] = Handler{Methods: handler.methods, Handle: create(handler)}
	}
	return hMap
}

func checkRets(f interface{}) bool {
	funcValue := reflect.ValueOf(f)
	// 参数检测 参数都要是指针
	for idx := 0; idx < funcValue.Type().NumIn(); idx++ {
		typ := funcValue.Type().In(idx)
		if typ.Kind() != reflect.Ptr {
			return false
		}
	}
	// 返回值检测
	if funcValue.Type().NumOut() != 1 {
		return false
	}
	typ := funcValue.Type().Out(0)
	switch reflect.New(typ).Interface().(type) {
	case *response.StdResp, *error:
		return true
	}
	return false
}

func create(handler *HandlerData) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer util.Recover()

		bePool := false
		proc := handler.proc
		var procArgs []reflect.Value
		var body string
		httpStatus := http.StatusOK
		defer func() {
			rw.WriteHeader(httpStatus)
			rw.Write([]byte(body))

			if !bePool || httpStatus != http.StatusOK {
				return
			}
			for idx, value := range procArgs {
				proc.PutObject(proc.Args[idx], value.Interface())
			}
		}()

		//  如果有预处理函数，就通过预处理函数生成proc函数的参数
		if before := handler.before; before != nil {
			args := make([]reflect.Value, len(before.Args))
			for idx, arg := range before.Args {
				switch arg.Kind() {
				case Request:
					args[idx] = reflect.ValueOf(r)
				case Params:
					args[idx] = reflect.ValueOf(ps)
				}
			}
			procArgs = reflect.ValueOf(before.F).Call(args)
		}

		// 自己创建参数，并在body里读取数据赋值
		if len(procArgs) == 0 && len(proc.Args) > 0 {
			bePool = true
			procArgs = make([]reflect.Value, len(proc.Args))
			for idx, arg := range proc.Args {
				procArgs[idx] = reflect.ValueOf(proc.NewObject(arg))
			}
		}

		// 参数赋值
		inBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if len(inBody) > 0 {
			err := json.Unmarshal(inBody, procArgs[0].Interface())
			if err != nil {
				httpStatus = http.StatusInternalServerError
				body = fmt.Sprintf("body format error(%s)", err.Error())
				return
			}
		}

		respList := reflect.ValueOf(proc.F).Call(procArgs)
		if len(respList) == 0 || respList[0].Kind() == reflect.Invalid {
			return
		}
		bodyFunc, num := getBodyFunc(proc.Resp[0], respList[0])
		if !bodyFunc.IsValid() {
			return
		}
		result := bodyFunc.Call(respList[:num])
		body = result[0].Interface().(string)
	}

}
