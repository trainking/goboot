package mongodbx

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewCollection 创建一个集合的的操作
func NewCollection(url string, database string, collection string) *mongo.Collection {
	// 设定10秒建立连接不成功，则超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	return client.Database(database).Collection(collection)
}
