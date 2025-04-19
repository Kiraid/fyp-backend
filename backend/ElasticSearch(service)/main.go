package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"search.com/m/routes"
	rpcserver "search.com/m/rpc_server"
	"search.com/m/storing"
)

// ...
func main() {
	storing.InitES()

	go func() {
		rpcserver.Server()
	}()
	
	router := gin.Default()
	routes.RegisterRoutes(router)
	log.Println("Search Service running on http://localhost:8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("Search Service failed:", err)
	}

}
