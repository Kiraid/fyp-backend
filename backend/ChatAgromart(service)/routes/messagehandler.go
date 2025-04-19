package routes

import (
	"net/http"
	"strconv"

	"example.com/chat/models"
	"github.com/gin-gonic/gin"
)

func NewMessage(context *gin.Context) {
	var message models.Message
	err := context.ShouldBindJSON(&message)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Couldnt parse request data"})
		return
	}
	userId := context.GetInt64("userId")
	message.SenderID = userId
	err = message.CreateMessage()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt save message", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Message saved", "event": message})
}

func GetMessages(context *gin.Context) {
	receiverID := context.Query("receiverId")     // Get receiverId from query
	limit := context.DefaultQuery("limit", "20")  // Default to 20 if not provided
	offset := context.DefaultQuery("offset", "0") // Default to 0 if not provided
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	receiverIDint, err := strconv.Atoi(receiverID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reciever Id"})
		return
	}

	userId := context.GetInt64("userId")
	messages, err := models.GetAllMessages(int(userId), receiverIDint, limitInt, offsetInt)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get messages", "details": err.Error()})
		return
	}

	var formattedMessages []map[string]string
	for _, message := range messages {
		username, err := models.GetUserNamebyId(message.SenderID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't get sender username"})
			return
		}

		formattedMessages = append(formattedMessages, map[string]string{
			"user":    username,
			"content": message.Content,
		})
	}
	context.JSON(http.StatusOK, gin.H{"messages": formattedMessages})

}
