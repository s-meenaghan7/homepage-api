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
func InitDynamoClient(ctx context.Context, endpoint string) *dynamodb.Client {
	containerEndpoint := "http://" + endpoint
	fmt.Printf("Initializing DynamoDB client pointing to local instance at [%s]...\n", containerEndpoint)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("dummy-region"),
		config.WithBaseEndpoint(containerEndpoint),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DynamoDB client initialized successfully!\n")

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Region = "dummy-region"
		o.BaseEndpoint = &containerEndpoint
	})
}

// Setup the DynamoDB visitors table used for testing. Automatically calls a cleanup function to delete the table after the test completes.
func SetupVisitorTable(client *dynamodb.Client, t *testing.T) {
	t.Helper()

	tableName := "test-visitors"
	fmt.Printf("Creating DynamoDB table [%s] for testing...\n", tableName)

	_, err := client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
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
	waiter := dynamodb.NewTableExistsWaiter(client)
	err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 30*time.Second)
	if err != nil {
		panic(fmt.Sprintf("failed to wait for table to be created: %v", err))
	}
	fmt.Printf("Table [%s] created successfully for test [%s].\n", tableName, t.Name())

	t.Cleanup(func() {
		_, err := client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Table [%s] cleaned up successfully.\n", tableName)
	})
}
