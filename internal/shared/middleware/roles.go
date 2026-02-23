package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

var StaffRoles = map[string]bool{
	"SUPER_ADMIN":  true,
	"ADMIN_CABANG": true,
	"SALES":        true,
	"SURVEYOR":     true,
	"FINANCE":      true,
	"COLLECTION":   true,
}

// Middleware untuk membatasi akses hanya untuk role tertentu
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		roles, ok := GetRoles(c)
		if !ok || len(roles) == 0 {
			response.Error(c, http.StatusForbidden, "Akses ditolak: tidak ada role yang ditemukan")
			c.Abort()
			return
		}

		for _, role := range roles {
			if allowed[role] {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Akses ditolak: role Anda tidak memiliki izin untuk endpoint ini")
		c.Abort()
	}
}

func RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, ok := GetRoles(c)
		if !ok || len(roles) == 0 {
			response.Error(c, http.StatusForbidden, "Akses ditolak")
			c.Abort()
			return
		}

		for _, role := range roles {
			if StaffRoles[role] {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Akses ditolak: hanya staff yang dapat mengakses endpoint ini")
		c.Abort()
	}
}
