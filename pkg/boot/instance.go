package boot

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		Config    *viper.Viper
		Addr      string // 监听地址
		IntanceID int64  // 实例ID
	}
)

// Init 初始化
func (b *BaseInstance) Init() error {
	// 初始化日志
	loggerConf := b.Config.GetStringMap("Logger")
	config := log.NewConfigByMap(loggerConf)
	config.ID = strconv.FormatInt(b.IntanceID, 10)
	log.InitLogger(config)

	// 初始化ID生成器
	idgen.InitNode(b.IntanceID)
	return nil
}

// BootServe 启动服务
func BootServe(instance Instance) error {
	// 初始化
	if err := instance.Init(); err != nil {
		return err
	}

	// 2. 优雅退出
	go func() {
		exitC := make(chan os.Signal, 1)
		signal.Notify(exitC, syscall.SIGINT, syscall.SIGTERM)

		<-exitC
		instance.Stop()
		os.Exit(0)
	}()

	return instance.Start()
}
