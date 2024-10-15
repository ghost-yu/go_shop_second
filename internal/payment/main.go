package main

import (
	"github.com/ghost-yu/go_shop_second/common/config"
	"github.com/ghost-yu/go_shop_second/common/logging"
	"github.com/ghost-yu/go_shop_second/common/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serverType := viper.GetString("payment.server-to-run")

	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
