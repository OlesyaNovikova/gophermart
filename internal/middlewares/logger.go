package middlewares

import (
	"net/http"
	"time"

	l "github.com/OlesyaNovikova/gophermart/internal/utils/logger"
)

type responseData struct {
	status int
	size   int
}

// loggingWriter реализует интерфейс http.ResponseWriter
type loggingWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	if r.responseData.status == 0 {
		r.responseData.status = http.StatusOK
	}
	return size, err
}

func (r *loggingWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLog(next http.Handler) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		resData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   resData,
		}
		next.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		l.Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", resData.status,
			"size", resData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
