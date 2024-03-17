package encoding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type encoding interface {
	encoder
	decoder
}

// encoder interface is used to encode the response
type encoder interface {
	Encode(v any) error
}

// decoder interface is used to decode the request
type decoder interface {
	Decode(v any) error
}

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	// Create a buffer to capture the encoded JSON
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		// If encoding fails, set the status header to indicate an error
		w.WriteHeader(http.StatusInternalServerError)
		return newEncodingError(EncodeJSON, err, "failed to encode JSON.")
	}

	// Write the buffer to the response writer
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		return newEncodingError(WriteResponse, err, "failed to write buffer to response writer.")
	}

	return nil
}

func Decode[T any](body io.ReadCloser) (T, error) {
	var v T
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		return v, newEncodingError(DecodeJSON, err, "failed to decode JSON.")
	}
	return v, nil
}
