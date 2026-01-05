package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
	"github.com/rvolykh/telegram-bot/apps/poweron/scrap"
)

func Today(ctx context.Context, cfg config.Config, t *reply.Telegram, ssmClient *ssm.Client, chatID int64) error {
	poweron := scrap.NewPoweron(ssmClient, cfg.SSMParamPowerScheduleCache, 1*time.Hour)

	powerSchedule, err := poweron.GetPowerSchedule(ctx)
	if err != nil {
		return fmt.Errorf("failed to get power schedule: %w", err)
	}

	return t.SendMessage(ctx, chatID, powerSchedule)
}
