package operations

import (
	"context"
	"fmt"

	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions"
)

func Unsubscribe(ctx context.Context, cfg *config.Config, chatID int64) error {
	t, err := reply.NewTelegram(cfg)
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	s, err := subscriptions.NewSubscriptionsDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions repository: %w", err)
	}

	if err := s.DeleteSubscription(ctx, chatID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	_, err = t.SendMessage(ctx, chatID, "Ви відписались від сповіщень для всіх груп")
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
