package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Received SQS Event with %d records", len(sqsEvent.Records))

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(cfg)

	// Process each SQS record
	for i, record := range sqsEvent.Records {
		log.Printf("Processing record %d:", i+1)
		log.Printf("  Message ID: %s", record.MessageId)
		log.Printf("  Receipt Handle: %s", record.ReceiptHandle)
		log.Printf("  Source ARN: %s", record.EventSourceARN)
		log.Printf("  Body: %s", record.Body)

		// Parse Telegram update from SQS message body
		if record.Body != "" {
			var update tgbotapi.Update
			if err := json.Unmarshal([]byte(record.Body), &update); err != nil {
				log.Printf("Error parsing Telegram update: %v", err)
				log.Printf("Raw body: %s", record.Body)
				continue
			}

			log.Printf("Parsed Telegram Update:")
			log.Printf("  Update ID: %d", update.UpdateID)

			if update.Message == nil {
				log.Printf("Skipping update: no message")
				continue
			}

			if !update.Message.IsCommand() {
				log.Printf("Skipping message: not a command")
				continue
			}

			log.Printf("  Message:")
			log.Printf("    Message ID: %d", update.Message.MessageID)
			log.Printf("    Date: %d", update.Message.Date)
			log.Printf("    Text: %s", update.Message.Text)
			log.Printf("    IsCommand: %v", update.Message.IsCommand())
			log.Printf("    Command: %s", update.Message.Command())
			log.Printf("    Command Arguments: %s", update.Message.CommandArguments())
			if update.Message.From != nil {
				log.Printf("    From:")
				log.Printf("      ID: %d", update.Message.From.ID)
				log.Printf("      Username: %s", update.Message.From.UserName)
				log.Printf("      First Name: %s", update.Message.From.FirstName)
				log.Printf("      Last Name: %s", update.Message.From.LastName)
				log.Printf("      Is Bot: %v", update.Message.From.IsBot)
			}

			if update.Message.Chat == nil {
				log.Printf("Skipping message: no chat")
				continue
			}

			log.Printf("    Chat:")
			log.Printf("      ID: %d", update.Message.Chat.ID)
			log.Printf("      Type: %s", update.Message.Chat.Type)
			log.Printf("      Title: %s", update.Message.Chat.Title)
			log.Printf("      Username: %s", update.Message.Chat.UserName)

			queueUrl := os.Getenv(fmt.Sprintf("SQS_COMMAND_%s_QUEUE_URL", strings.ToUpper(update.Message.Command())))
			if queueUrl == "" {
				log.Printf("No SQS queue URL found for command: %s, Skipping %d", update.Message.Command(), update.UpdateID)
				continue
			}

			resp, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:               aws.String(queueUrl),
				MessageBody:            aws.String(record.Body),
				MessageGroupId:         aws.String(fmt.Sprintf("%d", update.Message.Chat.ID)),
				MessageDeduplicationId: aws.String(fmt.Sprintf("%d", update.UpdateID)),
			})
			if err != nil {
				log.Printf("Error sending message to SQS: %v", err)
				continue
			}
			log.Printf("Message %d sent to SQS: %s / %s", update.UpdateID, queueUrl, *resp.MessageId)
		}
	}

	// SQS events don't require a response
	return nil
}

func main() {
	lambda.Start(handler)
}
