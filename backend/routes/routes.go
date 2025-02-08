package routes

import (
	"fyp.com/m/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	userRoutes := server.Group("/users")
	{
		userRoutes.GET("", getUsers)
		userRoutes.GET("/:id", getUser)
		userRoutes.POST("/about", create_about)
		userRoutes.PUT("/:id", updateUser)
		userRoutes.DELETE("/:id", deleteUser)
	}

	// Authentication routes
	server.POST("/signup", signUp) 
	server.POST("/login", login)   
	server.POST("/refresh-token", refreshToken)
	server.POST("/logout", logout)

	server.GET("/search-bar", searchProduct)
	server.GET("/getallcategories", getProductsCategories)
	server.GET("/getallproducts", getProducts)
	server.GET("Product/:id", getProductbyID)

	authenticated := server.Group("/products")
	authenticated.Use(middleware.Authenticate)
	authenticated.POST("/new-product", createProduct)
	authenticated.PUT("/update-product/:id", updateProduct)
	authenticated.DELETE("/delete-product/:id", deleteProduct)
	authenticated.GET("/get-Product/:id", getProductsbyuserid)
	

	categoryRoutes := server.Group("/category")
	//categoryRoutes.Use(middleware.Authenticate)
	{
		categoryRoutes.POST("/new-category", createProductCategory)
		categoryRoutes.PUT("/update-category/:id", updateProductCategory)
		categoryRoutes.DELETE("delete-category/:id", deleteProductCategory)
		categoryRoutes.GET("/get-Category/:id", getProductCategorybyuserid)
	}

	orderRoutes := server.Group("/order")
	{
		orderRoutes.GET("orders", getOrders)
		orderRoutes.POST("new-order", saveOrder)
		orderRoutes.GET("buyer-orders/:id", getOrdersByBuyerid)
		orderRoutes.GET("seller-orders/:id", getOrdersBySellerid)
		orderRoutes.GET("order-detail/:id", getOrderByid)
		orderRoutes.PUT("update-status/:id", UpdateOrderStatus)
		orderRoutes.DELETE("delete-order/:id", deleteOrder)
		orderRoutes.GET("seller-stats/:id", getsellerStats)
		orderRoutes.GET("monthly-stats/:id", getsellermonthlyStats)
	}

	ReviewRoutes := server.Group("/review")
	{
		ReviewRoutes.GET("reviews", getReviews)
		ReviewRoutes.GET("product-review/:id", getReviewsByProductId)
		ReviewRoutes.POST("new-review", saveReview)
	}
}
