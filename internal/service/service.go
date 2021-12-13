package service

import (
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/xylonx/x-telebot/internal/bot"
	"github.com/xylonx/x-telebot/internal/config"
	"github.com/xylonx/x-telebot/internal/core"
	"gopkg.in/tucnak/telebot.v2"
)

var (
	tbot        *bot.Bot
	cronJobPool *ants.Pool
)

func Start() (err error) {
	tbot, err = bot.NewBot(&bot.Option{
		BotAPIToken: config.Config.TelegramBot.ApiKey,
	})
	if err != nil {
		return err
	}

	cronJobPool, err = ants.NewPool(int(config.Config.Application.CronjobPoolSize), ants.WithNonblocking(false))
	if err != nil {
		return err
	}

	// register service here
	registerPictureService(tbot)

	go func() {
		tbot.TelegramBotClient.Start()
	}()

	// do cronJobs
	go func() {
		for {
			time.Sleep(time.Second)
			jobs := tbot.GetAllCronJobs()
			for i := range jobs {
				select {
				case <-jobs[i].Ticker.C:
					cronJobPool.Submit(func() {
						jobs[i].Handler(tbot.TelegramBotClient, core.SessionManager)
					})
				default:
				}
			}
		}
	}()

	// register handlers
	handlers := tbot.GetAllMessageHandler()
	for ep := range handlers {
		switch ep {
		case telebot.OnQuery:
			tbot.TelegramBotClient.Handle(ep, func(q *telebot.Query) {
				// TODO:
			})
		default:
			tbot.TelegramBotClient.Handle(ep, func(msg *telebot.Message) {
				for i := range handlers[ep] {
					handlers[ep][i](tbot.TelegramBotClient, msg, core.SessionManager)
				}
			})
		}
	}

	return nil
}

func Stop() {
	tbot.TelegramBotClient.Stop()
}
