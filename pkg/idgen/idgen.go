package idgen

import (
	"github.com/bwmarrin/snowflake"
)

var _node *snowflake.Node

// InitNode 初始化生成ID生成器
func InitNode(instanceID int64) {
	var err error
	_node, err = snowflake.NewNode(instanceID)
	if err != nil {
		panic(err)
	}
}

// GenerateInt64 创建一个int64的ID
func GenerateInt64() int64 {
	return _node.Generate().Int64()
}

// GenerateString 创建一个string的ID
func GenerateString() string {
	return _node.Generate().String()
}
