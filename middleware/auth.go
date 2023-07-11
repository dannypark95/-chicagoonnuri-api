package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

var UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

    fields := strings.Fields(authHeader)
    if len(fields) != 2 || fields[0] != "Bearer" {
        http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
        return
    }


		tokenString := fields[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// The middleware has no knowledge of the user details now, as we don't decode the claims.
		// The next handler will need to perform any user-specific logic.
		next.ServeHTTP(w, r)
	})
}