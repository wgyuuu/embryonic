package middle

import (
	"net/http"

	"github.com/urfave/negroni"

	"github.com/wgyuuu/embryonic/helper/gls"
	"github.com/wgyuuu/embryonic/helper/util"
)

// 服务化http请求串联
func Context(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	traceID := r.Header.Get("header-rid")
	spanID := r.Header.Get("header-spanid")
	if spanID == "" {
		spanID = util.GenerateSpanID(r)
	}

	gls.SetGls(traceID, spanID, func() {
		next(rw, r)
	})
}

func NewContext() negroni.Handler {
	return negroni.HandlerFunc(Context)
}
