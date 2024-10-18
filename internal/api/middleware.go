package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/negroni"
)

func (h *APIHandler) LogAPI(handler http.Handler) http.Handler {
	h.loger.Debugln("LogAPI middleware")
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		uri := request.RequestURI
		method := request.Method
		reqID := middleware.GetReqID(request.Context())
		lrw := negroni.NewResponseWriter(writer)

		handler.ServeHTTP(lrw, request)

		statusCode := lrw.Status()
		duration := time.Since(start).Milliseconds()

		h.loger.Infoln("Request received: ",
			"requestID:", reqID,
			"uri:", uri,
			"method:", method,
			"duration:", duration,
		)
		h.loger.Infoln("Response sent: ",
			"requestID:", reqID,
			"statusCode:", statusCode,
			"size:", lrw.Size(),
		)
	}
	return http.HandlerFunc(logFn)
}
