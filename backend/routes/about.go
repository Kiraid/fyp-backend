package routes

import (
	"net/http"

	"fyp.com/m/middleware"
	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func create_about(context *gin.Context) {
	var about models.Seller_about
	err := context.ShouldBindJSON(&about)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Couldnt parse request data for seller desc"})
		return
	}
	err = about.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt save seller desc", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Seller description added", "event": about})
}

func authCallbackHandler(c *gin.Context) {
	authCode := c.Query("code")
	if authCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing auth code"})
		return
	}

	// Step 1: Exchange Auth Code for Access Token
	tokenData, err := middleware.ExchangeAuthCode(authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange auth code", "error": err.Error()})
		return
	}

	// Step 2: Fetch User Info
	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid token response"})
		return
	}

	userInfo, err := middleware.GetUserInfo(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user info", "error": err.Error()})
		return
	}

	// Respond with user info
	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}
