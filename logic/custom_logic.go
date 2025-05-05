package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/ExquisiteCore/LagrangeGo-Template/config"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"

	qqbot "github.com/ExquisiteCore/LagrangeGo-Template/bot"
	tgbot "github.com/go-telegram/bot"
	//"github.com/sirupsen/logrus"
)

// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic(b *tgbot.Bot) {
	// 注册私聊消息处理逻辑
	Manager.RegisterPrivateMessageHandler(
		func(client *client.QQClient, event *message.PrivateMessage) {
			client.SendPrivateMessage(
				event.Sender.Uin,
				[]message.IMessageElement{message.NewText("だまれ！")},
			)
		},
	)

	// 注册群消息处理逻辑
	Manager.RegisterGroupMessageHandler(func(client *client.QQClient, event *message.GroupMessage) {
		logrus.Infof("QQ Group Number: %d\n", event.GroupUin)
		if tgn, ok := config.GlobalQQTGMap[event.GroupUin]; ok {
			b.SendMessage(context.TODO(), &tgbot.SendMessageParams{
				ChatID: tgn,
				Text:   fmt.Sprintf("[%s]: %s", event.Sender.Nickname, event.ToString()),
			})
		}
	})
}

func TGSetUpHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	logrus.Infof("TG Group Number: %d\n", update.Message.Chat.ID)
	if qqn, ok := config.GlobalTGQQMap[update.Message.Chat.ID]; ok {
		var s strings.Builder
		s.WriteString("[")
		if update.Message.From.LastName != "" {
			s.WriteString(update.Message.From.LastName)
			s.WriteString(" ")
		}
		s.WriteString(update.Message.From.FirstName)
		s.WriteString("]: ")
		s.WriteString(update.Message.Text)
		qqbot.QQClient.SendGroupMessage(
			qqn,
			[]message.IMessageElement{
				message.NewText(
					s.String(),
				),
			},
		)
	}
}
