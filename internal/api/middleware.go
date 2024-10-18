package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/negroni"
)

func (api *ApiHandler) LogAPI(h http.Handler) http.Handler {
	logFn := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		uri := request.RequestURI
		method := request.Method
		reqID := middleware.GetReqID(request.Context())
		lrw := negroni.NewResponseWriter(writer)

		h.ServeHTTP(lrw, request)

		statusCode := lrw.Status()
		duration := time.Since(start).Milliseconds()

		api.loger.Infoln("Request received: ",
			"requestID:", reqID,
			"uri:", uri,
			"method:", method,
			"duration:", duration,
		)
		api.loger.Infoln("Response sent: ",
			"requestID:", reqID,
			"statusCode:", statusCode,
			"size:", lrw.Size(),
		)
	}
	return http.HandlerFunc(logFn)
}
