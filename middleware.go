package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiters(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		host = ip
	}
	if limiter, ok := limiters[host]; ok {
		return limiter
	}

	limiter := rate.NewLimiter(1, 1)
	limiters[host] = limiter

	return limiter

}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := getLimiters(r.RemoteAddr)

		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Too many requests")
			return
		}
		next(w, r)

	}

}

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
		fmt.Println("err:", err) // ← add this
		fmt.Println("valid:", token.Valid)
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
