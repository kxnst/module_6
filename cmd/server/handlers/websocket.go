package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"guitar_processor/cmd/server/request"
	"guitar_processor/cmd/server/utils"
	_ "guitar_processor/docs"
	"guitar_processor/internal/effect/service"
	"guitar_processor/internal/entity"
)

type WebSocketHandler struct {
	sessions map[string]request.DeviceInfoRequest
	as       service.AudioService
	aus      utils.AuthService
}

// HandleInit godoc
// @Summary Initialize session
// @Description Receives I/O configuration
// @Tags Websocket
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body request.DeviceInfoRequest true "Інформація про пристрої"
// @Success 200 {string} string "ok"
// @Failure 400 {object} ErrorResponse "error"
// @Router /ws [post]
func (w *WebSocketHandler) HandleInit(c *fiber.Ctx) error {
	var req request.DeviceInfoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	user := c.Locals("user")

	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	w.sessions[user.(*entity.User).ID] = req

	return c.SendStatus(200)
}

func (w *WebSocketHandler) HandleStream(c *websocket.Conn) error {
	defer c.Close()

	token := c.Query("token")
	if token == "" {
		return w.GetJsonError(c, "Missing token")
	}

	user, err := w.aus.ValidateTokenAndGetUser(token)
	if err != nil {
		return w.GetJsonError(c, "Unauthorized")
	}

	deviceInfo, ok := w.sessions[user.ID]
	if !ok {
		return w.GetJsonError(c, "Device not initialized")
	}

	str, err := w.as.Start(&deviceInfo)

	defer func() {
		w.as.Terminate()
		c.Close()
		str.Stop()
		str.Close()
	}()

	if err != nil {
		return w.GetJsonError(c, "Error starting stream")
	}

	err = str.Start()

	if err != nil {
		return w.GetJsonError(c, "Error starting stream: "+err.Error())
	}

	for {
		var msg struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := c.ReadJSON(&msg); err != nil {
			c.WriteJSON(&WebSocketMessage{
				Type: "error",
				Data: err.Error(),
			})
			break
		}

		switch msg.Type {
		case "effects_list":
			effects := w.as.GetEffectsInfo()
			if err := c.WriteJSON(&WebSocketMessage{
				Type: "effect_list_response",
				Data: effects,
			}); err != nil {
				log.Warn("Failed to send effect list:", err)
				break
			}
		case "add_effect":
			var payload struct {
				Slug string `json:"slug"`
			}
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				c.WriteJSON(&WebSocketMessage{
					Type: "error",
					Data: err.Error(),
				})
				break
			}

			effect, err := w.as.AddEffectToChain(payload.Slug)
			if err != nil {
				c.WriteJSON(&WebSocketMessage{
					Type: "error",
					Data: err.Error(),
				})
			} else {
				c.WriteJSON(&EffectAddedMessage{
					Type:       "effect_added",
					Slug:       effect.GetSlug(),
					Parameters: effect.GetParameters(),
				})
			}

		case "set_param":
			var payload struct {
				UUID  string  `json:"uuid"`
				Value float32 `json:"value"`
			}
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				c.WriteJSON(&WebSocketMessage{
					Type: "error",
					Data: err.Error(),
				})
				break
			}

			err := w.as.SetEffectParameter(payload.UUID, payload.Value)
			if err != nil {
				c.WriteJSON(&WebSocketMessage{
					Type: "error",
					Data: err.Error(),
				})
				break
			}

			c.WriteJSON(&WebSocketMessage{
				Type: "param_updated",
				Data: payload,
			})
		}

	}

	return nil
}

func (w *WebSocketHandler) GetJsonError(c *websocket.Conn, msg string) error {
	return c.WriteJSON(WebSocketMessage{
		Type: "error",
		Data: msg,
	})

}

func NewWebSocketHandler(as service.AudioService, aus utils.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		sessions: make(map[string]request.DeviceInfoRequest),
		as:       as,
		aus:      aus,
	}
}
