package visitor

import (
	"context"
	"homepage-api/internal/responses"

	"github.com/aws/aws-lambda-go/events"
)

type VisitorCountResponseBody struct {
	Count int `json:"count"`
}

func Handler(svc VisitorService) func(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	return func(
		ctx context.Context,
		req events.APIGatewayV2HTTPRequest,
	) (events.APIGatewayV2HTTPResponse, error) {
		pageID, ok := req.PathParameters["page_id"]
		if !ok {
			return responses.Error(400, "Bad Request: missing page_id path parameter"), nil
		}

		if req.RequestContext.HTTP.Method == "GET" {
			count, err := svc.GetVisits(ctx, pageID)
			if err != nil {
				return responses.Error(400, err.Error()), nil
			}
			return responses.Success(200, VisitorCountResponseBody{Count: count}), nil
		}

		if req.RequestContext.HTTP.Method == "POST" {
			count, err := svc.HandleVisit(ctx, pageID)
			if err != nil {
				return responses.Error(400, err.Error()), nil
			}
			return responses.Success(200, VisitorCountResponseBody{Count: count}), nil
		}

		// default, method not allowed
		return responses.Error(405, "Method Not Allowed"), nil
	}
}
