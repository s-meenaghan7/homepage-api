package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Initialize the DynamoDB client to point to the local instance started by Docker.
func initDynamoClient(ctx context.Context) *dynamodb.Client {
	containerEndpoint := "http://localhost:8000"
	fmt.Printf("Initializing DynamoDB client pointing to local instance at [%s]...\n", containerEndpoint)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("dummy-region"),
		config.WithBaseEndpoint(containerEndpoint),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DynamoDB client initialized!\n")
	fmt.Printf("CFG BASE ENDPOINT: %q\n", *cfg.BaseEndpoint)
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Region = "dummy-region"
		o.BaseEndpoint = &containerEndpoint
	})
}

// Setup the DynamoDB visitors table used for testing. Automatically calls a cleanup function to delete the table after the test completes. Returns the initialized DynamoDB client used to create the table.
func SetupVisitorTable(t *testing.T) *dynamodb.Client {
	t.Helper()

	tableName := "test-visitors"
	testClient := initDynamoClient(context.Background())

	fmt.Printf("Creating DynamoDB table [%s] for testing...\n", tableName)
	output, err := testClient.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("page_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("page_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		panic(err)
	}
	// Wait for table to be created before proceeding.
	waiter := dynamodb.NewTableExistsWaiter(testClient)
	err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 30*time.Second)
	if err != nil {
		panic(fmt.Sprintf("failed to wait for table to be created: %v", err))
	}
	fmt.Printf("Table [%s] created successfully: %v\n", tableName, output)

	t.Cleanup(func() {
		_, err := testClient.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			panic(err)
		}
	})

	return testClient
}
