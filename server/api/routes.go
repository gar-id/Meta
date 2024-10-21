package api

import (
	"MetaHandler/server/api/process"

	"github.com/gofiber/fiber/v2"
)

func routes(app *fiber.App) {

	api := app.Group("v1")
	apiMetaHandler := api.Group("/meta")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})
	apiMetaHandler.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})

	// Route for handle meta server
	apiMetaServerHandler := apiMetaHandler.Group("/server")
	apiMetaServerHandler.Get("/info", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServerHandler.Post("/add", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServerHandler.Post("/update", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServerHandler.Post("/delete", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})

	// Route for handle meta service
	apiMetaServiceHandler := apiMetaHandler.Group("/service")
	apiMetaServiceHandler.Get("/info", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServiceHandler.Post("/stop", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServiceHandler.Post("/start", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaServiceHandler.Post("/restart", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})

	// Route for handle meta server stunnel
	apiMetaStunnelHandler := apiMetaHandler.Group("/stunnel")
	apiMetaStunnelHandler.Get("/info", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaStunnelHandler.Post("/add", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaStunnelHandler.Post("/update", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})
	apiMetaStunnelHandler.Post("/delete", func(c *fiber.Ctx) error {
		return process.GetServerInfo(c)
	})

}
