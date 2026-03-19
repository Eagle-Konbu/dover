package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Eagle-Konbu/dover/application"
	"github.com/Eagle-Konbu/dover/infrastructure/discord"
	"github.com/Eagle-Konbu/dover/infrastructure/octopus"
)

func main() {
	log.SetFlags(0)
	send := flag.Bool("send", false, "send notification to Discord webhook")
	flag.Parse()

	energy := &octopus.Client{
		Email:         os.Getenv("OCTOPUS_EMAIL"),
		Password:      os.Getenv("OCTOPUS_PASSWORD"),
		APIURL:        os.Getenv("OCTOPUS_API_URL"),
		AccountNumber: os.Getenv("OCTOPUS_ACCOUNT_NUMBER"),
	}

	svc := &application.ReportService{Energy: energy}
	if *send {
		svc.Notifier = &discord.Webhook{URL: os.Getenv("DISCORD_WEBHOOK_URL")}
	} else {
		svc.Notifier = &stdoutNotifier{}
	}

	if err := svc.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
