package main

import (
	"backend/canvas"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var connectionsMap = sync.Map{}

var addr = flag.String("addr", "localhost:3001", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var state = canvas.CreateNewState()

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/updates", connectionHandler)
	http.HandleFunc("/set", setPixel)
	http.HandleFunc("/state", getState)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func getState(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(state.GetState())
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
	}

	return
}

func setPixel(w http.ResponseWriter, r *http.Request) {
	pixelData := canvas.PixelDto{}
	err := json.NewDecoder(r.Body).Decode(&pixelData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pixel, err := state.SetPixel(pixelData)
	if err != nil {
		_, _ = fmt.Fprint(w, err)
		return
	}
	broadcastUpdate(pixel)
}

func broadcastUpdate(pixel canvas.PixelDto) {
	connectionsMap.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)
		_ = conn.WriteJSON(pixel)

		return true
	})

	return
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade connection error:", err)
		return
	}
	defer c.Close()

	connectionsMap.Store(c, true)

	for {
		fmt.Println(connectionsMap)
		_, _, err := c.ReadMessage()
		if errors.Is(err, websocket.ErrCloseSent) {
			connectionsMap.Delete(c)
			break
		} else if err != nil {
			fmt.Println("Read message error:", err)
			break
		}
	}
}
