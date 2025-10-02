package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.RWMutex
	RequestsTotal int64
	ErrorsTotal   int64
	PerPath       map[string]int64
}

var AppMetrics = &Metrics{PerPath: map[string]int64{}}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("← %s %s %d %s", r.Method, r.URL.Path, lrw.status, duration)
		AppMetrics.mu.Lock()
		AppMetrics.RequestsTotal++
		AppMetrics.PerPath[r.URL.Path]++
		if lrw.status >= 400 {
			AppMetrics.ErrorsTotal++
		}
		AppMetrics.mu.Unlock()
	})
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	AppMetrics.mu.RLock()
	defer AppMetrics.mu.RUnlock()
	_ = json.NewEncoder(w).Encode(AppMetrics)
}
