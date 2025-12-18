package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TelegramAPIToken         = "TELEGRAM_APITOKEN"
	SSMParamTelegramAPIToken = "SSM_PARAM_TELEGRAM_APITOKEN"
	TelegramChatID           = "TELEGRAM_CHAT_ID"
)

type LambdaFunction struct {
	SSMClient *ssm.Client
}

func (f *LambdaFunction) Handler(ctx context.Context, event events.SNSEvent) error {
	log.Printf("Received events: %d", len(event.Records))

	apiToken, err := f.getAPIToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get API token: %w", err)
	}
	chatID := os.Getenv(TelegramChatID)
	if chatID == "" {
		return fmt.Errorf("failed to get chat ID")
	}
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chat ID: %w", err)
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	for _, record := range event.Records {
		handlerFn := handleUnexpectedSource
		if strings.Contains(record.SNS.Message, "AlarmName") {
			handlerFn = handleCloudWatchAlert
		}
		text, err := handlerFn(record)
		if err != nil {
			return fmt.Errorf("failed to handle alert: %w", err)
		}

		message := tgbotapi.NewMessage(chatIDInt, text)
		message.ParseMode = tgbotapi.ModeHTML

		_, err = bot.Send(message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		log.Printf("Message sent to user: %d", chatIDInt)
	}

	return nil
}

func handleCloudWatchAlert(record events.SNSEventRecord) (string, error) {
	var alert = make(map[string]any)
	err := json.Unmarshal([]byte(record.SNS.Message), &alert)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal alert: %w", err)
	}

	state, err := getStringValueFromKV(alert, "NewStateValue")
	if err != nil {
		return "", fmt.Errorf("failed to get new state value: %w", err)
	}

	changedAt, err := getStringValueFromKV(alert, "StateChangeTime")
	if err != nil {
		return "", fmt.Errorf("failed to get state change time: %w", err)
	}

	reason, err := getStringValueFromKV(alert, "NewStateReason")
	if err != nil {
		return "", fmt.Errorf("failed to get new state reason: %w", err)
	}

	var icon = "🔴"
	if state == "OK" {
		icon = "🟢"
	}

	text := fmt.Sprintf(
		"%s <b>Alert</b>: <u>%s</u>!\n\n%s\n%s",
		icon, record.SNS.Subject, changedAt, reason,
	)
	return text, nil
}

func handleUnexpectedSource(record events.SNSEventRecord) (string, error) {
	log.Println("Unexpected source format", record.SNS.Message)

	text := fmt.Sprintf(
		"🆘 <b>Alert</b>: <u>%s</u>!\n\n%s",
		record.SNS.Subject, record.SNS.Message,
	)
	return text, nil
}

func (f *LambdaFunction) getAPIToken(ctx context.Context) (string, error) {
	if apiToken, ok := os.LookupEnv(TelegramAPIToken); ok {
		return apiToken, nil
	}

	apiToken, err := f.SSMClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(os.Getenv(SSMParamTelegramAPIToken)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get SSM parameter: %w", err)
	}
	return *apiToken.Parameter.Value, nil
}

func getStringValueFromKV(kv map[string]any, key string) (string, error) {
	if value, ok := kv[key]; ok {
		if value, ok := value.(string); ok {
			return value, nil
		}
		return "", fmt.Errorf("value is not a string")
	}
	return "", fmt.Errorf("failed to get value from map")
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load SDK config: %v", err)
	}

	fn := &LambdaFunction{
		SSMClient: ssm.NewFromConfig(cfg),
	}
	lambda.Start(fn.Handler)
}
