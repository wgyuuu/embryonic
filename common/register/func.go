package register

import (
	"reflect"
	"sync"

	"github.com/wgyuuu/embryonic/common/pool"
)

type FuncData struct {
	F          interface{}
	Args       []*Args
	Resp       []*Args
	factoryMap map[string]*sync.Pool // 是pool对象，暂时用interface代替
}

func newFuncData(f interface{}, args, resp []*Args) *FuncData {
	funcData := FuncData{
		F:          f,
		Args:       args,
		Resp:       resp,
		factoryMap: make(map[string]*sync.Pool),
	}
	for _, t := range args {
		funcData.factoryMap[t.Type()] = pool.NewPool(func() interface{} {
			return newObject(t.Typ())
		})
	}

	return &funcData
}

func (f *FuncData) NewObject(a *Args) interface{} {
	pool, ok := f.factoryMap[a.Type()]
	if ok {
		obj := pool.Get()
		return obj
	}

	return newObject(a.Typ())
}

func (f *FuncData) PutObject(a *Args, obj interface{}) {
	pool, ok := f.factoryMap[a.Type()]
	if !ok {
		return
	}

	pool.Put(obj)
}

func newObject(typ reflect.Type) interface{} {
	var value reflect.Value
	switch typ.Kind() {
	case reflect.Ptr:
		value = reflect.New(typ.Elem())
	default:
		value = reflect.New(typ).Elem()
	}

	return value.Interface()
}

func funcAnalysis(f interface{}) (args, resp []*Args) {
	funcType := reflect.ValueOf(f).Type()
	if funcType.Kind() != reflect.Func {
		return
	}

	args = make([]*Args, funcType.NumIn())
	resp = make([]*Args, funcType.NumOut())
	for index := range args {
		args[index] = newArgs(funcType.In(index))
	}
	for index := range resp {
		resp[index] = newArgs(funcType.Out(index))
	}
	return
}
