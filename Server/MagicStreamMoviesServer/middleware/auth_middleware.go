package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization token not provided"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()
	}
}
