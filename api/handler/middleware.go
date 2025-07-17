package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/Akrom0181/Auth/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const SUPER_USER_ID = "14be02d4-71fb-49e1-847a-5e0fabd99bd0"

func (h *Handler) AuthorizerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseJWT(token, h.Config.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token 1"})
			return
		}

		userID, ok1 := claims["user_id"].(string)
		userType, ok2 := claims["user_type"].(string)
		expFloat, ok3 := claims["exp"].(float64)
		if !ok1 || !ok2 || !ok3 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		if time.Now().Unix() > int64(expFloat) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		endpoint := c.FullPath()

		restricted := map[string]bool{
			"/roles/create":    true,
			"/roles/update":    true,
			"/roles/list":      true,
			"/sysusers/create": true,
		}

		if restricted[endpoint] {
			if userType != "sysuser" || userID != SUPER_USER_ID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
				return
			}
		}

		c.Set("user_id", userID)
		c.Set("user_type", userType)

		c.Next()
	}
}
