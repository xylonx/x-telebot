package core

import (
	"github.com/xylonx/x-telebot/pkg/session"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

var (
	SessionManager session.SessionManager
)

func Setup() (err error) {
	SessionManager, err = session.NewSQLLiteSessionManager("x-telebot.db")
	if err != nil {
		zapx.Error("new sqlite session manager failed", zap.Error(err))
		return err
	}

	return nil
}
