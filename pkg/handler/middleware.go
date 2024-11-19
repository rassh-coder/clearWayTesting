package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

const authHeader = "Authorization"
const typeAuthHeader = "Bearer"
const expTime = time.Hour * 24
const ctxtUidKey = "userId"

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tokenHeader := r.Header.Get(authHeader)
		tokenArr := strings.Split(tokenHeader, " ")

		if len(tokenArr) < 2 {
			errRes := ErrorResponse{Error: "unauthorized"}
			h.NewError(w, &errRes, http.StatusUnauthorized)
			return
		}

		if tokenArr[0] != typeAuthHeader {
			errRes := ErrorResponse{Error: "unauthorized"}
			h.NewError(w, &errRes, http.StatusUnauthorized)
			return
		}

		token := tokenArr[1]

		session, err := h.Services.Authorization.GetSessionByToken(token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errRes := ErrorResponse{Error: "unauthorized"}
				h.NewError(w, &errRes, http.StatusUnauthorized)
				return
			}
			log.Printf("Can't get session: %s", err)

			errRes := ErrorResponse{Error: err.Error()}
			h.NewError(w, &errRes, http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), ctxtUidKey, session.UID)
		r = r.WithContext(ctx)

		dif := time.Since(*session.CreatedAt)

		if !session.IsActive || dif > expTime {
			errRes := ErrorResponse{Error: "unauthorized"}
			h.NewError(w, &errRes, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
