package eventsource

import (
	"encoding/json"
	"fmt"
	"io"
)

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
