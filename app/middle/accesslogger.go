package middle

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"github.com/wgyuuu/embryonic/helper/dlog"
)

func AccessLogger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	defer func() {
		dlog.Info("access", "url", r.URL.Path, "time", time.Since(start))
	}()

	next(rw, r)
}

func NewAccessLogger() negroni.Handler {
	return negroni.HandlerFunc(AccessLogger)
}
