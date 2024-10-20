package api

import (
	"errors"
	"fmt"

	"MetaHandler/server/config/caches"
	"MetaHandler/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start() {
	fibercfg := fiber.Config{
		Prefork:               false,
		ServerHeader:          "MetaHandler by Gar",
		Concurrency:           256 * 1024 * 30,
		DisableStartupMessage: true,

		// // Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	}
	app := fiber.New(fibercfg)

	// Setup cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// // Setup basicAuth
	// app.Use(basicauth.New(basicauth.Config{
	// 	Users: map[string]string{
	// 		caches.MetaHandlerConfig.MetaHandler.API.UserAuth: caches.MetaHandlerConfig.MetaHandler.API.PasswordAuth,
	// 	},
	// 	Realm: "Forbidden",
	// 	Authorizer: func(user, pass string) bool {
	// 		if user == caches.MetaHandlerConfig.MetaHandler.API.UserAuth && pass == caches.MetaHandlerConfig.MetaHandler.API.PasswordAuth {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// 	Unauthorized: func(c *fiber.Ctx) error {
	// 		return c.JSON("Forbidden")
	// 	},
	// 	ContextUsername: "_user",
	// 	ContextPassword: "_pass",
	// }))

	// routes config
	routes(app)

	err := app.Listen(fmt.Sprintf("%v:%v", caches.MetaHandlerServer.MetaHandlerServer.API.HTTPHost, caches.MetaHandlerServer.MetaHandlerServer.API.HTTPPort))
	if err != nil {
		tools.ZapLogger("both", "server").Fatal(fmt.Sprintf("Failed to start MetaHandler RestAPI Server. %v", err))
	}
}
