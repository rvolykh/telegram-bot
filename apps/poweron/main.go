package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	PowerScheduleCacheTTL      = 5 * time.Minute
	TelegramAPIToken           = "TELEGRAM_APITOKEN"
	SSMParamTelegramAPIToken   = "SSM_PARAM_TELEGRAM_APITOKEN"
	SSMParamPowerScheduleCache = "SSM_PARAM_POWERON_CACHE"
)

type LambdaFunction struct {
	SSMClient *ssm.Client
}

func (f *LambdaFunction) Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Received SQS Event with %d records", len(sqsEvent.Records))

	apiToken, err := f.getAPIToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get API token: %w", err)
	}
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	schedule, err := f.getPowerSchedule(ctx)
	if err != nil {
		return fmt.Errorf("failed to get power on: %w", err)
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

		log.Printf("Parsed Telegram Update:")
		log.Printf("  Update ID: %d", update.UpdateID)

		if update.Message == nil {
			log.Printf("Skipping update: no message")
			continue
		}

		if update.Message.Chat == nil {
			log.Printf("Skipping message: no chat")
			continue
		}

		chatID := update.Message.Chat.ID

		_, err := bot.Send(tgbotapi.NewMessage(chatID, schedule))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}
		log.Printf("Message sent to user: %s", update.FromChat().UserName)
	}

	return nil
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

func (f *LambdaFunction) getPowerSchedule(ctx context.Context) (string, error) {
	powerSchedule, err := f.SSMClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(os.Getenv(SSMParamPowerScheduleCache)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to get power schedule from cache: %v", err)
	} else {
		isCacheValid := powerSchedule.Parameter.Value != nil &&
			*powerSchedule.Parameter.Value != "none" &&
			powerSchedule.Parameter.LastModifiedDate.After(time.Now().Add(-PowerScheduleCacheTTL))

		if isCacheValid {
			log.Printf("Power schedule is still valid, returning cached value")
			return *powerSchedule.Parameter.Value, nil
		}
		log.Printf("Power schedule is outdated, refreshing cache")
	}

	powerScheduleText, err := getPowerSchedule()
	if err != nil {
		return "", fmt.Errorf("failed to get power schedule: %w", err)
	}

	_, err = f.SSMClient.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(os.Getenv(SSMParamPowerScheduleCache)),
		Value:     aws.String(powerScheduleText),
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to put power schedule to cache: %v", err)
	} else {
		log.Printf("Updated power schedule in cache")
	}

	return powerScheduleText, nil
}

type PowerOnMenuItem struct {
	Name          string `json:"name"`
	RawMobileHTML string `json:"rawMobileHtml"`
}

type PowerOnMember struct {
	MenuItems []PowerOnMenuItem `json:"menuItems"`
}

type PowerOnResponse struct {
	Member []PowerOnMember `json:"hydra:member"`
}

func getPowerSchedule() (string, error) {
	resp, err := http.Get("https://api.loe.lviv.ua/api/menus?page=1&type=photo-grafic")
	if err != nil {
		return "", fmt.Errorf("failed to get power on: %w", err)
	}
	defer resp.Body.Close()

	var powerOnResponse PowerOnResponse
	err = json.NewDecoder(resp.Body).Decode(&powerOnResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode power on response: %w", err)
	}

	if len(powerOnResponse.Member) == 0 {
		return "", fmt.Errorf("no power on found")
	}
	if len(powerOnResponse.Member[0].MenuItems) == 0 {
		return "", fmt.Errorf("no menu items found")
	}

	text := cleanHTML(powerOnResponse.Member[0].MenuItems[0].RawMobileHTML)
	if text == "" {
		text = "Немає запланованих відключень електроенергії на сьогодні."
	}

	return text, nil
}

func cleanHTML(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`
	for i, c := range text {
		if (i+1) == len(text) && end >= start {
			builder.WriteString(text[end:])
		}
		if c != '<' && c != '>' {
			continue
		}
		if c == '<' {
			if !in {
				start = i
				builder.WriteString(text[end:start])
			}
			in = true
			continue
		}
		end, in = i+1, false
	}

	return builder.String()
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
