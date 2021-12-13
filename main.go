package main

import (
	"os"

	"github.com/xylonx/x-telebot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
