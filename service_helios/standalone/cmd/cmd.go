package main

import (
	"github.com/alecthomas/kong"
	_ "github.com/go-sql-driver/mysql"
	logrus "github.com/helios/go-sdk/proxy-libs/helioslogrus"
	"microservice/config"
	"microservice/service_helios/standalone/serve"
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
