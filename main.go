package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/brxie/OdpadyAlertBot/config"
	"github.com/brxie/OdpadyAlertBot/notification"
	"github.com/brxie/OdpadyAlertBot/pkg/client/telegram"
	"github.com/brxie/OdpadyAlertBot/scheduler"
)

//go:embed assets/logo.png
var logo []byte

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Telegram bot token not set")
	}

	client, err := telegram.NewTelegramClient(token, nil)
	if err != nil {
		log.Fatal(err)
	}

	regions, err := config.LoadRegions()
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error)
	for _, r := range regions {
		events, err := config.LoadEvents(r.CalendarFile)
		if err != nil {
			log.Fatal(err)
		}
		go scheduler.Run(r, events, errCh, notification.SchedulerCallback(client, r, logo))
	}

	for {
		err = <-errCh
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
	}
}
