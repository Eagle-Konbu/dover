# dover

[![CI](https://github.com/Eagle-Konbu/dover/actions/workflows/ci.yml/badge.svg)](https://github.com/Eagle-Konbu/dover/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Eagle-Konbu/dover/branch/main/graph/badge.svg)](https://codecov.io/gh/Eagle-Konbu/dover)

## Overview

A Discord bot that collects electric usage information and sends a daily notification.

This is a small, personal-use bot designed for simplicity and maintainability.

---

## How It Works

1. Amazon EventBridge triggers a Lambda function once per day
2. Fetches electric usage information
3. Formats the collected data
4. Sends a notification via Discord Webhook

---

## Architecture

- Language: Go
- Runtime: AWS Lambda
- Scheduler: Amazon EventBridge
- Notification: Discord Webhook

---

## Execution Schedule

- Frequency: Once per day
- Trigger: Amazon EventBridge scheduled execution

---

## Environment Variables

| Name | Description |
| ---- | ----------- |
| DISCORD_WEBHOOK_URL | Discord Webhook URL |

---

## Notes

- This bot is intended for personal use
- Source website structure changes may cause data retrieval failures

---

## License

MIT
