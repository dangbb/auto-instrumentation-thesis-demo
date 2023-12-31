package main

import (
	"github.com/alecthomas/kong"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"microservice/config"
	"microservice/service/standalone/serve"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infof("Start standalone service")

	cliConfig := config.Config{}
	kongCtx := kong.Parse(&cliConfig)

	switch kongCtx.Command() {
	case "server":
		serve.Run(cliConfig)
	}
}
