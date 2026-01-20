package main

import (
	"context"
	"fmt"
	"os"

	"homepage-api/lambdas/visitor"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}

	repo := visitor.NewDynamoRepository(
		dynamodb.NewFromConfig(cfg),
		os.Getenv("TABLE_NAME"),
	)
	svc := visitor.NewService(repo)
	handler := visitor.Handler(svc)

	lambda.Start(handler)
}
