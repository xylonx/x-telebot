package bot

import (
	"sync"
	"time"

	"github.com/xylonx/x-telebot/pkg/session"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gopkg.in/tucnak/telebot.v2"
)

// CronHandler - do cron jobs by passing Bot and the session about a communication
type CronHandler func(*telebot.Bot, session.SessionManager)

type CronJob struct {
	Handler CronHandler
	Ticker  *time.Ticker
}

type CallbackHandler func(*telebot.Bot, *telebot.Message, session.SessionManager)

type Bot struct {
	TelegramBotClient *telebot.Bot

	// cronJobs
	cronRWlock *sync.RWMutex
	cronJobs   []CronJob

	// messageHandler
	// telegramMessageEndpoint -> many message handles
	handlerRWLock *sync.RWMutex
	handlers      map[interface{}][]CallbackHandler
}

type Option struct {
	BotAPIToken string
}

func NewBot(opt *Option) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		URL:   "https://api.telegram.org",
		Token: opt.BotAPIToken,
		Poller: &telebot.LongPoller{
			Timeout: 5 * time.Second,
		},
	})

	if err != nil {
		zapx.Error("create telegram bot client failed", zap.Error(err))
		return nil, err
	}

	return &Bot{
		TelegramBotClient: bot,
		cronRWlock:        &sync.RWMutex{},
		cronJobs:          nil,
		handlerRWLock:     &sync.RWMutex{},
		handlers:          map[interface{}][]CallbackHandler{},
	}, nil
}

func (b *Bot) RegisterCronJob(endpint interface{}, cronjob CronHandler, ticker *time.Ticker) {
	b.cronRWlock.Lock()
	b.cronJobs = append(b.cronJobs, CronJob{
		Handler: cronjob,
		Ticker:  ticker,
	})
	b.cronRWlock.Unlock()
}

func (b *Bot) GetAllCronJobs() []CronJob {
	b.cronRWlock.RLock()
	jobs := b.cronJobs
	b.cronRWlock.RUnlock()
	return jobs
}

func (b *Bot) RegisterMessageHandler(endpint interface{}, handler CallbackHandler) {
	b.handlerRWLock.Lock()
	b.handlers[endpint] = append(b.handlers[endpint], handler)
	b.handlerRWLock.Unlock()
}

func (b *Bot) GetAllMessageHandler() map[interface{}][]CallbackHandler {
	b.handlerRWLock.RLock()
	handlers := b.handlers
	b.handlerRWLock.RUnlock()
	return handlers
}
