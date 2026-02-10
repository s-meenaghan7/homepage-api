package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Start a local DynamoDB instance using Docker. Returns the container endpoint to communicate with DynamoDB and a cleanup function to stop the container after tests complete.
func RunDockerDynamoDB(ctx context.Context) (string, func()) {
	image := "amazon/dynamodb-local"
	port := "8000/tcp"
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(30 * time.Second),
	}
	fmt.Printf("Initializing test container w/ image [%q]...\n", image)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
	if err != nil {
		panic(err)
	}
	
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test container initialized successfully, endpoint: [%s]\n", endpoint)

	return endpoint, func() {
		fmt.Printf("Terminating test container w/ image [%q]...\n", image)
		err := container.Terminate(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println("Test container terminated successfully")
	}
}
