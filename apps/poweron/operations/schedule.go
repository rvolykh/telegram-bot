package operations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/reply"
	"github.com/rvolykh/telegram-bot/apps/poweron/scrap"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions"
)

func DeliverScheduleUpdates(ctx context.Context, cfg *config.Config) error {
	poweron := scrap.NewPoweron(cfg)

	powerSchedule, err := poweron.GetPowerSchedule(ctx)
	if err != nil {
		return fmt.Errorf("failed to get power schedule: %w", err)
	}

	s, err := subscriptions.NewSubscriptionsDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions repository: %w", err)
	}

	t, err := reply.NewTelegram(cfg)
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	subscribers, err := s.ListSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}

	for _, subscriber := range subscribers {
		message := prepareMessage(powerSchedule, subscriber.Groups)

		prev, err := t.GetPinnedMessage(ctx, subscriber.ChatID)
		if err != nil {
			log.Printf("Failed to get pinned message %d: %s", subscriber.ChatID, err)
		}

		// telegram trims last \n so should we before compare
		if prev == message[:len(message)-1] {
			log.Printf("No updates, skipping for chat %d", subscriber.ChatID)
			continue
		}

		messageID, err := t.SendMessage(ctx, subscriber.ChatID, message)
		if err != nil {
			log.Printf("Failed to send message to %d: %s", subscriber.ChatID, err)
			continue
		}
		log.Printf("Message sent to chat: %d", subscriber.ChatID)

		if err := t.PinMessage(ctx, subscriber.ChatID, messageID); err != nil {
			log.Printf("Failed to pin message in %d: %s", subscriber.ChatID, err)
			continue
		}
		log.Printf("Message %d is pinned in chat %d", messageID, subscriber.ChatID)
	}
	return nil
}

func prepareMessage(powerSchedule scrap.Schedule, groups []string) string {
	var message strings.Builder

	message.WriteString("Сьогодні:\n")
	filterPowerScheduleGroups(&message, powerSchedule.Today, groups)

	message.WriteString("\nЗавтра:\n")
	filterPowerScheduleGroups(&message, powerSchedule.Tomorrow, groups)

	return message.String()
}

func filterPowerScheduleGroups(b *strings.Builder, schedule string, groups []string) {
	for i, line := range strings.Split(schedule, "\n") {
		if i > 1 {
			var skip = true
			for _, group := range groups {
				if strings.Contains(line, group) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}

		b.WriteString(line + "\n")
	}
}
