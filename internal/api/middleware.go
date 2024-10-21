package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/negroni"
)

func (h *Handler) LogAPI(handler http.Handler) http.Handler {
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

func (h *Handler) gzipMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.loger.Debugln("gzipMiddleware")

		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.loger.Errorln(err)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		handler.ServeHTTP(ow, r)
	})
}
