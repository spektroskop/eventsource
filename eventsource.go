package eventsource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Flusher); !ok {
			http.Error(w, "ResponseWriter is not a http.Flusher", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")

		next.ServeHTTP(w, r)
	})
}

func Flush(w http.ResponseWriter) {
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func SendEvent(w io.Writer, name string, data []byte) error {
	_, err := w.Write([]byte(
		fmt.Sprintf("event: %s\ndata: %s\n\n", name, string(data)),
	))
	return err
}

func SendEventf(w io.Writer, data []byte, f string, args ...interface{}) error {
	return SendEvent(w, fmt.Sprintf(f, args...), data)
}

func EncodeEvent(w io.Writer, name string, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return SendEvent(w, name, data)
}

func EncodeEventf(w io.Writer, object interface{}, f string, args ...interface{}) error {
	return EncodeEvent(w, fmt.Sprintf(f, args...), object)
}

func SendMessage(w io.Writer, data []byte) error {
	_, err := w.Write([]byte(
		fmt.Sprintf("data: %s\n\n", string(data)),
	))
	return err
}

func EncodeMessage(w io.Writer, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return SendMessage(w, data)
}
