package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"search.com/m/storing"
)

func savetoES(context *gin.Context) {
	var product storing.Product

	// Bind JSON request to product struct
	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save the product to Elasticsearch
	if err := product.Save(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product in Elasticsearch"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Product saved successfully"})
}

func searchES(context *gin.Context) {
	queryText := context.Query("q") // Get query text from URL parameters

	if queryText == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	results, err := storing.SearchingJson(queryText)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results from Elasticsearch"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"results": results})
}
