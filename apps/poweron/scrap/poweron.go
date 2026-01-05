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
)

type Poweron struct {
	ssmClient *ssm.Client
	cacheKey  string
	cacheTTL  time.Duration
}

func NewPoweron(ssmClient *ssm.Client, cacheKey string, cacheTTL time.Duration) *Poweron {
	return &Poweron{
		ssmClient: ssmClient,
		cacheKey:  cacheKey,
		cacheTTL:  cacheTTL,
	}
}

func (p *Poweron) GetPowerSchedule(ctx context.Context) (string, error) {
	powerSchedule, err := p.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(p.cacheKey),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to get power schedule from cache: %v", err)
	} else {
		isCacheValid := powerSchedule.Parameter.Value != nil &&
			*powerSchedule.Parameter.Value != "none" &&
			powerSchedule.Parameter.LastModifiedDate.After(time.Now().Add(-p.cacheTTL))

		if isCacheValid {
			log.Printf("Power schedule is still valid, returning cached value")
			return *powerSchedule.Parameter.Value, nil
		}
		log.Printf("Power schedule is outdated, refreshing cache")
	}

	powerScheduleText, err := getPowerSchedule()
	if err != nil {
		return "", fmt.Errorf("failed to get power schedule: %w", err)
	}

	_, err = p.ssmClient.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(p.cacheKey),
		Value:     aws.String(powerScheduleText),
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		log.Printf("Failed to put power schedule to cache: %v", err)
	} else {
		log.Printf("Updated power schedule in cache")
	}

	return powerScheduleText, nil
}

func getPowerSchedule() (string, error) {
	resp, err := http.Get("https://api.loe.lviv.ua/api/menus?page=1&type=photo-grafic")
	if err != nil {
		return "", fmt.Errorf("failed to get power on: %w", err)
	}
	defer resp.Body.Close()

	var powerOnResponse PowerOnResponse
	err = json.NewDecoder(resp.Body).Decode(&powerOnResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode power on response: %w", err)
	}

	if len(powerOnResponse.Member) == 0 {
		return "", fmt.Errorf("no power on found")
	}
	if len(powerOnResponse.Member[0].MenuItems) == 0 {
		return "", fmt.Errorf("no menu items found")
	}

	text := cleanHTML(powerOnResponse.Member[0].MenuItems[0].RawMobileHTML)
	if text == "" {
		text = "Немає запланованих відключень електроенергії"
	}

	return text, nil
}
