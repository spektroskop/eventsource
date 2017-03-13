package eventsource

import (
	"encoding/json"
	"fmt"
	"io"
)

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
