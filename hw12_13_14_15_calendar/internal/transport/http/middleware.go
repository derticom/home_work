package internalhttp

import (
	"net/http"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.serverLog.Info(r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}
