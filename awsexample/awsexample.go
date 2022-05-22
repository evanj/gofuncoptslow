package awsexample

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CallDynamoDB demonstrates passing functional options to dynamodb.Client.GetItem.
func CallDynamoDB(ctx context.Context, client *dynamodb.Client) (string, error) {
	keyValue, err := attributevalue.Marshal("item_primary_key")
	if err != nil {
		return "", err
	}
	const tableName = "example_table"
	tableNameValue := tableName
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"primary_key": keyValue,
		},

		TableName: &tableNameValue,
	}

	demoEndpointResolver := dynamodb.EndpointResolverFromURL("https://example.com")
	// inlining call to dynamodb.WithEndpointResolver
	// WithEndpointResolver: dynamodb@v1.15.3/api_client.go:176:14: o does not escape
	// func literal does not escape
	endpointResolverOption := dynamodb.WithEndpointResolver(demoEndpointResolver)
	// GetItem: api_op_GetItem.go:23:69: optFns does not escape
	// ... argument does not escape
	out, err := client.GetItem(ctx, input, endpointResolverOption)
	if err != nil {
		return "", err
	}

	// do something with the result: return the first key (if any)
	for key := range out.Item {
		return key, nil
	}
	return "", nil
}
