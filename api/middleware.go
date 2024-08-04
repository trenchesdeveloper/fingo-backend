package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")
		if len(tokenSplit) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		log.Println("got token33 ", token)
		log.Println("got token34 ", tokenSplit[1])
		log.Println("got token35 ", tokenSplit[0])
		// Check if token is valid with bearer
		if tokenSplit[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Check if token is valid
		userID, err := tokenController.VerifyToken(tokenSplit[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
