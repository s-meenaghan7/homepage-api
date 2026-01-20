package visitor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type mockService struct {
	getVisitsFn   func(ctx context.Context, pageID string) (int, error)
	handleVisitFn func(ctx context.Context, pageID string) (int, error)
}

func (m *mockService) GetVisits(ctx context.Context, pageID string) (int, error) {
	return m.getVisitsFn(ctx, pageID)
}

func (m *mockService) HandleVisit(ctx context.Context, pageID string) (int, error) {
	return m.handleVisitFn(ctx, pageID)
}

func makeHandlerRequest(
	method string,
	pathParams map[string]string,
) events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		PathParameters: pathParams,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: method,
			},
		},
	}
}

func TestHandler_GET_Success(t *testing.T) {
	req := makeHandlerRequest("GET", map[string]string{"page_id": "/"})

	expected := 5
	svc := &mockService{
		getVisitsFn: func(ctx context.Context, pageID string) (int, error) {
			return expected, nil
		},
	}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d. %s", response.StatusCode, response.Body)
	}

	var b VisitorCountResponseBody
	json.Unmarshal([]byte(response.Body), &b)

	if b.Count != expected {
		t.Fatalf("expected count to equal [%d], but got [%d]", expected, b.Count)
	}
}

func TestHandler_GET_Error(t *testing.T) {
	req := makeHandlerRequest("GET", map[string]string{"page_id": "/ "})
	svc := &mockService{
		getVisitsFn: func(ctx context.Context, pageID string) (int, error) {
			return -1, errors.New("Error in GetVisits")
		},
	}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("handler should not return errors to lambda function")
	}
	if response.StatusCode != 400 {
		t.Fatalf("expected response.StatusCode to be [%d], but got [%d]", 400, response.StatusCode)
	}

	var body map[string]string
	json.Unmarshal([]byte(response.Body), &body)

	_, ok := body["error"]
	if !ok {
		t.Fatalf("expected key error is missing from response body")
	}
}

func TestHandler_POST_Success(t *testing.T) {
	req := makeHandlerRequest("POST", map[string]string{"page_id": "home"})

	expected := 2
	svc := &mockService{
		handleVisitFn: func(ctx context.Context, pageID string) (int, error) {
			return expected, nil
		},
	}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d. %s", response.StatusCode, response.Body)
	}

	var b VisitorCountResponseBody
	json.Unmarshal([]byte(response.Body), &b)

	if b.Count != expected {
		t.Fatalf("expected count to equal [%d], but got [%d]", expected, b.Count)
	}
}

func TestHandler_POST_Error(t *testing.T) {
	req := makeHandlerRequest("POST", map[string]string{"page_id": "home"})

	svc := &mockService{
		handleVisitFn: func(ctx context.Context, pageID string) (int, error) {
			return -1, errors.New("Error in HandleVisit")
		},
	}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("handler should not return errors to lambda function")
	}
	if response.StatusCode != 400 {
		t.Fatalf("expected response.StatusCode to be [%d], but got [%d]", 400, response.StatusCode)
	}

	var body map[string]string
	json.Unmarshal([]byte(response.Body), &body)

	_, ok := body["error"]
	if !ok {
		t.Fatalf("expected key error is missing from response body")
	}
}

func TestMissingPageID(t *testing.T) {
	req := makeHandlerRequest("GET", nil)

	svc := &mockService{}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("handler should not return errors to lambda function")
	}
	if response.StatusCode != 400 {
		t.Fatalf("expected response.StatusCode to be [%d], but got [%d]", 400, response.StatusCode)
	}

	var body map[string]string
	json.Unmarshal([]byte(response.Body), &body)

	_, ok := body["error"]
	if !ok {
		t.Fatalf("expected key error is missing from response body")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	req := makeHandlerRequest("DELETE", map[string]string{"page_id": "home"})

	svc := &mockService{}
	handler := Handler(svc)

	response, err := handler(context.Background(), req)

	if err != nil {
		t.Fatalf("handler should not return errors to lambda function")
	}
	if response.StatusCode != 405 {
		t.Fatalf("expected response.StatusCode to be [%d], but got [%d]", 405, response.StatusCode)
	}

	var body map[string]string
	json.Unmarshal([]byte(response.Body), &body)

	_, ok := body["error"]
	if !ok {
		t.Fatalf("expected key error is missing from response body")
	}
}
