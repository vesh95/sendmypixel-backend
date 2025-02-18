package main

import (
	v1 "backend/http/api/v1"
	"backend/internal/api"
	"backend/internal/canvas"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"github.com/rs/cors"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var connectionsMap = sync.Map{}

var addr = flag.String("addr", "localhost:3001", "http service address")

var State = canvas.NewSyncCanvas()

func main() {
	flag.Parse()
	log.SetFlags(0)

	v1PixelController := v1.NewPixelController(State, map[string]interface{}{}, &slog.Logger{})
	r := pat.New()
	r.Get("/state", v1PixelController.GetState)
	r.Post("/set", v1PixelController.SetPixel)
	r.Get("/updates", v1PixelController.WSUpgrader)

	m := mux.NewRouter()
	m.Handle("/api/v1", cors.Default().Handler(r))
	tbHandler := api.NewTokenBucket(1000, 1000, 100*time.Millisecond).RateLimitMiddleware(m)

	log.Fatal(http.ListenAndServe(*addr, tbHandler))
}
