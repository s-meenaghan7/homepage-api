package responses

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func Error(statusCode int, message string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"error": "%s"}`, message),
	}
}

func Success(statusCode int, body interface{}) events.APIGatewayV2HTTPResponse {
	jsonBody, _ := json.Marshal(body)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonBody),
	}
}
