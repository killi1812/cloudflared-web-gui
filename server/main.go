package main

import (
	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/controller"
	"github.com/killi1812/cloudflared-web-gui/service"
	"github.com/killi1812/cloudflared-web-gui/util/seed"

	"go.uber.org/zap"
)

//	@securitydefinitions.bearerauth	BearerAuth

func init() {
	app.Setup()
}

func main() {
	// Provide logger
	app.Provide(zap.S)

	app.Provide(service.NewUserCrudService)
	app.Provide(service.NewAuthService)
	app.Provide(service.NewTunelSrv)

	app.RegisterController(controller.NewInfoCnt)
	app.RegisterController(controller.NewUserCtn)
	app.RegisterController(controller.NewAuthCtn)
	app.RegisterController(controller.NewTunnelCtn)

	seed.Insert()

	app.Start()
}
