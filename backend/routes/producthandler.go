package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func getProducts(context *gin.Context) {
	products, err := models.GetAllProducts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch products", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, products)
}

func createProduct(context *gin.Context) {
	name := context.PostForm("name")
	description := context.PostForm("description")
	categoryName := context.PostForm("categoryName")
	userId, err := strconv.ParseInt(context.PostForm("userId"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	price, err := strconv.ParseInt(context.PostForm("price"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid price"})
		return
	}
	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Image upload failed"})
		return
	}
	randomFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := fmt.Sprintf("static/images/%s", randomFileName)
	if err := context.SaveUploadedFile(file, filePath); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
		return
	}

	product := models.Product{
		Name:        name,
		Description: description,
		ImagePath:   filePath,
		UserID:      userId,
		Category_name: categoryName,
		Price: price,
	}

	err = product.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt save product", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

func updateProduct(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product id"})
	}
	role := context.GetString("role")
	if role != "admin" {
	userId := context.GetInt64("userId")
	product, err := models.GetProductByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product"})
	}
	if product.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update product"})
		return
	}
	}
	var updatedproduct models.Product
	updatedproduct.Name = context.PostForm("name")
	updatedproduct.Description = context.PostForm("description")
	updatedproduct.Category_name = context.PostForm("categoryName")
	updatedproduct.Price, err = strconv.ParseInt(context.PostForm("price"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid price"})
		return
	}

	file, err := context.FormFile("image")
	if err == nil {
		randomFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("static/images/%s", randomFileName)
		if err := context.SaveUploadedFile(file, filePath); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
			return
		}
		updatedproduct.ImagePath = filePath
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt handle image"})
		return
	}
	updatedproduct.ID = id
	err = updatedproduct.UpdateProduct()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Couldn't update product, error: %s", err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func deleteProduct(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse product id for delete"})
		return
	}

	product, err := models.GetProductByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product for delete"})
	}
	role := context.GetString("role")
	if role != "admin" {
	userID := context.GetInt64("userId")
	if product.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete product"})
		return
	}
}
	err = product.DeleteProduct()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product deleted"})

}

func getProductsbyuserid(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse product id for read"})
		return
	}
	products, err := models.GetProductsbyUserID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch products", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, products)
}

func getProductbyID(context *gin.Context){
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse product id for read"})
		return
	}
	product, err := models.GetProductByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt get product"})
		return
	}
	user, err := models.GetUserbyID(product.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch product", "error": err.Error()})
		return
	}
	type Response struct {
		Product  *models.Product `json:"product"`
		Username string          `json:"username"`
	}
	response := Response{
		Product:  product,
		Username: user.Username,
	}
	context.JSON(http.StatusOK, response)
}

func searchProduct(context *gin.Context) {
	var request struct {
		Search string `json:"search" binding:"required"`
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	products, err := models.GetProductsbySearch(request.Search)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch products", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, products)
}
