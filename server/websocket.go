package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mzyy94/nscon"
)

var upgrader = websocket.Upgrader{}
var controller *nscon.Controller

func controllerHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer c.Close()

	defer controller.Close()
	log.Println("Connecting to controller..")
	controller.Connect()
	log.Println("Connected to controller")

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
		json.Unmarshal(message, &controller.Input)
	}
}
