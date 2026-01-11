package config

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_WithAPITokenEnvVar(t *testing.T) {
	// Set test environment variables
	testRegion := "us-east-1"
	testAPIToken := "test-api-token-123"
	testSSMParamCache := "/test/cache/param"
	testDynamoDBTable := "test-subscriptions-table"

	os.Setenv("AWS_DEFAULT_REGION", testRegion)
	os.Setenv(EnvAPIToken, testAPIToken)
	os.Setenv(EnvSSMParamCache, testSSMParamCache)
	os.Setenv(EnvDynamoDBSubscriptionsTable, testDynamoDBTable)
	os.Unsetenv(EnvSSMParamAPIToken) // Ensure SSM param is not set

	ctx := context.Background()
	cfg, err := LoadConfig(ctx)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify API token
	assert.Equal(t, testAPIToken, cfg.TelegramAPIToken)

	// Verify SSM param cache
	assert.Equal(t, testSSMParamCache, cfg.SSMParamCache)

	// Verify DynamoDB table
	assert.Equal(t, testDynamoDBTable, cfg.DynamoDBSubscriptionsTable)

	// Verify default cache TTL
	assert.Equal(t, PowerScheduleCacheTTL, cfg.PowerScheduleCacheTTL)

	// Verify AWS config is set
	assert.Equal(t, testRegion, cfg.AWSConfig.Region)
}

func TestLoadConfig_WithEmptyEnvVars(t *testing.T) {
	// Unset all environment variables
	os.Unsetenv(EnvAPIToken)
	os.Unsetenv(EnvSSMParamAPIToken)
	os.Unsetenv(EnvSSMParamCache)
	os.Unsetenv(EnvDynamoDBSubscriptionsTable)

	ctx := context.Background()
	cfg, err := LoadConfig(ctx)

	// When API token is not set and SSM param is not set, it should fail
	// because it will try to get from SSM but the parameter name is empty
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestConfigConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{
			name:     "EnvAPIToken",
			constant: EnvAPIToken,
			want:     "TELEGRAM_APITOKEN",
		},
		{
			name:     "EnvSSMParamAPIToken",
			constant: EnvSSMParamAPIToken,
			want:     "SSM_PARAM_TELEGRAM_APITOKEN",
		},
		{
			name:     "EnvSSMParamCache",
			constant: EnvSSMParamCache,
			want:     "SSM_PARAM_CACHE",
		},
		{
			name:     "EnvDynamoDBSubscriptionsTable",
			constant: EnvDynamoDBSubscriptionsTable,
			want:     "DYNAMODB_TABLE_SUBSCRIPTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.constant)
		})
	}
}
