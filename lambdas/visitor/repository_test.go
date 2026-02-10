//go:build integration

package visitor

import (
	"context"
	"os"
	"testing"

	"homepage-api/internal/testutil"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const tableName = "test-visitors"

var testClient *dynamodb.Client

// TestMain sets up the testing environment by starting a local DynamoDB instance using Docker and initializing the local DynamoDB client.
func TestMain(m *testing.M) {
	ctx := context.Background()
	endpoint, cleanup := testutil.RunDockerDynamoDB(ctx)
	testClient = testutil.InitDynamoClient(ctx, endpoint)

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestIncrementCount(t *testing.T) {
	ctx := context.Background()
	testutil.SetupVisitorTable(testClient, t)

	pageId := "test-page"
	repo := NewDynamoRepository(testClient, tableName)
	count, err := repo.IncrementCount(ctx, pageId)

	if err != nil {
		t.Fatalf("IncrementCount failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count of 1, got %d", count)
	}
}

func TestIncrementCountMultiple(t *testing.T) {
	ctx := context.Background()
	testutil.SetupVisitorTable(testClient, t)

	repo := NewDynamoRepository(testClient, tableName)

	iterations := 5
	pageID := "test-page-multiple"
	for i := range iterations {
		count, err := repo.IncrementCount(ctx, pageID)
		if err != nil {
			t.Fatalf("IncrementCount failed: %v", err)
		}
		if count != i+1 {
			t.Errorf("IncrementCount expected %d got %d", i+1, count)
		}
	}
}

func TestRetrieveCountEmpty(t *testing.T) {
	ctx := context.Background()
	testutil.SetupVisitorTable(testClient, t)

	repo := NewDynamoRepository(testClient, tableName)
	count, err := repo.RetrieveCount(ctx, "nonexistent-page-12345")

	if err != nil {
		t.Fatalf("RetrieveCount failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0 for nonexistent page, got %d", count)
	}
}
