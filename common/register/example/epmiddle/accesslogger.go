package epmiddle

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/urfave/negroni"
	"github.com/wgyuuu/embryonic/helper/dlog"
)

var totalCnt int
var totalTime time.Duration
var mutex sync.RWMutex

func AvgTimeInfo() string {
	if totalCnt == 0 {
		return "accessLogger"
	}
	return fmt.Sprintf("accessLogger, TotalTime=%v||TotalCnt=%d||AvgTime=%v\n", totalTime, totalCnt, totalTime/time.Duration(totalCnt))
}

func AccessLogger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	defer func() {
		mutex.Lock()
		totalCnt++
		totalTime += time.Since(start)
		mutex.Unlock()

		dlog.Info("access", "url", r.URL.Path, "time", time.Since(start))
	}()

	next(rw, r)
}

func NewAccessLogger() negroni.Handler {
	return negroni.HandlerFunc(AccessLogger)
}
