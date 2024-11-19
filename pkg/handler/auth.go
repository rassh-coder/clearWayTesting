package handler

import (
	"clearWayTest/pkg/models"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

type Resp struct {
	Status string `json:"status"`
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	var authInputs models.UserInputs
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.NewNotFound(w)
		return
	}

	errRes := DecodeJSONBody(w, r, &authInputs)
	if errRes != nil {
		h.NewError(w, errRes, http.StatusBadRequest)
		return
	}
	if authInputs.Login == nil {
		errRes := ErrorResponse{Error: "login is required"}
		h.NewError(w, &errRes, http.StatusBadRequest)
		return
	}
	if authInputs.Password == nil {
		errRes := ErrorResponse{Error: "password is required"}
		h.NewError(w, &errRes, http.StatusBadRequest)
		return
	}

	ip, err := h.ReadUserIP(r)
	if err != nil {
		log.Printf("Can't read user ip")
	}

	token, err := h.Services.Authorization.SignIn(*authInputs.Login, *authInputs.Password, ip)

	if err != nil {
		if err.Error() == "unauthorized" {
			errRes := ErrorResponse{Error: "invalid login/password"}
			h.NewError(w, &errRes, http.StatusUnauthorized)
			return
		}

		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	res := TokenResponse{Token: token}

	jsonRes, err := json.Marshal(&res)
	if err != nil {
		log.Printf("Can't marshall response to json: %s", err)
		http.Error(w, "can't response json body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonRes)
	if err != nil {
		log.Printf("Can't write response body: %s", err)
		http.Error(w, "can't write response", http.StatusInternalServerError)
		return
	}

	return
}

func (h *Handler) ReadUserIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
