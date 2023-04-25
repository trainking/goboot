package boot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/trainking/goboot/pkg/log"
)

// BootServe 启动服务
func BootServe(instance Instance) error {
	var err error

	// 初始化
	if err = instance.Init(); err != nil {
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

	defer func() {
		if err != nil {
			log.Errorf("BootServe error: %s", err)
		}

		e := recover()
		if e != nil {
			log.Errorf("BootServe unknow error: %v", e)
		}
	}()

	if err = instance.Start(); err != nil {
		instance.Stop()
		return err
	}

	return nil
}
