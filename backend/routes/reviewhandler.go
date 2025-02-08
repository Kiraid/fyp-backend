package routes

import (
	"net/http"
	"strconv"

	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func getReviews(context *gin.Context) {
	reviews, err := models.GetAllReviews()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reviews", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, reviews)
}

func saveReview(context *gin.Context) {
	var review models.Review
	err := context.ShouldBindJSON(&review)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	err = review.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save review", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Review created", "review": review})
}

func getReviewsByProductId(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't parse product ID for read"})
		return
	}
	reviews, err := models.GetReviewsByProductID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't fetch reviews", "error": err.Error()})
		return
	}

	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}
	avgRating := 0
	if len(reviews) > 0 {
		avgRating = totalRating / len(reviews)
	}

	stats := struct {
		AvgRating    int `json:"avg_rating"`
		TotalReviews int `json:"total_reviews"`
	}{
		AvgRating:    avgRating,
		TotalReviews: len(reviews),
	}

	context.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"stats":   stats,
	})
}


