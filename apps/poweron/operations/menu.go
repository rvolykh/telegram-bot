package operations

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
)

var (
	keyboardMainMenu = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron Сьогодні"),
			tgbotapi.NewKeyboardButton("/poweron Завтра"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron Підписатись"),
			tgbotapi.NewKeyboardButton("/poweron Відписатись"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron Закрити"),
		),
	)

	keyboardGroupsMenu = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 1.1"),
			tgbotapi.NewKeyboardButton("/poweron 1.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 2.1"),
			tgbotapi.NewKeyboardButton("/poweron 2.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 3.1"),
			tgbotapi.NewKeyboardButton("/poweron 3.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 4.1"),
			tgbotapi.NewKeyboardButton("/poweron 4.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 5.1"),
			tgbotapi.NewKeyboardButton("/poweron 5.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron 6.1"),
			tgbotapi.NewKeyboardButton("/poweron 6.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/poweron Закрити"),
		),
	)

	keyboardCloseMenu = tgbotapi.NewRemoveKeyboard(true)
)

func ShowMainMenu(ctx context.Context, t *reply.Telegram, chatID int64) error {
	err := t.SendMenu(ctx, chatID, "Оберіть операцію", keyboardMainMenu)
	if err != nil {
		return fmt.Errorf("failed to send main menu: %w", err)
	}
	return nil
}

func ShowSelectGroupMenu(ctx context.Context, t *reply.Telegram, chatID int64) error {
	err := t.SendMenu(ctx, chatID, "Оберіть групу", keyboardGroupsMenu)
	if err != nil {
		return fmt.Errorf("failed to send select group menu: %w", err)
	}
	return nil
}

func CloseMenu(ctx context.Context, t *reply.Telegram, chatID int64) error {
	err := t.SendMenu(ctx, chatID, "Закрито", keyboardCloseMenu)
	if err != nil {
		return fmt.Errorf("failed to send close menu: %w", err)
	}
	return nil
}
