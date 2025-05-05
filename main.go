package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	qqbot "github.com/ExquisiteCore/LagrangeGo-Template/bot"
	"github.com/ExquisiteCore/LagrangeGo-Template/config"
	"github.com/ExquisiteCore/LagrangeGo-Template/logic"
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
	tgbot "github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
)

// 创建 protocolLogger 实例
var logger = utils.ProtocolLogger{}

func init() {
	config.Init()
	utils.Init()
	qqbot.Init(&logger)
}

func main() {

	// QQ Login
	qqbot.Login()
	qqbot.Listen()

	// TG Login
	b, err := tgbot.New(
		config.GlobalConfig.TGBot.Token,
		tgbot.WithDefaultHandler(logic.TGSetUpHandler),
	)
	if err != nil {
		logrus.Errorf("Telegram bot login failed: %s\n", err)
		return
	}

	logic.RegisterCustomLogic(b)
	logic.SetupLogic()
	defer qqbot.QQClient.Release()
	defer qqbot.Dumpsig()

	b.Start(context.TODO())

	// setup the main stop channel
	mc := make(chan os.Signal, 2)
	signal.Notify(mc, os.Interrupt, syscall.SIGTERM)
	for {
		switch <-mc {
		case os.Interrupt, syscall.SIGTERM:
			return
		}
	}
}
