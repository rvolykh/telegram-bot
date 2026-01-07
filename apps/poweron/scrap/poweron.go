package scrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
)

type Poweron struct {
	ssmClient *ssm.Client
	cacheKey  string
	cacheTTL  time.Duration
}

func NewPoweron(cfg *config.Config) *Poweron {
	return &Poweron{
		ssmClient: ssm.NewFromConfig(cfg.AWSConfig),
		cacheKey:  cfg.SSMParamCache,
		cacheTTL:  cfg.PowerScheduleCacheTTL,
	}
}

func (p *Poweron) GetPowerSchedule(ctx context.Context) (Schedule, error) {
	cacheEntry, err := p.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(p.cacheKey),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to get power schedule from cache: %v", err)
	} else {
		isCacheValid := cacheEntry.Parameter.Value != nil &&
			*cacheEntry.Parameter.Value != "none" &&
			cacheEntry.Parameter.LastModifiedDate.After(time.Now().Add(-p.cacheTTL))

		if isCacheValid {
			log.Printf("Power schedule is still valid, returning cached value")

			var powerSchedule Schedule
			if err := json.Unmarshal([]byte(*cacheEntry.Parameter.Value), &powerSchedule); err != nil {
				return Schedule{}, fmt.Errorf("failed to unmarshal power schedule: %w", err)
			}
			return powerSchedule, nil
		}
		log.Printf("Power schedule is outdated, refreshing cache")
	}

	powerSchedule, err := getPowerSchedule()
	if err != nil {
		return Schedule{}, fmt.Errorf("failed to get power schedule: %w", err)
	}

	powerScheduleJSON, err := json.Marshal(powerSchedule)
	if err != nil {
		return Schedule{}, fmt.Errorf("failed to marshal power schedule: %w", err)
	}

	_, err = p.ssmClient.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(p.cacheKey),
		Value:     aws.String(string(powerScheduleJSON)),
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to put power schedule to cache: %v", err)
	} else {
		log.Printf("Updated power schedule in cache")
	}

	return powerSchedule, nil
}

func getPowerSchedule() (Schedule, error) {
	var schedule Schedule

	resp, err := http.Get("https://api.loe.lviv.ua/api/menus?page=1&type=photo-grafic")
	if err != nil {
		return schedule, fmt.Errorf("failed to get power on: %w", err)
	}
	defer resp.Body.Close()

	var powerOnResponse PowerOnResponse
	err = json.NewDecoder(resp.Body).Decode(&powerOnResponse)
	if err != nil {
		return schedule, fmt.Errorf("failed to decode power on response: %w", err)
	}

	if len(powerOnResponse.Member) == 0 {
		return schedule, fmt.Errorf("no power on found")
	}
	if len(powerOnResponse.Member[0].MenuItems) == 0 {
		return schedule, fmt.Errorf("no menu items found")
	}

	for _, menuItem := range powerOnResponse.Member[0].MenuItems {
		switch menuItem.Name {
		case "Today":
			schedule.Today = cleanHTML(menuItem.RawMobileHTML)
		case "Tomorrow":
			schedule.Tomorrow = cleanHTML(menuItem.RawMobileHTML)
		}
	}

	if schedule.Today == "" {
		schedule.Today = "Немає запланованих відключень електроенергії"
	}
	if schedule.Tomorrow == "" {
		schedule.Tomorrow = "Немає запланованих відключень електроенергії"
	}

	return schedule, nil
}
