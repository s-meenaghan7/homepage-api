package visitor

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type VisitorRepository interface {
	IncrementCount(ctx context.Context, pageID string) (int, error)
	RetrieveCount(ctx context.Context, pageID string) (int, error)
}

type VisitorCountItem struct {
	Count  int    `json:"count"`
	PageID string `json:"page_id"`
}

type DynamoVisitorRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoRepository(client *dynamodb.Client, tableName string) *DynamoVisitorRepository {
	return &DynamoVisitorRepository{client: client, tableName: tableName}
}

func (r *DynamoVisitorRepository) IncrementCount(
	ctx context.Context,
	pageID string,
) (int, error) {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"page_id": &types.AttributeValueMemberS{
				Value: pageID,
			},
		},
		UpdateExpression: aws.String(
			"SET #count = if_not_exists(#count, :zero) + :inc",
		),
		ExpressionAttributeNames: map[string]string{
			"#count": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc":  &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
		ReturnValues: types.ReturnValueAllNew,
	}

	out, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return -1, err
	}

	var item VisitorCountItem
	err = attributevalue.UnmarshalMap(out.Attributes, &item)
	if err != nil {
		return -1, err
	}

	return item.Count, nil
}

func (r *DynamoVisitorRepository) RetrieveCount(
	ctx context.Context,
	pageID string,
) (int, error) {
	output, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"page_id": &types.AttributeValueMemberS{
				Value: pageID,
			},
		},
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return -1, err
	}

	if output.Item == nil {
		return 0, nil
	}

	var item VisitorCountItem
	err = attributevalue.UnmarshalMap(output.Item, &item)
	if err != nil {
		return -1, err
	}

	return item.Count, nil
}
