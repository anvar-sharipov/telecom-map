package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/anvar-sharipov/telecom-map/internal/utils"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func Auth(next func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		// 1️⃣ Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return utils.NewUnauthorized("missing Authorization header")
		}

		// 2️⃣ Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.NewUnauthorized("invalid Authorization header")
		}

		tokenString := parts[1]

		// 3️⃣ Parse & validate ACCESS token
		userID, err := utils.ParseToken(tokenString)
		if err != nil {
			return utils.NewUnauthorized("invalid or expired access token")
		}

		// 4️⃣ Save user_id to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// 5️⃣ Next
		return next(w, r.WithContext(ctx))
	}
}
