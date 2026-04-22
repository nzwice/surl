package transport

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/nzwice/surl/pkg/ctxutils"
)

const (
	xRequestIdHeader string = "x-request-id"
)

func withRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get(xRequestIdHeader)
		if reqId == "" {
			reqId = uuid.NewString()
		}
		ctx := ctxutils.WithRequestId(r.Context(), reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := newProxyWriter(w).(*proxyWriter)
		defer func(start time.Time) {
			end := time.Now()
			slog.InfoContext(r.Context(), "request log",
				slog.Time("start_time", start.UTC()),
				slog.Time("end_time", end.UTC()),
				slog.String("user_agent", r.Header.Get("User-Agent")),
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("referer", r.Referer()),
				slog.String("ip", r.RemoteAddr),
				slog.Duration("duration", end.Sub(start)),
				slog.Int("status", pw.StatusCode),
			)
		}(time.Now())
		next.ServeHTTP(pw, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				stack := debug.Stack()
				slog.ErrorContext(r.Context(), "panic detected", slog.Any("stack", string(stack)), slog.Any("error", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type proxyWriter struct {
	next       http.ResponseWriter
	StatusCode int
}

// Header implements [http.ResponseWriter].
func (p *proxyWriter) Header() http.Header {
	return p.next.Header()
}

// Write implements [http.ResponseWriter].
func (p *proxyWriter) Write(b []byte) (int, error) {
	if p.StatusCode == 0 {
		p.StatusCode = http.StatusOK
	}
	return p.next.Write(b)
}

// WriteHeader implements [http.ResponseWriter].
func (p *proxyWriter) WriteHeader(statusCode int) {
	p.StatusCode = statusCode
	p.next.WriteHeader(statusCode)
}

func newProxyWriter(next http.ResponseWriter) http.ResponseWriter {
	return &proxyWriter{
		next: next,
	}
}
