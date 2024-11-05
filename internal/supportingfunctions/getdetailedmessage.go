package supportingfunctions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// GetDetailedMessage возвращает детальное описание сообщения которое может содержатся
// в http.Body, а если такого сообщения нет то пустое значение
func GetDetailedMessage(body io.ReadCloser) (string, error) {
	var (
		b   []byte
		err error
	)

	defer func(body io.ReadCloser) {
		if errClose := body.Close(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}(body)

	var msg string

	b, err = io.ReadAll(body)
	if err != nil {
		return msg, err
	}

	var r map[string]interface{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return msg, err
	}

	if message, ok := r["message"]; ok {
		msg = fmt.Sprint(message)
	}

	return msg, nil
}
