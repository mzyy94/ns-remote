package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func controllerHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			switch e := err.(type) {
			case *websocket.CloseError:
				log.Println("WebSocket closed:", e.Code)
			default:
				log.Println("WebSocket read error:", e)
			}
			break
		}
		log.Printf("WebSocket recv data: %s", message)
	}
}
