package handlers

import (
	"github.com/gofiber/fiber/v2"
	"guitar_processor/internal/effect/service"
)

type EffectsHandler struct {
	as service.AudioService
}

func NewEffectsHandler(as service.AudioService) *EffectsHandler {
	return &EffectsHandler{as: as}
}

// Handle godoc
// @Summary Отримати список доступних ефектів
// @Tags Effects
// @Produce json
// @Success 200 {array} dsp.EffectInfo
// @Failure 401 {object} ErrorResponse
// @Router /effects [get]
func (h *EffectsHandler) Handle(c *fiber.Ctx) error {
	effects := h.as.GetEffectsInfo()

	return c.JSON(effects)
}
