package json

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func ReadJson(w http.ResponseWriter, r *http.Request, v interface{}) error {
	const maxBytes = 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is empty")
		}
		return err
	}

	if decoder.More() {
		return errors.New("request body contains unexpected extra data")
	}

	return nil
}

func WriteJson(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	status := "error"
	if code >= 200 && code < 300 {
		status = "success"
	}
	w.WriteHeader(code)
	if status == "error" {
		response := ErrResponse{
			Code:    code,
			Message: v,
		}
		return json.NewEncoder(w).Encode(response)
	} else {
		response := Response{
			Code:   code,
			Data:   v,
			Status: status,
		}
		return json.NewEncoder(w).Encode(response)
	}
}

func WriteJsonError(w http.ResponseWriter, status int, mssg string) {
	if err := WriteJson(w, status, mssg); err != nil {
		panic(err)
	}
}
