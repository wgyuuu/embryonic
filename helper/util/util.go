package util

import (
	"encoding/binary"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wgyuuu/embryonic/helper/dlog"
)

func GenerateSpanID(r *http.Request) string {
	strAddr := strings.Split(r.RemoteAddr, ":")
	ipLong, _ := Ip2Long(strAddr[0])
	time := uint64(time.Now().UnixNano())

	spanId := ((time ^ uint64(ipLong)) << 32) | uint64(rand.Int31())
	return strconv.FormatUint(spanId, 16)
}

func Ip2Long(ip string) (uint32, error) {
	addr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(addr.IP.To4()), nil
}

func SafeFunc(f func()) {
	defer func() {
		Recover()
	}()

	f()
}

func Recover() {
	if rcv := recover(); rcv == nil {
		return
	} else {
		dlog.Error("panic", rcv)
	}

	i := 0
	funcPtr, file, line, ok := runtime.Caller(i)
	for ok {
		dlog.Error("func", runtime.FuncForPC(funcPtr).Name(), "file", file, "line", line)
		i++
		funcPtr, file, line, ok = runtime.Caller(i)
	}
}
