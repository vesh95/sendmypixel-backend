package v1

import (
	"backend/pkg/authentication"
	"backend/pkg/canvas"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type Application struct {
	gameState   canvas.Canvas
	configs     map[string]interface{}
	logger      *slog.Logger
	connections sync.Map
	upgrader    websocket.Upgrader
}

type connectionInfo struct {
	userId int64
}

func NewPixelController(gameState canvas.Canvas, configs map[string]interface{}, logger *slog.Logger) *Application {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return &Application{
		gameState:   gameState,
		configs:     configs,
		logger:      logger,
		connections: sync.Map{},
		upgrader:    upgrader,
	}
}

type badRequestResponse struct {
	Message   error  `json:"message"`
	RequestId string `json:"request_id"`
}

// GetState обработчик GET /api/v1/state
func (c *Application) GetState(w http.ResponseWriter, r *http.Request) {
	//u := c.extractUser(r)
	w.Header().Set("Content-Type", "application/json")
	requestUuid := uuid.New()
	slog.Debug("Getting state", slog.Group(requestUuid.String()))

	err := json.NewEncoder(w).Encode(c.gameState.GetFull())
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Warn("Bad request response sending error", slog.Group(requestUuid.String(), "path", r.URL.Path, "error", err))
	}
}

// SetPixel обработчик POST /api/v1/set
func (c *Application) SetPixel(w http.ResponseWriter, r *http.Request) {
	u := c.extractUser(r)
	pixelData := canvas.PixelDto{}
	pixelData.UserId = u.Id

	err := json.NewDecoder(r.Body).Decode(&pixelData)
	if err != nil {
		c.logger.Error("Decode error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = c.gameState.SetPixel(pixelData)
	c.logger.Info("Pixel", pixelData)
	if err != nil {
		c.logger.Error("pixel", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}
	go c.broadcastUpdate(pixelData)
	if err != nil {
		c.logger.Error("Not send update pixel")
	}

	w.WriteHeader(http.StatusOK)
	c.logger.Info("Pixel send")
}

func (c *Application) broadcastUpdate(pixel canvas.PixelDto) {
	c.connections.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)
		err := conn.WriteJSON(pixel)
		if err != nil {
			ci := value.(connectionInfo)
			c.logger.Warn("can't send message", "user_id", ci.userId)
		}

		return true
	})
}

// WSUpgrader обработчик GET /api/v1/updates
func (c *Application) WSUpgrader(w http.ResponseWriter, r *http.Request) {
	requestUuid := uuid.New()

	c.logger.Info("Connected")
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error("Upgrade connection error", slog.Group(requestUuid.String(), "error", err))
		return
	}
	defer conn.Close()

	c.connections.Store(conn, &connectionInfo{userId: 1})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			c.connections.Delete(c)
			fmt.Println("Disconnect: ", err)
			break
		}
	}
}

func (c *Application) extractUser(r *http.Request) authentication.User {
	val, ok := r.Context().Value(authentication.UserContextKey).(authentication.User)
	if ok {
		return val
	}

	c.logger.Error("can't get user data from request context", "route", r.URL, "context", r.Context())
	panic("can't get user data from request context")
}
