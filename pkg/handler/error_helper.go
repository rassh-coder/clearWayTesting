package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) NewError(w http.ResponseWriter, errRes *ErrorResponse, status int) {
	w.WriteHeader(status)
	jsonErrRes, err := json.Marshal(&errRes)

	if err != nil {
		log.Printf("Can't marshall response to json: %s", err)
		http.Error(w, "can't response json body", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonErrRes)
	if err != nil {
		log.Printf("Can't write response body: %s", err)
		http.Error(w, "can't write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewNotFound(w http.ResponseWriter) {
	errRes := ErrorResponse{Error: "not found"}
	h.NewError(w, &errRes, http.StatusNotFound)
}
