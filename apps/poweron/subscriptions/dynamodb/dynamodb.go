package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rvolykh/telegram-bot/apps/poweron/config"
	"github.com/rvolykh/telegram-bot/apps/poweron/subscriptions/models"
)

type DynamoDB struct {
	tableName *string
	client    *dynamodb.Client
}

func NewDynamoDB(cfg *config.Config) *DynamoDB {
	return &DynamoDB{
		tableName: aws.String(cfg.DynamoDBSubscriptionsTable),
		client:    dynamodb.NewFromConfig(cfg.AWSConfig),
	}
}

func (d *DynamoDB) InsertSubscription(ctx context.Context, chatID int64, group string) error {
	input := &dynamodb.PutItemInput{
		TableName: d.tableName,
		Item: map[string]types.AttributeValue{
			"ChatId": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", chatID),
			},
			"Groups": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: group},
				},
			},
		},
		ReturnValues: types.ReturnValueNone,
	}

	_, err := d.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to insert subscription: %w", err)
	}
	return nil
}

func (d *DynamoDB) GetSubscription(ctx context.Context, chatID int64) ([]string, error) {
	input := &dynamodb.GetItemInput{
		TableName: d.tableName,
		Key: map[string]types.AttributeValue{
			"ChatId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", chatID)},
		},
	}
	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if result.Item == nil {
		return []string{}, nil
	}

	var subscription models.Subscription
	err = attributevalue.UnmarshalMap(result.Item, &subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return subscription.Groups, nil
}

func (d *DynamoDB) UpdateSubscription(ctx context.Context, chatID int64, group string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: d.tableName,
		Key: map[string]types.AttributeValue{
			"ChatId": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", chatID),
			},
		},
		UpdateExpression: aws.String("SET #gs = list_append(if_not_exists(#gs, :empty_list), :new_group)"),
		ExpressionAttributeNames: map[string]string{
			"#gs": "Groups",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":new_group": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: group},
				},
			},
			":empty_list": &types.AttributeValueMemberL{},
		},
		ReturnValues: types.ReturnValueNone,
	}

	_, err := d.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

func (d *DynamoDB) DeleteSubscription(ctx context.Context, chatID int64) error {
	input := &dynamodb.DeleteItemInput{
		TableName: d.tableName,
		Key: map[string]types.AttributeValue{
			"ChatId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", chatID)},
		},
		ReturnValues: types.ReturnValueNone,
	}

	_, err := d.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (d *DynamoDB) ListSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	input := &dynamodb.ScanInput{
		TableName: d.tableName,
	}
	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	var subscriptions []models.Subscription
	err = attributevalue.UnmarshalListOfMaps(result.Items, &subscriptions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscriptions: %w", err)
	}

	return subscriptions, nil
}
