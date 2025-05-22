package middleware

import (
	"log"
	"net/http"
	"time"
)

func Chain(middlewares ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// Allows the user to modify CORS rules,
// by default Umi follows a zero trust approach,
// this means that it will block any request from
// different origins
// This function expects the user to specify each part of its Cors rules
func Cors(origin, credentials, headers, methods string) Middleware {
	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", credentials)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Allow-Methods", methods)

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	return cors
}

// Defaults to Origin *
func FlexibleCors() Middleware {
	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", FLEXIBLE_ORIGIN)
			w.Header().Set("Access-Control-Allow-Credentials", FLEXIBLE_COR_CREDENTIALS)
			w.Header().Set("Access-Control-Allow-Headers", FLEXIBLE_COR_HEADERS)
			w.Header().Set("Access-Control-Allow-Methods", FLEXIBLE_COR_METHODS)

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	return cors
}

func Logger() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request at %v", time.Now())
			log.Printf("Requested from %s", r.RemoteAddr)
			next(w, r)
		}
	}
}
