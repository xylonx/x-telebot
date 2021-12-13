package service

import (
	"context"
	"regexp"
	"strconv"

	"github.com/xylonx/x-telebot/internal/bot"
	"github.com/xylonx/x-telebot/internal/config"
	"github.com/xylonx/x-telebot/internal/picture/lolicon"
	"github.com/xylonx/x-telebot/pkg/session"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gopkg.in/tucnak/telebot.v2"
)

var porn = map[string]bool{"色": false, "涩": true}

var (
	// 来(2)张(黑丝)色图  -> general
	// 来(5)张(萝莉)涩图  -> pornographic
	PornPicReg *regexp.Regexp
)

func init() {
	var err error

	PornPicReg, err = regexp.Compile(`来(\d+)张(.*)(涩|色)图`)
	if err != nil {
		zapx.Error("compile regular string failed", zap.Error(err))
	}
}

func registerPictureService(bot *bot.Bot) {
	bot.RegisterMessageHandler(telebot.OnText, func(b *telebot.Bot, msg *telebot.Message, _ session.SessionManager) {
		subStr := PornPicReg.FindStringSubmatch(msg.Text)
		if subStr == nil {
			return
		}

		zapx.Info("text match pornographic regexp", zap.Any("sender", msg.Sender))

		num, _ := strconv.ParseInt(subStr[1], 10, 16)
		tag := subStr[2]
		r18 := porn[subStr[3]]

		tags := make([]string, 1)
		if tag != "" {
			tags[0] = tag
		} else {
			tags = nil
		}

		pics, err := lolicon.DefaultLoliconClient.GetPictures(context.Background(), r18, int(num), tags)
		if err != nil {
			return
		}

		if config.Config.Picture.Sync {
			// TODO: sync using syncer
		}

		b.Send(msg.Chat, "少女祈祷中....")

		photos := make(telebot.Album, 0, len(pics.Data))
		for i := range pics.Data {
			photos = append(photos, &telebot.Photo{File: telebot.File{FileURL: pics.Data[i].URLs.Original}})
		}

		_, err = b.SendAlbum(msg.Chat, photos)
		if err != nil {
			zapx.Error("album message send failed", zap.Error(err))
			return
		}
	})
}
