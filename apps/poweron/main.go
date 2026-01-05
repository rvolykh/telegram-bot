package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/operations"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
)

type LambdaFunction struct {
	SSMClient *ssm.Client
	Config    config.Config
}

func (f *LambdaFunction) Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Received SQS Event with %d records", len(sqsEvent.Records))

	apiToken, err := f.getAPIToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get API token: %w", err)
	}

	telegram, err := reply.NewTelegram(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create telegram: %w", err)
	}

	for i, record := range sqsEvent.Records {
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
		case "":
			err = operations.ShowMainMenu(ctx, telegram, chatID)
			if err != nil {
				log.Printf("Error showing main menu: %v", err)
				continue
			}
			log.Printf("Main menu sent to chat: %d", chatID)

		case "Сьогодні":
			err = operations.Today(ctx, f.Config, telegram, f.SSMClient, chatID)
			if err != nil {
				log.Printf("Error showing today schedule: %v", err)
				continue
			}
			log.Printf("Today schedule sent to chat: %d", chatID)

		case "Завтра":
			err = operations.Tomorrow(ctx, f.Config, telegram, f.SSMClient, chatID)
			if err != nil {
				log.Printf("Error showing tomorrow schedule: %v", err)
				continue
			}
			log.Printf("Tomorrow schedule sent to chat: %d", chatID)

		case "Підписатись":
			err = operations.ShowSelectGroupMenu(ctx, telegram, chatID)
			if err != nil {
				log.Printf("Error subscribing: %v", err)
				continue
			}
			log.Printf("Select group menu sent to chat: %d", chatID)

		case "Відписатись":
			err = operations.Unsubscribe(ctx, f.Config, telegram, f.SSMClient, chatID)
			if err != nil {
				log.Printf("Error unsubscribing: %v", err)
				continue
			}
			log.Printf("Unsubscribed from chat: %d", chatID)

		case "Закрити":
			err = operations.CloseMenu(ctx, telegram, chatID)
			if err != nil {
				log.Printf("Error closing menu: %v", err)
				continue
			}
			log.Printf("Closed menu in chat: %d", chatID)

		case "1.1", "1.2", "2.1", "2.2", "3.1", "3.2", "4.1", "4.2", "5.1", "5.2", "6.1", "6.2":
			err = operations.Subscribe(ctx, f.Config, telegram, f.SSMClient, chatID, args)
			if err != nil {
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

func (f *LambdaFunction) getAPIToken(ctx context.Context) (string, error) {
	if f.Config.TelegramAPIToken != "" {
		return f.Config.TelegramAPIToken, nil
	}

	apiToken, err := f.SSMClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(f.Config.SSMParamTelegramAPIToken),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get SSM parameter: %w", err)
	}
	return *apiToken.Parameter.Value, nil
}

func main() {
	cfg, err := config.LoadAWSConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load AWS config: %v", err)
	}

	fn := &LambdaFunction{
		SSMClient: ssm.NewFromConfig(cfg),
		Config:    config.LoadAppConfig(),
	}
	lambda.Start(fn.Handler)
}
