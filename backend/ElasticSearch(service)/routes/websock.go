package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"search.com/m/storing"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (Modify if needed for security)
	},
}

func searchESWS(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println("Failed to set WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from WebSocket (user query)
		_, queryText, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Println("Received query:", string(queryText))

		// Call Searching function to get Elasticsearch results
		results, err := storing.SearchingJson(string(queryText))
		if err != nil {
			errorMessage := "Failed to fetch results from Elasticsearch"
			log.Println(errorMessage)
			conn.WriteJSON(gin.H{"error": errorMessage})
			continue
		}

		// Send results back to WebSocket client
		err = conn.WriteJSON(gin.H{"results": results})
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
