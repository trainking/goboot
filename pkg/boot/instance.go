package boot

import (
	"strconv"

	"github.com/trainking/goboot/pkg/idgen"
	"github.com/trainking/goboot/pkg/log"

	"github.com/spf13/viper"
)

type (
	// Instance 定义一个实例
	Instance interface {
		// Start 开启方法
		Start() error

		// 初始化
		Init() error

		// Stop 结束方法
		Stop()
	}

	// BaseInstance 基础实例结构
	BaseInstance struct {
		Name      string // 实例名
		Config    *viper.Viper
		Addr      string // 监听地址
		IntanceID int64  // 实例ID
	}
)

// Init 初始化
func (b *BaseInstance) Init() error {
	// 初始化日志
	loggerConf := b.Config.GetStringMap("Logger")

	if len(loggerConf) > 0 {
		config := log.NewConfigByMap(loggerConf)
		config.ID = strconv.FormatInt(b.IntanceID, 10)
		log.InitLogger(config)
	} else {
		log.InitLogger(log.Config{
			Target: b.Name,
			ID:     strconv.FormatInt(b.IntanceID, 10),
		})
	}

	// 初始化ID生成器
	idgen.InitNode(b.IntanceID)
	return nil
}
