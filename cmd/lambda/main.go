package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Eagle-Konbu/dover/application"
	"github.com/Eagle-Konbu/dover/infrastructure/discord"
	"github.com/Eagle-Konbu/dover/infrastructure/octopus"
)

func handler(ctx context.Context) error {
	energy := &octopus.Client{
		Email:         os.Getenv("OCTOPUS_EMAIL"),
		Password:      os.Getenv("OCTOPUS_PASSWORD"),
		APIURL:        os.Getenv("OCTOPUS_API_URL"),
		AccountNumber: os.Getenv("OCTOPUS_ACCOUNT_NUMBER"),
	}
	notifier := &discord.Webhook{
		URL: os.Getenv("DISCORD_WEBHOOK_URL"),
	}
	svc := &application.ReportService{
		Energy:   energy,
		Notifier: notifier,
	}
	return svc.Run(ctx)
}

func main() {
	log.SetFlags(0)
	lambda.Start(handler)
}
