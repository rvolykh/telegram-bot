package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

const (
	EnvAPIToken                   = "TELEGRAM_APITOKEN"
	EnvSSMParamAPIToken           = "SSM_PARAM_TELEGRAM_APITOKEN"
	EnvSSMParamCache              = "SSM_PARAM_CACHE"
	EnvDynamoDBSubscriptionsTable = "DYNAMODB_TABLE_SUBSCRIPTIONS"

	PowerScheduleCacheTTL = 5 * time.Minute
)

type Config struct {
	AWSConfig                  aws.Config
	TelegramAPIToken           string
	SSMParamCache              string
	DynamoDBSubscriptionsTable string
	PowerScheduleCacheTTL      time.Duration
}

func LoadConfig(ctx context.Context) (*Config, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	apiToken := os.Getenv(EnvAPIToken)
	if apiToken == "" {
		response, err := ssm.NewFromConfig(awsConfig).GetParameter(ctx, &ssm.GetParameterInput{
			Name:           aws.String(os.Getenv(EnvSSMParamAPIToken)),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get Telegram API token from SSM parameter: %w", err)
		}
		apiToken = *response.Parameter.Value
	}

	return &Config{
		AWSConfig:                  awsConfig,
		TelegramAPIToken:           apiToken,
		SSMParamCache:              os.Getenv(EnvSSMParamCache),
		DynamoDBSubscriptionsTable: os.Getenv(EnvDynamoDBSubscriptionsTable),
		PowerScheduleCacheTTL:      PowerScheduleCacheTTL,
	}, nil
}
