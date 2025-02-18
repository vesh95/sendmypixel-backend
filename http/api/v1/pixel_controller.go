package v1

import (
	"backend/internal/canvas"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type PixelController struct {
	gameState   canvas.Canvas
	configs     map[string]interface{}
	logger      *slog.Logger
	connections sync.Map
	upgrader    websocket.Upgrader
}

func NewPixelController(gameState canvas.Canvas, configs map[string]interface{}, logger *slog.Logger) *PixelController {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return &PixelController{
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
func (c *PixelController) GetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestUuid := uuid.New()

	err := json.NewEncoder(w).Encode(c.gameState.GetFull())
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Error("Bad request response sending error", slog.Group(requestUuid.String(), "path", r.URL.Path, "error", err))
	}
}

// SetPixel обработчик POST /api/v1/set
func (c *PixelController) SetPixel(w http.ResponseWriter, r *http.Request) {
	pixelData := canvas.PixelDto{}
	err := json.NewDecoder(r.Body).Decode(&pixelData)
	if err != nil {
		slog.Error("Decode error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = c.gameState.SetPixel(pixelData)
	slog.Info("Pixel", pixelData)
	if err != nil {
		slog.Error("pixel", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}
	_, err = c.broadcastUpdate(pixelData)
	if err != nil {
		slog.Error("Not send update pixel")
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Pixel send")
}

func (c *PixelController) broadcastUpdate(pixel canvas.PixelDto) (bool, error) {
	var err error
	c.connections.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)
		_ = conn.WriteJSON(pixel)

		return true
	})

	return err != nil, err
}

// WSUpgrader обработчик GET /api/v1/updates
func (c *PixelController) WSUpgrader(w http.ResponseWriter, r *http.Request) {
	slog.Info("Connected")
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade connection error:", err)
		return
	}
	defer conn.Close()

	c.connections.Store(c, true)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			c.connections.Delete(c)
			slog.Info("Disconnect: ", err)
			break
		}
	}
}
