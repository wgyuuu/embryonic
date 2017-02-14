package handler

import "github.com/wgyuuu/embryonic/common/register"

var handlerMap *register.HandlerMap = register.NewHandlerMap()

func GetHandleMap() map[string]register.Handler {
	return handlerMap.HandleMap()
}
