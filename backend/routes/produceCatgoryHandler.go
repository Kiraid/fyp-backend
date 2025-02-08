package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func createProductCategory(context *gin.Context) {
	name := context.PostForm("name")
	description := context.PostForm("description")
	userId, err := strconv.ParseInt(context.PostForm("userId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid userId", "error:": err.Error()})
		return
	}
	file, err := context.FormFile("image")
	var filePath string
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Image upload failed", "error": err})
		return
	}
	randomFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath = fmt.Sprintf("static/images/%s", randomFileName)
	if err := context.SaveUploadedFile(file, filePath); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
		return
	}

	productcategory := models.ProductCategory{
		Name:        name,
		Description: description,
		ImagePath:   filePath,
		UserID:      userId,
	}

	err = productcategory.SaveCategory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt save product category", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Product category created successfully", "product": productcategory})
}

func getProductsCategories(context *gin.Context) {
	products, err := models.GetAllProductsCategory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch products categories", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, products)
}

func updateProductCategory(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product category id"})
	}
	var updatedproductcategory models.ProductCategory
	updatedproductcategory.Name = context.PostForm("name")
	updatedproductcategory.Description = context.PostForm("description")

	file, err := context.FormFile("image")
	if err == nil {
		randomFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("static/images/%s", randomFileName)
		if err := context.SaveUploadedFile(file, filePath); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
			return
		}
		updatedproductcategory.ImagePath = filePath
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt handle image"})
		return
	}
	updatedproductcategory.ID = id
	err = updatedproductcategory.UpdateProductCategory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Couldn't update product category, error: %s", err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product Category updated"})
}

func deleteProductCategory(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product category id"})
	}
	// userId := context.GetInt64("userId")
	productcategory, err := models.GetProductCategoryByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product category"})
	}
	// if productcategory.UserID != userId {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete product category"})
	// 	return
	// }
	err = productcategory.DeleteProductCategory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product category", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product category deleted"})

}

func getProductCategorybyuserid(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse product id for read"})
		return
	}
	products, err := models.GetProductsCategorybyUserID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product categories", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, products)
}
