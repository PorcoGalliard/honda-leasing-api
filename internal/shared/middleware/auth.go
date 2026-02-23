package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRoles  = "roles"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Format Authorization tidak valid, gunakan: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Token claims tidak valid")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Token tidak mengandung user_id")
			c.Abort()
			return
		}
		c.Set(ContextKeyUserID, int64(userIDFloat))

		if rolesRaw, exists := claims["roles"]; exists {
			if rolesSlice, ok := rolesRaw.([]interface{}); ok {
				roles := make([]string, 0, len(rolesSlice))
				for _, r := range rolesSlice {
					if s, ok := r.(string); ok {
						roles = append(roles, s)
					}
				}
				c.Set(ContextKeyRoles, roles)
			}
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}

func GetRoles(c *gin.Context) ([]string, bool) {
	val, exists := c.Get(ContextKeyRoles)
	if !exists {
		return nil, false
	}
	roles, ok := val.([]string)
	return roles, ok
}
