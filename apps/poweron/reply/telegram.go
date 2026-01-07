package reply

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(cfg *config.Config) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramAPIToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Telegram{
		bot: bot,
	}, nil
}
