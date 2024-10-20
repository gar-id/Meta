package api

import (
	"MetaHandler/server/api/process"

	"github.com/gofiber/fiber/v2"
)

func routes(app *fiber.App) {

	api := app.Group("v1")
	apiMetaHandler := api.Group("/centralissh")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})
	apiMetaHandler.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(process.WelcomeGeneral(200, "success", "Welcome to MetaHandler API"))
	})

	// Route for User
	apiMetaHandler.Get("/user/info", func(c *fiber.Ctx) error {
		return process.UserInfo(c)
	})
	apiMetaHandler.Post("/user/update", func(c *fiber.Ctx) error {
		return process.UserUpdate(c)
	})
	apiMetaHandler.Post("/user/delete", func(c *fiber.Ctx) error {
		return process.UserDelete(c)
	})

	// Route for TOTP
	apiMetaHandler.Post("/user/totp", func(c *fiber.Ctx) error {
		return process.TOTPInit(c)
	})

}
