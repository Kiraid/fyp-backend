package main

import (
	"example.com/chat/db"
	"example.com/chat/models"
	"example.com/chat/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// WebSocket Client Struct
type Client struct {
	conn       *websocket.Conn
	senderID   int64
	receiverID int64
	send       chan []byte
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *models.Message)
var mutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket Connection Handler
func wsEndpoint(c *gin.Context) {
	senderID, err := strconv.ParseInt(c.Query("senderID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	receiverID, err := strconv.ParseInt(c.Query("receiverID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not open WebSocket connection"})
		return
	}

	client := &Client{conn: ws, senderID: senderID, receiverID: receiverID, send: make(chan []byte)}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	go handleMessages(client)
	go client.writeMessages()
}

// Handle Incoming Messages
func handleMessages(client *Client) {
	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		// Save message to database
		message := models.Message{
			SenderID:   client.senderID,
			RecieverID: client.receiverID,
			Content:    string(msg),
			Time:       time.Now().Unix(),
		}
		message.CreateMessage()

		// Send message to the intended recipient
		broadcast <- &message
	}
}

// Send Messages to WebSocket Clients
func (c *Client) writeMessages() {
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// Broadcast Messages to the Intended Receiver
func handleBroadcast() {
	for {
		message := <-broadcast
		msgBytes := []byte(message.Content)

		mutex.Lock()
		for client := range clients {
			if (client.receiverID == message.RecieverID && client.senderID == message.SenderID) ||
				(client.receiverID == message.SenderID && client.senderID == message.RecieverID) {
				select {
				case client.send <- msgBytes:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
}

func main() {
	// Initialize the database
	db.InitDB()
	db.CreateTable()

	// Create a new Gin router
	router := gin.Default()

	// Register WebSocket endpoint
	router.GET("/ws", wsEndpoint)

	// Register other API routes
	routes.RegisterRoutes(router)

	// Start WebSocket broadcasting in a separate goroutine
	go handleBroadcast()

	// Start the chat service on port 8081
	log.Println("Chat Service running on http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Chat Service failed:", err)
	}
}
