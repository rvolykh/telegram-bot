package config

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	TelegramAPIToken           string
	SSMParamTelegramAPIToken   string
	SSMParamPowerScheduleCache string
	PowerScheduleCacheTTL      time.Duration
}

func LoadAppConfig() Config {
	return Config{
		TelegramAPIToken:           os.Getenv("TELEGRAM_APITOKEN"),
		SSMParamTelegramAPIToken:   os.Getenv("SSM_PARAM_TELEGRAM_APITOKEN"),
		SSMParamPowerScheduleCache: os.Getenv("SSM_PARAM_POWERON_CACHE"),
		PowerScheduleCacheTTL:      5 * time.Minute,
	}
}

func LoadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}
