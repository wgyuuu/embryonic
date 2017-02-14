package gls

import "github.com/jtolds/gls"

const (
	keyTraceID = "trace_id"
	keySpanID  = "span_id"
)

var mgr *gls.ContextManager

func init() {
	mgr = gls.NewContextManager()
}

func SetGls(traceID, spanID string, contextCall func()) {
	mgr.SetValues(gls.Values{
		keyTraceID: traceID,
		keySpanID:  spanID,
	}, contextCall)
}

func GetGlsInfo() (traceID, spanID string) {
	traceObj, _ := mgr.GetValue(keyTraceID)
	spanObj, _ := mgr.GetValue(keySpanID)

	traceID, _ = traceObj.(string)
	spanID, _ = spanObj.(string)
	return
}
