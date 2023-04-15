package boot

import (
	"os"
	"os/signal"
	"syscall"
)

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

	if err := instance.Start(); err != nil {
		instance.Stop()
		return err
	}

	return nil
}
