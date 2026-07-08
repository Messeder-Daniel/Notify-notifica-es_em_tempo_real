package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			ctx.Abort()
			return
		}

		expiresAt, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(expiresAt) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
			ctx.Abort()
			return
		}

		userEmail, _ := claims["email"].(string)

		ctx.Set("userID", userID)
		ctx.Set("userEmail", userEmail)

		ctx.Next()
	}
}
