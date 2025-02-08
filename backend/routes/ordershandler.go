package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyp.com/m/models"
	"github.com/gin-gonic/gin"
)

func getOrders(context *gin.Context) {
	orders, err := models.GetAllOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch orders", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, orders)
}

func saveOrder(context *gin.Context) {
	var order models.Order
	err := context.ShouldBindJSON(&order)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	err = order.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save order", "error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Order created", "order": order})
}

func getOrdersByBuyerid(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse order id for read"})
		return
	}
	orders, err := models.GetOrderByBuyerID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch orders", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, orders)
}

func getOrdersBySellerid(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse order id for read"})
		return
	}
	orders, err := models.GetOrderBySellerID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch orders", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, orders)
}

func getOrderByid(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse order id for read"})
		return
	}
	order, err := models.GetOrderByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch orders", "error": err.Error()})
		return
	}
	product, err := models.GetProductByID(order.ProductID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch orders", "error": err.Error()})
		return
	}
	type Response struct {
		Order   *models.Order
		Product *models.Product
	}
	response := Response{
		Order:   order,
		Product: product,
	}

	context.JSON(http.StatusOK, response)
}

func UpdateOrderStatus(ctx *gin.Context) {
	orderIDParam := ctx.Param("id")
	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID", "error": err.Error()})
		return
	}
	var requestBody struct {
		OrderStatus string `json:"orderStatus"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}
	order, err := models.GetOrderByID(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch order", "error": err.Error()})
		return
	}
	order.OrderStatus = strings.ToLower(requestBody.OrderStatus)

	if err := order.UpdateOrderStatus(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "order": order})
}

func deleteOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt parse Order id"})
	}
	// userId := context.GetInt64("userId")
	order, err := models.GetOrderByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch order"})
	}
	// if productcategory.UserID != userId {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete product category"})
	// 	return
	// }
	err = order.DeleteOrder()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt delete order", "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Order deleted"})

}

func getsellerStats(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt convert parse order id for read"})
		return
	}
	orders, err := models.GetOrderBySellerID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch orders", "error": err.Error()})
		return
	}
	type sellerStats struct {
		TotalOrders  int
		Revenue      int
		ActiveOrders int
		TotalSales   int
	}
	var stats sellerStats
	for _, order := range orders {
		stats.TotalOrders += 1
		stats.Revenue += int(order.CheckoutPrice)
		if order.OrderStatus == "completed" {
			stats.TotalSales += int(order.CheckoutPrice)
		} else {
			stats.ActiveOrders += 1
		}
	}
	context.JSON(http.StatusOK, stats)
}

func getsellermonthlyStats(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't parse seller ID"})
		return
	}

	orders, err := models.GetOrderBySellerID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't fetch orders", "error": err.Error()})
		return
	}

	type MonthlyStats struct {
		Year            int     `json:"year"`
		Month           int     `json:"month"`
		TotalOrders     int     `json:"total_orders"`
		CompletedOrders int     `json:"completed_orders"`
		TotalRevenue    float64 `json:"total_revenue"`
		RevenueGrowth   float64 `json:"revenue_growth_percentage"`
		OrderGrowth     float64 `json:"order_growth_percentage"`
	}

	type YearlyStats struct {
		Year            int     `json:"year"`
		TotalOrders     int     `json:"total_orders"`
		CompletedOrders int     `json:"completed_orders"`
		TotalRevenue    float64 `json:"total_revenue"`
		RevenueGrowth   float64 `json:"revenue_growth_percentage"`
		OrderGrowth     float64 `json:"order_growth_percentage"`
	}

	monthlyStats := make(map[string]*MonthlyStats)
	yearlyStats := make(map[int]*YearlyStats)

	currentYear, currentMonth, _ := time.Now().Date()
	var currentMonthRevenue float64
	var currentYearRevenue float64

	for _, order := range orders {
		orderTime := time.Unix(order.Time, 0)
		year, month := orderTime.Year(), int(orderTime.Month())

		monthKey := strconv.Itoa(year) + "-" + strconv.Itoa(month)

		if _, exists := monthlyStats[monthKey]; !exists {
			monthlyStats[monthKey] = &MonthlyStats{Year: year, Month: month, TotalOrders: 0, CompletedOrders: 0, TotalRevenue: 0}
		}
		monthlyStats[monthKey].TotalOrders++

		if order.OrderStatus == "completed" {
			monthlyStats[monthKey].CompletedOrders++
			monthlyStats[monthKey].TotalRevenue += float64(order.CheckoutPrice)

			if year == currentYear && month == int(currentMonth) {
				currentMonthRevenue += float64(order.CheckoutPrice)
			}
		}

		if _, exists := yearlyStats[year]; !exists {
			yearlyStats[year] = &YearlyStats{Year: year, TotalOrders: 0, CompletedOrders: 0, TotalRevenue: 0}
		}
		yearlyStats[year].TotalOrders++

		if order.OrderStatus == "completed" {
			yearlyStats[year].CompletedOrders++
			yearlyStats[year].TotalRevenue += float64(order.CheckoutPrice)

			if year == currentYear {
				currentYearRevenue += float64(order.CheckoutPrice)
			}
		}
	}

	var monthlyResults []MonthlyStats
	var yearlyResults []YearlyStats
	var previousMonthRevenue, previousMonthOrders float64
	var previousYearRevenue, previousYearOrders float64

	for _, stats := range monthlyStats {
		if previousMonthRevenue > 0 {
			stats.RevenueGrowth = ((stats.TotalRevenue - previousMonthRevenue) / previousMonthRevenue) * 100
		} else if previousMonthRevenue == 0 && stats.TotalRevenue > 0 {
			stats.RevenueGrowth = 100
		}

		if previousMonthOrders > 0 {
			stats.OrderGrowth = ((float64(stats.CompletedOrders) - previousMonthOrders) / previousMonthOrders) * 100
		} else if previousMonthOrders == 0 && stats.CompletedOrders > 0 {
			stats.OrderGrowth = 100
		}

		previousMonthRevenue = stats.TotalRevenue
		previousMonthOrders = float64(stats.CompletedOrders)

		monthlyResults = append(monthlyResults, *stats)
	}

	for _, stats := range yearlyStats {
		if previousYearRevenue > 0 {
			stats.RevenueGrowth = ((stats.TotalRevenue - previousYearRevenue) / previousYearRevenue) * 100
		} else if previousYearRevenue == 0 && stats.TotalRevenue > 0 {
			stats.RevenueGrowth = 100
		}

		if previousYearOrders > 0 {
			stats.OrderGrowth = ((float64(stats.CompletedOrders) - previousYearOrders) / previousYearOrders) * 100
		} else if previousYearOrders == 0 && stats.CompletedOrders > 0 {
			stats.OrderGrowth = 100
		}

		previousYearRevenue = stats.TotalRevenue
		previousYearOrders = float64(stats.CompletedOrders)

		yearlyResults = append(yearlyResults, *stats)
	}

	context.JSON(http.StatusOK, gin.H{
		"monthly_stats":         monthlyResults,
		"yearly_stats":          yearlyResults,
		"current_month_revenue": currentMonthRevenue,
		"current_year_revenue":  currentYearRevenue,
	})
}
