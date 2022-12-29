package dynamodbx

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// GetClient 获取daynamoDB的客户端实例;
// 必须在环境变量中设置:
// - AWS_ACCESS_KEY_ID
// - AWS_SECRET_ACCESS_KEY
// - AWS_REGION
func GetClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	_dynamoDB := dynamodb.NewFromConfig(cfg)
	return _dynamoDB, nil
}
