package main

import (
	"cana.io/clap/api"
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	base.Init()
	app := base.App
	conf := base.Now()

	app.Static("/", conf.Service.StaticPath)
	static := func(c *fiber.Ctx) error {
		return c.SendFile(conf.Service.StaticPath + "/index.html")
	}
	app.Get("/login", static)
	app.Get("/app/**", static)
	app.Options("/**", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Use("/select", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Use("/api", func(c *fiber.Ctx) error {
		return api.CheckAuth(c)
	})
	app.Post(api.LoginApi, api.LoginUser)
	app.Get(api.ConfigApi, api.ConfigBase)
	app.Get(api.StaticApi, api.StaticRes)
	app.Get(api.CleanApi, api.CleanAll)

	app.Post(api.ResApiPre+api.ListSuffix, api.ListRes)
	app.Get(api.ResApiPre+api.SimpleSuffix, api.SimpleRes)
	app.Get(api.ResApiPre+"/:id", api.GetRes)
	app.Post(api.PowApiPre+api.ListSuffix, api.ListPow)
	app.Post(api.PowApiPre+api.SimpleSuffix, api.SimplePow)
	app.Get(api.PowApiPre+"/:id", api.GetPow)

	app.Post(api.UserApiPre+api.ListSuffix, api.ListUser)
	app.Get(api.UserApiPre+"/:id", api.GetUser)
	app.Put(api.UserApiPre, api.CreateUser)
	app.Post(api.UserApiPre+"/:id", api.ChangeUser)

	app.Post(api.RoleApiPre+api.ListSuffix, api.ListRole)
	app.Get(api.RoleApiPre+api.SimpleSuffix, api.SimpleRole)
	app.Get(api.RoleApiPre+"/:id", api.GetRole)
	app.Put(api.RoleApiPre, api.CreateRole)

	app.Post(api.EnvApiPre+api.ListSuffix, api.ListEnv)
	app.Get(api.EnvApiPre+api.SimpleSuffix, api.SimpleEnv)
	app.Get(api.EnvApiPre+"/:id", api.GetEnv)

	app.Post(api.SpaceApiPre+api.ListSuffix, api.ListSpace)
	app.Get(api.SpaceApiPre+api.SimpleSuffix, api.SimpleSpace)
	app.Get(api.SpaceApiPre+"/:id", api.GetSpace)

	app.Post(api.AppApiPre+api.ListSuffix, api.ListApp)
	app.Get(api.AppApiPre+api.SimpleSuffix, api.SimpleApp)
	app.Get(api.AppApiPre+"/:id", api.GetApp)
	app.Put(api.AppApiPre, api.CreateApp)
	app.Post(api.AppApiPre+"/:id", api.UpdateApp)
	app.Delete(api.AppApiPre+"/:id", api.DeleteApp)

	app.Post(api.RenderApiPre+api.ListSuffix, api.ListTemplate)
	app.Get(api.RenderApiPre+api.SimpleSuffix, api.SimpleTemplate)
	app.Get(api.RenderApiPre+"/:id", api.GetTemplate)

	app.Post(api.DeployApiPre+api.ListSuffix, api.ListDeploy)
	app.Get(api.DeployApiPre+"/:id", api.GetDeploy)
	app.Put(api.DeployApiPre, api.CreateDeploy)
	app.Post(api.DeployApiPre+"/:id", api.UpdateDeploy)

	app.Get("/prop/:type/:id", api.GetProp)
	app.Put("/prop/:type", api.CreateProp)
	app.Post("/prop/:type/:id", api.UpdateProp)

	app.Get("/pod/space/:id", api.ListSpacePod)
	app.Get("/pod/deploy/:id", api.ListDeployPod)
	app.Get("/pod/restart/:id", api.RestartDeployPod)
	// type=check, build, deploy
	app.Get("/deploy/:type/:deploy", api.ExecDeploy)
	// type=jsonnet, property
	app.Get("/render/:type/:deploy", api.ExecRender)
	// type=exec, attach, inner
	app.Get("/select/:type/:pod", websocket.New(api.ExecSelect))

	api.InitTask()
	log.Panic(app.Listen(":" + conf.Service.Port))
}
