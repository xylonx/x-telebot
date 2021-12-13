package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/xylonx/x-telebot/internal/config"
	"github.com/xylonx/x-telebot/internal/core"
	"github.com/xylonx/x-telebot/internal/service"
	"github.com/xylonx/zapx"
)

var version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:     "x-telebot",
	Short:   "",
	Version: version,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = config.Setup(cfgFile)
		if err != nil {
			return err
		}

		err = core.Setup()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.default.yaml", "specify config file path")
}

func Execute() error {
	return rootCmd.Execute()
}

func run() (err error) {
	err = service.Start()
	if err != nil {
		return err
	}

	zapx.Info("service starts. Press Ctrl+C to stop")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig

	zapx.Info("receive SIGTERM. stopping service....")
	service.Stop()
	zapx.Info("service stopped")
	return nil
}
