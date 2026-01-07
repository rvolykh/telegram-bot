package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/operations"
)

type LambdaFunction struct {
	Config *config.Config
}

func (f *LambdaFunction) Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Received SQS Event with %d records", len(sqsEvent.Records))

	if len(sqsEvent.Records) > 0 {
		return f.handleSQS(ctx, sqsEvent.Records)
	}

	return f.handleSchedule(ctx)
}

func (f *LambdaFunction) handleSQS(ctx context.Context, records []events.SQSMessage) error {
	for i, record := range records {
		log.Printf("Processing record %d:", i+1)
		log.Printf("  Message ID: %s", record.MessageId)
		log.Printf("  Receipt Handle: %s", record.ReceiptHandle)
		log.Printf("  Source ARN: %s", record.EventSourceARN)
		log.Printf("  Body: %s", record.Body)

		var update tgbotapi.Update
		if err := json.Unmarshal([]byte(record.Body), &update); err != nil {
			log.Printf("Error parsing Telegram update: %v", err)
			continue
		}

		if update.Message == nil {
			log.Printf("Skipping update: no message")
			continue
		}

		if update.Message.Chat == nil {
			log.Printf("Skipping message: no chat")
			continue
		}

		chatID := update.Message.Chat.ID

		args := update.Message.CommandArguments()
		switch args {
		case "Сьогодні":
			if err := operations.Today(ctx, f.Config, chatID); err != nil {
				log.Printf("Error showing today schedule: %v", err)
				continue
			}
			log.Printf("Today schedule sent to chat: %d", chatID)

		case "Завтра":
			if err := operations.Tomorrow(ctx, f.Config, chatID); err != nil {
				log.Printf("Error showing tomorrow schedule: %v", err)
				continue
			}
			log.Printf("Tomorrow schedule sent to chat: %d", chatID)

		case "Підписатись":
			if err := operations.ShowSelectGroupMenu(ctx, f.Config, chatID); err != nil {
				log.Printf("Error subscribing: %v", err)
				continue
			}
			log.Printf("Select group menu sent to chat: %d", chatID)

		case "Відписатись":
			if err := operations.Unsubscribe(ctx, f.Config, chatID); err != nil {
				log.Printf("Error unsubscribing: %v", err)
				continue
			}
			log.Printf("Unsubscribed from chat: %d", chatID)

		case "":
			if err := operations.ShowMainMenu(ctx, f.Config, chatID); err != nil {
				log.Printf("Error showing main menu: %v", err)
				continue
			}
			log.Printf("Main menu sent to chat: %d", chatID)

		case "Закрити":
			if err := operations.CloseMenu(ctx, f.Config, chatID); err != nil {
				log.Printf("Error closing menu: %v", err)
				continue
			}
			log.Printf("Closed menu in chat: %d", chatID)

		case "1.1", "1.2", "2.1", "2.2", "3.1", "3.2", "4.1", "4.2", "5.1", "5.2", "6.1", "6.2":
			if err := operations.Subscribe(ctx, f.Config, chatID, args); err != nil {
				log.Printf("Error showing group schedule: %v", err)
				continue
			}
			log.Printf("Subscribed to group %s in chat: %d", args, chatID)

		default:
			log.Printf("Unknown command args: %s", args)

		}
	}

	return nil
}

func (f *LambdaFunction) handleSchedule(ctx context.Context) error {
	return operations.DeliverScheduleUpdates(ctx, f.Config)
}

func main() {
	cfg, err := config.LoadConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load config: %v", err)
	}

	lambda.Start((&LambdaFunction{Config: cfg}).Handler)
}
