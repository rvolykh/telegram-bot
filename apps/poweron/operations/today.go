package operations

import (
	"context"
	"fmt"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
	"github.com/rvolykh/telegram-bot/apps/poweron/scrap"
)

func Today(ctx context.Context, cfg *config.Config, chatID int64) error {
	poweron := scrap.NewPoweron(cfg)

	powerSchedule, err := poweron.GetPowerSchedule(ctx)
	if err != nil {
		return fmt.Errorf("failed to get power schedule: %w", err)
	}

	t, err := reply.NewTelegram(cfg)
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	_, err = t.SendMessage(ctx, chatID, powerSchedule.Today)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
