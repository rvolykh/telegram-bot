package reply

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) GetPinnedMessage(ctx context.Context, chatID int64) (Message, error) {
	input := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}

	output, err := t.bot.GetChat(input)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get pinned message: %w", err)
	}

	if output.PinnedMessage == nil {
		return Message{}, nil
	}

	return Message{
		ID:   output.PinnedMessage.MessageID,
		Text: output.PinnedMessage.Text,
	}, nil
}
