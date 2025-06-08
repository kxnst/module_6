package handlers

import (
	"github.com/gofiber/fiber/v2"
	_ "guitar_processor/docs"
	_ "guitar_processor/internal/effect/dsp"
	"guitar_processor/internal/effect/service"
)

type DevicesHandler struct {
	as service.AudioService
}

// Handle godoc
// @Summary Get devices list
// @Description Returns list of available devices
// @Tags Devices
// @Accept json
// @Produce json
// @Success 200 {array} []dsp.Device "devices list"
// @Failure 400 {object} ErrorResponse "error"
// @Router /devices [get]
func (dh *DevicesHandler) Handle(c *fiber.Ctx) error {
	devices, err := dh.as.GetDevices()

	if err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(devices)
}

func NewDevicesHandler(as service.AudioService) *DevicesHandler {
	return &DevicesHandler{as: as}
}
