package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nanasuryana335/honda-leasing-api/internal/response"
)

type JWTClaims struct {
	UserID      int64
	Username    string
	Email       string
	Role        string
	IsActive    bool
	LockedUntil *time.Time
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendError(c, http.StatusUnauthorized, "mising auth token")
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) < 2 {
			response.SendError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
	}
}
