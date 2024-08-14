package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func tokenFromRequest(r *http.Request) (string, error) {
	val := r.Header.Get("Authorization")
	splitted := strings.Split(val, " ")

	var t string

	if len(splitted) != 2 {
		return t, errors.New("invalid header structure")
	}

	if strings.ToUpper(splitted[0]) != "BEARER" {
		return t, errors.New("invalid header structure")
	}

	return splitted[1], nil
}

type contextKey string

const (
	userCtx contextKey = "user"
	orgCtx  contextKey = "org"
)

func requireAuthentication(logger *logrus.Entry,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", r.Context().Value(middleware.RequestIDKey).(string))
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }

func jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
