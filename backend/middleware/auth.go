package middleware

import (
	"net/http"

	"fyp.com/m/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	authheader := context.Request.Header.Get("Authorization")
	if authheader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	const bearerPrefix = "Bearer "
    if len(authheader) < len(bearerPrefix) || authheader[:len(bearerPrefix)] != bearerPrefix {
        context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
        return
    }
	token := authheader[len(bearerPrefix):]
	id, role, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Wrong token"})
		return
	}


	context.Set("userId", id)
	context.Set("role", role)
	context.Next()

}