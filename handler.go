package slogawslambda

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LambdaHandler struct {
	slog.JSONHandler
}

func NewAWSLambdaHandler(ctx context.Context, opts *slog.HandlerOptions) slog.Handler {
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
	return slog.NewJSONHandler(os.Stdout, opts).
		WithAttrs([]slog.Attr{slog.String("function_arn", arn)}).
		WithAttrs([]slog.Attr{slog.String("request_id", requestID)})
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
