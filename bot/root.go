package bot

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
)

var (
	Bot *tb.Bot
)

func Start() error {
	var err error
	setting := tb.Settings{
		Token:   viper.GetString("telegram.token"),
		Updates: 100,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second, AllowedUpdates: []string{
			"message",
			"chat_member",
			"inline_query",
			"callback_query",
		}},
		OnError: func(err error, context tb.Context) {
			fmt.Printf("%+v\n", err)
		},
	}
	if viper.GetString("telegram.proxy") != "" {
		setting.URL = viper.GetString("telegram.proxy")
	}
	Bot, err = tb.NewBot(setting)
	if err != nil {
		return err
	}
	RegisterCommands()
	RegisterHandle()
	Bot.Start()
	return nil
}

func RegisterCommands() {
	_ = Bot.SetCommands([]tb.Command{
		{
			Text:        StartCmd,
			Description: "Hello🙌",
		},
		{
			Text:        AddAdCmd,
			Description: "添加广告",
		},
		{
			Text:        AllAdCmd,
			Description: "查看广告",
		},
		{
			Text:        DelAdCmd,
			Description: "删除广告",
		},
	})
}

func RegisterHandle() {
	Bot.Handle(StartCmd, func(c tb.Context) error {
		// 发送消息并检查错误
		msg, err := Bot.Send(c.Chat(), "🙋 欢迎进群, 我是防广告机器人, 请勿发送广告, 昵称也不要带有推广信息, 谢谢. 重要信息已置顶, 请留意查看.")
		if err != nil {
			return err
		}

		// 设置定时器，30秒后删除消息
		time.AfterFunc(30*time.Second, func() {
			if err := Bot.Delete(msg); err != nil {
				fmt.Println("Error deleting the message:", err)
			}
		})

		return nil
	}, PreCmdMiddleware)

	creatorOnly := Bot.Group()
	creatorOnly.Use(CreatorCmdMiddleware)
	creatorOnly.Handle(AllAdCmd, AllAd)
	creatorOnly.Handle(AddAdCmd, AddAd)
	creatorOnly.Handle(DelAdCmd, DelAd)

	groupOnly := Bot.Group()
	groupOnly.Use(PreGroupMiddleware)
	groupOnly.Handle(tb.OnText, OnTextMessage)
	groupOnly.Handle(tb.OnSticker, OnStickerMessage)
	groupOnly.Handle(tb.OnPhoto, OnPhotoMessage)

	Bot.Handle(tb.OnChatMember, OnChatMemberMessage)
}
