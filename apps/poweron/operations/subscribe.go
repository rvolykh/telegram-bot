package operations

import (
	"context"
	"fmt"
	"slices"

	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions"
)

func Subscribe(ctx context.Context, cfg *config.Config, chatID int64, group string) error {
	t, err := reply.NewTelegram(cfg)
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	s, err := subscriptions.NewSubscriptionsDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions repository: %w", err)
	}

	groups, err := s.GetSubscription(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if len(groups) == 0 {
		if err = s.InsertSubscription(ctx, chatID, group); err != nil {
			return fmt.Errorf("failed to insert subscription: %w", err)
		}
	} else {
		if slices.Contains(groups, group) {
			_, err = t.SendMessage(ctx, chatID, "Ви вже підписані на цю групу")
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		}

		if err = s.UpdateSubscription(ctx, chatID, group); err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
	}

	_, err = t.SendMessage(ctx, chatID, "Ви підписалися на сповіщення для групи "+group)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
