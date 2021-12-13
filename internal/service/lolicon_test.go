package service_test

import (
	"testing"

	"github.com/xylonx/x-telebot/internal/service"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

func TestRegular(t *testing.T) {
	cases := []struct {
		Text      string
		SubString []string
	}{
		{"来三张色图", nil},
		{"来3张萝莉色图", []string{"来3张萝莉色图", "3", "萝莉", "色"}},
		{"来3张涩图", []string{"来3张黑丝涩图", "3", "", "涩"}},
	}

	for i := range cases {
		m := service.PornPicReg.FindStringSubmatch(cases[i].Text)
		eq := stringSliceEqual(m, cases[i].SubString)
		zapx.Info("match result", zap.Bool("match", eq), zap.Strings("result", m))
		if !eq {
			t.FailNow()
		}
	}
}

func stringSliceEqual(src, dst []string) bool {
	if len(src) != len(dst) {
		return false
	}

	for i := range src {
		if src[i] != dst[i] {
			return false
		}
	}

	return true
}
