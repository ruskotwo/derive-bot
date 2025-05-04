package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ruskotwo/derive-bot/cmd/factory"
)

var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "Start telegram bot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start telegram bot")

		bot, cleaner, err := factory.InitTelegramBot()
		defer func() {
			if cleaner != nil {
				cleaner()
			}
		}()

		if err != nil {
			log.Printf("error initialize bot: %v", err)
			return
		}

		if err := bot.Start(); err != nil {
			log.Printf("error start bot: %v", err)
			return
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
	},
}

func init() {
	rootCmd.AddCommand(telegramCmd)
}
