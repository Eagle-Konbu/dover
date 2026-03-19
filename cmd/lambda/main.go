package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Eagle-Konbu/dover/application"
	"github.com/Eagle-Konbu/dover/infrastructure/discord"
	"github.com/Eagle-Konbu/dover/infrastructure/octopus"
	"github.com/Eagle-Konbu/dover/infrastructure/secretsmanager"
)

type octopusCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handler(ctx context.Context) error {
	secrets, err := secretsmanager.New(ctx)
	if err != nil {
		return fmt.Errorf("init secrets manager: %w", err)
	}

	raw, err := secrets.GetSecretValue(ctx, os.Getenv("OCTOPUS_SECRET_NAME"))
	if err != nil {
		return fmt.Errorf("get octopus credentials: %w", err)
	}

	var creds octopusCredentials
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return fmt.Errorf("parse octopus credentials: %w", err)
	}

	energy := &octopus.Client{
		Email:         creds.Email,
		Password:      creds.Password,
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
