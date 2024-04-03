package slogawslambda

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name             string
		envVariableValue string
		want             slog.Leveler
	}{
		{
			name:             "Should return Info given not set",
			envVariableValue: "",
			want:             slog.LevelInfo,
		},
		{
			name:             "Should return Info given invalid value",
			envVariableValue: "not-a-valid-value",
			want:             slog.LevelInfo,
		},
		{
			name:             "Should return Debug given lowercase debug",
			envVariableValue: "debug",
			want:             slog.LevelDebug,
		},
		{
			name:             "Should return Warn given uppercase warn",
			envVariableValue: "WARN",
			want:             slog.LevelWarn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVariableValue != "" {
				os.Setenv("LOG_LEVEL", tt.envVariableValue)
			}

			got := getLogLevel()

			os.Unsetenv("LOG_LEVEL")

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestUsingNonAwsContextShouldNotError(t *testing.T) {
	ctx := context.Background()

	h := NewAWSLambdaHandler(ctx, nil)

	assert.NotNil(t, h)
}
