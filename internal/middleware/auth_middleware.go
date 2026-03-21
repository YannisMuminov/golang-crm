package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/YannisMuminov/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserID      = "user_id"
	ContextUserEmail   = "user_email"
	ContextUserRole    = "user_role"
	ContextPermissions = "permissions"
)

type jwtClaims struct {
	UserID      int64    `json:"user_id"`
	Role        string   `json:"role"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func RequireAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, expected: Bearer <token>",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if strings.ToLower(parts[0]) != "bearer" || len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, expected: Bearer <token>",
			})
			return
		}

		tokenStr := parts[1]

		claims := &jwtClaims{}

		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(cfg.Secret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserEmail, claims.Email)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextPermissions, claims.Permissions)

		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextPermissions)

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no permission found",
			})
			return
		}

		permList, ok := val.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid permission format",
			})
			return
		}

		for _, p := range permList {
			if p == permission {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})

	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return 0, false
	}

	id, ok := val.(int64)

	return id, ok
}

func GetUserRole(c *gin.Context) string {
	val, _ := c.Get(ContextUserRole)
	role, _ := val.(string)
	return role
}
