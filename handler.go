package eventsource

import (
	"context"
	"net/http"
)

func Flush(w http.ResponseWriter) {
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.CloseNotifier); !ok {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if _, ok := w.(http.Flusher); !ok {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")

		ctx, cancel := context.WithCancel(r.Context())
		go func() {
			<-w.(http.CloseNotifier).CloseNotify()
			cancel()
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
