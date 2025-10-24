package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	// File server for static HTML
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", webSocketHandler)

	log.Println("🚀 Server starting on http://localhost:8080")
	log.Println("📁 Serving static files from ./static/")
	log.Println("🔌 WebSocket endpoint: ws://localhost:8080/ws")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Get client address for better logging
	clientAddr := conn.RemoteAddr().String()
	log.Printf("✅ Client connected: %s", clientAddr)

	// Echo loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check if it's a normal close
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("👋 Client disconnected normally: %s", clientAddr)
			} else {
				log.Printf("❌ Read error from %s: %v", clientAddr, err) // ← Fixed: %v not %e
			}
			return
		}

		log.Printf("📨 Received from %s: %s", clientAddr, string(message))

		// Echo message back
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("❌ Write error to %s: %v", clientAddr, err)
			return
		}

		log.Printf("📤 Echoed to %s: %s", clientAddr, string(message))
	}
}
