package middleware

import (
	"context"
	"github.com/drTragger/MykroTask/utils"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type contextKey string

type JwtClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

const UserIDKey contextKey = "userID"

func JWTMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSONResponse(w, http.StatusUnauthorized, &utils.ErrorResponse{
					Status:  false,
					Message: "Missing Authorization header",
				})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.WriteJSONResponse(w, http.StatusUnauthorized, &utils.ErrorResponse{
					Status:  false,
					Message: "Invalid Authorization header format",
				})
				return
			}

			tokenStr := parts[1]
			claims := &JwtClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				utils.WriteJSONResponse(w, http.StatusUnauthorized, &utils.ErrorResponse{
					Status:  false,
					Message: "Invalid token",
					Errors:  err.Error(),
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
