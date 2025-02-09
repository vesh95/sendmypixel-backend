package main

import (
	"backend/internal/canvas"
	"flag"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var connectionsMap = sync.Map{}

var addr = flag.String("addr", "localhost:3001", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var State = canvas.NewSyncCanvas()

func main() {
	//mux := http.NewServeMux()
	//flag.Parse()
	//log.SetFlags(0)
	//mux.HandleFunc("/updates", connectionHandler)
	//mux.HandleFunc("/set", setPixel)
	//mux.HandleFunc("/state", getState)
	//handler := cors.Default().Handler(mux)
	//
	//log.Fatal(http.ListenAndServe(*addr, handler))
}

//func getState(w http.ResponseWriter, _ *http.Request) {
//	err := json.NewEncoder(w).Encode(State.GetFull())
//	w.Header().Set("Content-Type", "application/json")
//	if err != nil {
//		w.WriteHeader(500)
//		log.Println(err)
//	}
//
//	return
//}
//
//// setPixel
//func setPixel(w http.ResponseWriter, r *http.Request) {
//	pixelData := canvas.PixelDto{}
//	err := json.NewDecoder(r.Body).Decode(&pixelData)
//	if err != nil {
//		slog.Error("Decode error: %s", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	_, err = State.SetPixel(pixelData)
//	slog.Info("Pixel", pixelData)
//	if err != nil {
//		slog.Error("pixel", err)
//		w.WriteHeader(http.StatusBadRequest)
//		_, _ = fmt.Fprint(w, err)
//		return
//	}
//	_, err = broadcastUpdate(pixelData)
//	if err != nil {
//		slog.Error("Not send update pixel")
//	}
//
//	w.WriteHeader(http.StatusOK)
//	slog.Info("Pixel send")
//}
//
//func broadcastUpdate(pixel canvas.PixelDto) (bool, error) {
//	var err error
//	connectionsMap.Range(func(key, value any) bool {
//		conn := key.(*websocket.Conn)
//		_ = conn.WriteJSON(pixel)
//
//		return true
//	})
//
//	return err != nil, err
//}
//
//func connectionHandler(w http.ResponseWriter, r *http.Request) {
//	slog.Info("Connected")
//	c, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		fmt.Println("Upgrade connection error:", err)
//		return
//	}
//	defer c.Close()
//
//	connectionsMap.Store(c, true)
//
//	for {
//		_, _, err := c.ReadMessage()
//		if err != nil {
//			connectionsMap.Delete(c)
//			slog.Info("Disconnect: ", err)
//			break
//		}
//	}
//}
