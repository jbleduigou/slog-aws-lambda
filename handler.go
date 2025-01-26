package slogawslambda

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LambdaHandler struct {
	slog.JSONHandler
}

// NewAWSLambdaHandler creates a new AWS Lambda handler
// ctx is the context of the lambda function
// opts are the options for the handler
// envVars is an optional list of environment variables to add to the handler
func NewAWSLambdaHandler(ctx context.Context, opts *slog.HandlerOptions, envVars ...string) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	// Set the level based on the environment variable `LOG_LEVEL`
	if opts.Level == nil {
		opts.Level = getLogLevel()
	}

	// Retrieve AWS Request ID and lambda function arn
	lc, found := lambdacontext.FromContext(ctx)

	if !found {
		return slog.NewJSONHandler(os.Stdout, opts)
	}

	requestID := lc.AwsRequestID
	arn := lc.InvokedFunctionArn

	// Create the Handler using the attributes from lambda context
	h := slog.NewJSONHandler(os.Stdout, opts).
		WithAttrs([]slog.Attr{slog.String("function_arn", arn)}).
		WithAttrs([]slog.Attr{slog.String("request_id", requestID)})

	if len(envVars) == 0 {
		return h
	}

	// Add the environment variables to the handler as attributes
	for _, name := range envVars {
		val, found := os.LookupEnv(name)
		if found && len(name) > 0 {
			h = h.WithAttrs([]slog.Attr{slog.String(strings.ToLower(name), val)})
		}
	}

	return h
}

func getLogLevel() slog.Leveler {
	str, found := os.LookupEnv("LOG_LEVEL")

	// If no value is set, use Info as default Level
	if !found {
		return slog.LevelInfo
	}

	var l slog.Level
	err := l.UnmarshalText([]byte(str))

	// If invalid value is set, use Info as default Level
	if err != nil {
		return slog.LevelInfo
	}

	return l
}
