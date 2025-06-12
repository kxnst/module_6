package main

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"guitar_processor/internal"
	"log"
)

func main() {
	app := fx.New(
		fx.Provide(internal.GetProviders()...),
		fx.Invoke(func() {
			fiberApp := fiber.New()
			fiberApp.Static("/client/auth", "./static/auth.html")
			fiberApp.Static("/client/effects", "./static/effects.html")

			log.Fatal(fiberApp.Listen(":3001"))
		}))

	app.Run()

}
