package mongodbx

import "go.mongodb.org/mongo-driver/bson"

// Map2BsonD 将map[string]interface{}结构转换未bson.D
func Map2BsonD(param map[string]interface{}) bson.D {
	var _cols = bson.D{}
	for k, v := range param {
		_cols = append(_cols, bson.E{Key: k, Value: v})
	}
	return _cols
}
