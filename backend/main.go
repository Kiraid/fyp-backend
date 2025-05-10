package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fyp.com/m/db"
	"fyp.com/m/kafka_module"
	"fyp.com/m/models"
	"fyp.com/m/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
//main function
func main() {
	db.InitDB()
	db.CreateTable()
	models.InitRedis()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka_module.ConsumerwithShutdown(ctx)
	}()

	server := gin.Default()
	server.Static("/static", "./static")
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	routes.RegisterRoutes(server)

	// Create an HTTP server from Gin
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	wg.Add(1)
	// Start the HTTP server in a goroutine
	go func() {
		defer wg.Done()
		fmt.Println("Server is running on port 8080...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan // Wait for a shutdown signal

	// Cleanup
	fmt.Println("\nShutdown signal received. Closing application...")
	cancel() // Notify goroutines to stop

	// Gracefully shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Error shutting down server:", err)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Application exited cleanly.")
}
