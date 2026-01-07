package subscriptions

import (
	"context"

	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions/dynamodb"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions/models"
)

type SubscriptionsDB interface {
	InsertSubscription(ctx context.Context, chatID int64, group string) error
	GetSubscription(ctx context.Context, chatID int64) ([]string, error)
	UpdateSubscription(ctx context.Context, chatID int64, group string) error
	DeleteSubscription(ctx context.Context, chatID int64) error
	ListSubscriptions(ctx context.Context) ([]models.Subscription, error)
}

func NewSubscriptionsDB(cfg *config.Config) (SubscriptionsDB, error) {
	return dynamodb.NewDynamoDB(cfg), nil
}
