package app

import "github.com/ghost-yu/go_shop_second/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
