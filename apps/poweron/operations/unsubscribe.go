package operations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
)

func Unsubscribe(ctx context.Context, cfg config.Config, t *reply.Telegram, ssmClient *ssm.Client, chatID int64) error {
	// TODO: Implement unsubscribe logic
	return t.SendMessage(ctx, chatID, "Відписалися від усіх груп")
}
