package main

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement authentication middleware
		// Verify the token from the Authorization header
		// If valid, call next.ServeHTTP(w, r)
		// If invalid, return an error response
		next.ServeHTTP(w, r)
	})
}
