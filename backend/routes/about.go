package routes

import (
	"net/http"

	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func create_about(context *gin.Context){
    var about models.Seller_about
    err := context.ShouldBindJSON(&about)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"message": "Couldnt parse request data for seller desc"})
        return
    }
    err =  about.Save()
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt save seller desc", "error":   err.Error()})
        return
    }
    context.JSON(http.StatusCreated, gin.H{"message": "Seller description added", "event": about})
}


