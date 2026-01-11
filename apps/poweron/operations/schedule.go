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

const (
	tomorrow  = "Завтра"
	today     = "Сьогодні"
	updatedAt = "Інформація станом на"
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

		// one to one match
		if prev.Text == message {
			log.Printf("No updates, skipping for chat %d", subscriber.ChatID)
			continue
		}
		// updatedAt change only
		if isMessagesEqual(prev.Text, message) {
			if err := t.EditMessage(ctx, subscriber.ChatID, prev.ID, message); err != nil {
				log.Printf("Failed to edit message %d in chat %d: %s", prev.ID, subscriber.ChatID, err)
			}
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

	message.WriteString(today + ":\n")
	filterPowerScheduleGroups(&message, powerSchedule.Today, groups)

	message.WriteString("\n" + tomorrow + ":\n")
	filterPowerScheduleGroups(&message, powerSchedule.Tomorrow, groups)

	// remove last new line as telegram will strip it anyway and it can create issue in cmp
	result := message.String()
	return result[:len(result)-1]
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

func isMessagesEqual(src, dst string) bool {
	var (
		srcLines = strings.Split(src, "\n")
		dstLines = strings.Split(dst, "\n")
	)

	if len(srcLines) != len(dstLines) {
		return false
	}

	for i := range srcLines {
		if srcLines[i] == dstLines[i] {
			continue
		}

		if strings.HasPrefix(srcLines[i], updatedAt) && strings.HasPrefix(dstLines[i], updatedAt) {
			continue
		}

		return false
	}

	return true
}
