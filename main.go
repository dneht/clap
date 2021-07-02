package main

import (
	"cana.io/clap/api"
	"cana.io/clap/pkg/base"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {
	base.Init()
	app := base.App
	conf := base.Now()
	app.Use("/api", func(c *fiber.Ctx) error {
		return api.CheckAuth(c)
	})
	app.Use("/select", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Static("/", conf.Service.StaticPath)
	app.Get("/app/**", func(c *fiber.Ctx) error {
		return c.SendFile(conf.Service.StaticPath + "/index.html")
	})
	app.Options("/**", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	app.Get("/clean/all", api.CleanAll)

	app.Get("/api/prop/base", api.PropBase)

	app.Post("/api/env/list", api.ListEnv)
	app.Get("/api/env/simple", api.SimpleEnv)
	app.Get("/api/env/:id", api.GetEnv)

	app.Post("/api/space/list", api.ListSpace)
	app.Get("/api/space/simple", api.SimpleSpace)
	app.Get("/api/space/:id", api.GetSpace)

	app.Post("/api/app/list", api.ListApp)
	app.Get("/api/app/simple", api.SimpleApp)
	app.Get("/api/app/:id", api.GetApp)
	app.Post("/api/app", api.CreateApp)
	app.Put("/api/app/:id", api.UpdateApp)
	app.Delete("/api/app/:id", api.DeleteApp)

	app.Post("/api/render/list", api.ListTemplate)
	app.Get("/api/render/simple", api.SimpleTemplate)
	app.Get("/api/render/:id", api.GetTemplate)
	app.Put("/api/render/:id", api.GetTemplate)
	app.Delete("/api/render/:id", api.GetTemplate)

	app.Post("/api/deploy/list", api.ListDeploy)
	app.Get("/api/deploy/:id", api.GetDeploy)
	app.Post("/api/deploy", api.CreateDeploy)
	app.Put("/api/deploy/:id", api.UpdateDeploy)

	app.Post("/api/config/list", api.ListDeploy)

	app.Get("/pod/space/:id", api.ListSpacePod)
	app.Get("/pod/deploy/:id", api.ListDeployPod)
	app.Get("/pod/restart/:id", api.RestartDeployPod)

	// type=check, build, deploy
	app.Get("/deploy/:type/:deploy", api.ExecDeploy)
	// type=jsonnet
	app.Get("/render/:type/:deploy", api.ExecRender)
	// type=exec, attach, inner
	app.Get("/select/:type/:pod", websocket.New(api.ExecSelect))

	log.Fatal(app.Listen(":" + conf.Service.Port))
}
