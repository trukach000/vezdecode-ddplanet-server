package main

import (
	"ddplanet-server/app"
	"ddplanet-server/docs"
	"ddplanet-server/pkg/config"
	"fmt"
	"net/http"

	_ "ddplanet-server/docs"

	"github.com/sirupsen/logrus"
)

// @title DDPlanet Server API
// @version 0.1
// @description An API for front-end
// @Schemes http https
func main() {
	config.Init()
	r := app.Setup()
	cfg := config.Get()

	docs.SwaggerInfo.BasePath = fmt.Sprintf("%s/api/v1", cfg.GlobalPrefix)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	var listenErr error
	if cfg.UseSSL {
		logrus.Infof("Starting listening server with SSL on %s", addr)
		listenErr = http.ListenAndServeTLS(addr, "cert.pem", "privkey.pem", r)
	} else {
		logrus.Infof("Starting listening server without SSL on %s", addr)
		listenErr = http.ListenAndServe(addr, r)
	}
	logrus.Fatalf("Listen err :%s ", listenErr)
}
