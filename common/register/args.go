package register

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/wgyuuu/embryonic/common/response"
)

type Kind int

const (
	Undefine Kind = iota
	Request
	Params
)

const (
	StringFunc = "String"
	ErrorFunc  = "Error"
	String     = "string"
	JsonFunc   = "Json"
)

var (
	typeRequest reflect.Type = reflect.TypeOf((*http.Request)(nil))
	typeParams  reflect.Type = reflect.TypeOf((httprouter.Params)(nil))
)

type Args struct {
	kind     Kind
	typ      reflect.Type
	bodyFunc string
}

func newArgs(typ reflect.Type) *Args {
	var kind Kind
	switch typ {
	case typeRequest:
		kind = Request
	case typeParams:
		kind = Params
	default:
		kind = Undefine
	}

	var bodyFunc string
	switch reflect.New(typ).Interface().(type) {
	case *response.StdResp:
		bodyFunc = StringFunc
	case *error:
		bodyFunc = ErrorFunc
	case *string:
		bodyFunc = String
	default:
		bodyFunc = JsonFunc
	}

	return &Args{
		typ:      typ,
		kind:     kind,
		bodyFunc: bodyFunc,
	}
}

func (a *Args) Equal(b *Args) bool {
	return a.typ == b.typ
}

func (a *Args) Typ() reflect.Type {
	return a.typ
}

func (a *Args) Type() string {
	return a.typ.String()
}

func (a *Args) Kind() Kind {
	return a.kind
}

func (a *Args) BodyFunc() string {
	return a.bodyFunc
}

func getBodyFunc(args *Args, value reflect.Value) (valueFunc reflect.Value, num int) {
	switch args.BodyFunc() {
	case StringFunc, ErrorFunc:
		valueFunc = value.MethodByName(args.BodyFunc())
	case String:
		valueFunc, num = reflect.ValueOf(FuncString), 1
	case JsonFunc:
		valueFunc, num = reflect.ValueOf(FuncJson), 1
	}
	return
}

func FuncString(body string) string {
	return body
}

func FuncJson(body interface{}) string {
	bytes, _ := json.Marshal(body)
	return string(bytes)
}
