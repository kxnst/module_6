package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/fx"
	"guitar_processor/cmd/server/handlers"
	"guitar_processor/cmd/server/middlewares"
	"guitar_processor/cmd/server/utils"
	_ "guitar_processor/docs"
	"guitar_processor/internal"
)

func main() {
	provides := internal.GetProviders()
	provides = append(provides, handlers.GetProviders()...)
	provides = append(provides, utils.NewAuthService)

	app := fx.New(
		fx.Provide(provides...),
		fx.Invoke(func(
			ah *handlers.AuthHandler,
			dh *handlers.DevicesHandler,
			wsh *handlers.WebSocketHandler,
			am *middlewares.AuthMiddleware,
			eh *handlers.EffectsHandler,
		) {
			fiberApp := fiber.New()

			authClosure := func(c *fiber.Ctx) error { return am.RequireAuth(c) }
			fiberApp.Use(middlewares.AllowCors())
			fiberApp.Post("/auth", func(c *fiber.Ctx) error { return ah.Handle(c) })
			fiberApp.Get("/devices", authClosure, func(c *fiber.Ctx) error { return dh.Handle(c) })
			fiberApp.Get("/effects", authClosure, func(c *fiber.Ctx) error { return eh.Handle(c) })
			fiberApp.Post("/ws", authClosure, func(c *fiber.Ctx) error {
				return wsh.HandleInit(c)
			})
			fiberApp.Get("/ws", websocket.New(func(conn *websocket.Conn) {
				if err := wsh.HandleStream(conn); err != nil {
					log.Error("WebSocket error: ", err)
				}
			}))
			fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

			log.Fatal(fiberApp.Listen(":3000"))
		}),
	)

	app.Run()
}
