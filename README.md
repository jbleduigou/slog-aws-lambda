# slog: AWS Lambda handler

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/jbleduigou/slog-aws-lambda?status.svg)](https://pkg.go.dev/github.com/jbleduigou/slog-aws-lambda)
![Build Status](https://github.com/jbleduigou/slog-aws-lambda/actions/workflows/go.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/jbleduigou/slog-aws-lambda)](https://goreportcard.com/report/github.com/jbleduigou/slog-aws-lambda)
[![Contributors](https://img.shields.io/github/contributors/jbleduigou/slog-aws-lambda)](https://github.com/jbleduigou/slog-aws-lambda/graphs/contributors)
[![License](https://img.shields.io/github/license/jbleduigou/slog-aws-lambda)](./LICENSE)

An [AWS Lambda Function](https://aws.amazon.com/lambda/) Handler for [slog](https://pkg.go.dev/log/slog) Go library.  
The idea is to provide a log handler that will use the attributes present in the lambda context.  
Also the handler is using Json format so that it can leverage the features offered by [AWS CloudWatch](https://aws.amazon.com/cloudwatch/).

## 🚀 Install

```sh
go get github.com/jbleduigou/slog-aws-lambda
```

**Compatibility**: go >= 1.21

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/jbleduigou/slog-aws-lambda](https://pkg.go.dev/github.com/jbleduigou/slog-aws-lambda)

The handler will log the following attributes by default:
* `request-id`: the unique request ID for the lambda function invocation.
* `function-arn`: the Amazon Resource Name (ARN) that's used to invoke the function, indicates if the invoker specified a version number or alias.

### Using the handler as default handler

```go
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	slogawslambda "github.com/jbleduigou/slog-aws-lambda"
)

func handleS3Event(ctx context.Context, s3Event events.S3Event) {
	slog.SetDefault(slog.New(slogawslambda.NewAWSLambdaHandler(ctx, nil)))
	//rest of the business logic
}
```

### Logging a message with a few attributes

Once the handler is defined as default, there is no additional configuration needed.  
We can start using the logger by directly using the `slog` package.

```go
import "log/slog"

slog.Info("Successfully downloaded configuration file from S3", "bucket", bucket, "object-key", objectKey)
```

### Changing the key for request ID or function ARN

The default case for the attributes is dash-case, i.e. we are using `request-id` and `function-arn`.  
If you want to switch to another case you can use the `ReplaceAttr` mechanism provided by the slog package.  
This way allows not only to switch to camel case but also renaming altogether the attribute.

```go
slog.SetDefault(slog.New(slogawslambda.NewAWSLambdaHandler(ctx, &slog.HandlerOptions{
    ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		// change request ID to camel case
        if a.Key == "request-id" {
            a.Key = "requestID"
            return a
        }
        // change the name of the attribute for the invoked function ARN
        if a.Key == "function-arn" {
            a.Key = "invokedFunctionARN"
            return a
        }
        return a
    },
})))
```

### Stop logging the function ARN

You might also want to stop logging the function ARN.  
For instance if you are not using it and want to save cost on your CloudWatch logs.

```go
slog.SetDefault(slog.New(slogawslambda.NewAWSLambdaHandler(ctx, &slog.HandlerOptions{
    ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
        if a.Key == "request-id" {
            // don't log function ARN
			return slog.Attr{}
        }
        return a
    },
})))
```

## 📓 Why should we want to log request ID and function ARN

Embracing serverless comes with scalability and resilience benefits, but logging across a distributed system can be challenging.  
AWS generates a unique request ID for each Lambda execution, serving as a critical identifier for tracing requests and facilitating troubleshooting.  

Implementing this logging strategy involves retrieving the request ID from lambdacontext, which is exactly what this handler does.  

Once this is configured it becomes easier to retrieve all the logs for a specific execution, whether you use the AWS Console or a third party tool like DataDog or SumoLogic.

## 🤝 Contributing

Feel free to contribute to this project, either my opening issues or submitting pull requests.  
Don't hesitate to contact me, by sending me a PM on [LinkedIn](www.linkedin.com/in/jbleduigou).

## 👤 Contributors

The only contributor so far is me, Jean-Baptiste Le Duigou.
Feel free to check my blog to about my other projects: http://www.jbleduigou.com

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/jbleduigou?style=for-the-badge)](https://github.com/sponsors/jbleduigou)

## 📝 License

This project is [MIT](./LICENSE) licensed.
