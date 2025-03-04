package main

import (
	v1 "backend/http/api/v1"
	authentication2 "backend/internal/authentication"
	"backend/internal/bot"
	"backend/internal/sync_canvas"
	"backend/pkg/data"
	"backend/pkg/middlewares/authentication"
	"flag"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var addr = flag.String("addr", "localhost:3001", "http service address")
var botToken = "6204994117:AAFYCOHYMxLXRcretwS3_jSkGp-c1MsE6es" // os.Getenv("TG_TOKEN")

var State = sync_canvas.NewSyncCanvas()

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
var dbConnector = data.NewDbConnector("localhost", "3001", "postgres", "example", "takemypixel")
var redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
var botWrapper = bot.NewBotWrapper(botToken, logger)

func main() {
	flag.Parse()

	go botWrapper.Start()

	r := mux.NewRouter()
	corsOnly := r.PathPrefix("/").Subrouter()
	corsOnly.Methods(http.MethodOptions)

	regular := r.PathPrefix("/").Subrouter()
	regular.Methods(http.MethodGet, http.MethodPost)

	application := v1.NewPixelController(State, map[string]interface{}{}, logger)

	corsOnly.HandleFunc("/api/v1/set", application.SetPixel)
	corsOnly.HandleFunc("/api/v1/state", application.GetState)
	corsOnly.HandleFunc("/api/v1/updates", application.WSUpgrader)

	regular.HandleFunc("/api/v1/set", application.SetPixel)
	regular.HandleFunc("/api/v1/state", application.GetState)

	auth := authentication.NewAuthenticationMiddleware(botToken, authentication2.NewTelegramAuthStorage(dbConnector.GetDb(), redisClient), logger)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(c.Handler)
	regular.Use(auth.Chain)

	logger.Info("start on :3001 port")
	log.Fatal(http.ListenAndServe(*addr, r))

	//rl := ratelimit.NewTokenBucket(3, 0, 3*time.Second)
	//rl.StartTokenRefill()
}
