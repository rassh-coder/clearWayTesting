package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *ErrorResponse {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &ErrorResponse{Error: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		if err.Error() == "http: request body too large" {
			msg := "request body too large"
			return &ErrorResponse{Error: msg}
		}

		msg := "invalid json body"
		log.Printf("Can't decode json: %s", err)
		return &ErrorResponse{Error: msg}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "request body must only contain a single JSON object"
		return &ErrorResponse{Error: msg}
	}

	return nil
}
