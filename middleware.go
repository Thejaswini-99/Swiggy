package main

import (
	"fmt"
	"net/http"
	"strings"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorisation := r.Header.Get("Authorization")

		if authorisation == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "No Token Provided")
			return
		}

		tokenStr := strings.TrimPrefix(authorisation, "Bearer ")

		token, err := validateToken(tokenStr)
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "It is not a valid Token")
			return
		}

		next(w, r)

	}

}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Orign", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		next(w, r)
	}

}
