package main

import (
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("Got the header as : " + authHeader)

		// // Extract the token from the header
		// tokenParts := strings.Split(authHeader, " ")
		// if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		token := authHeader

		log.Printf("Validating token : " + token)
		// Validate the token using Cognito
		valid, err := ValidateToken(token)
		log.Printf("Got output : %v ", valid)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
