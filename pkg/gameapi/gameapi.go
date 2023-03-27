package gameapi

import (
	"net"
	"sync"

	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/utils"
)

type (
	// App 游戏服务器实现
	App struct {
		boot.BaseInstance

		listener  net.Listener    // 网络监听
		exitChan  chan struct{}   // 退出信号
		waitGroup *sync.WaitGroup // 等待携程控制
		closeOnce sync.Once       // 保证关闭只执行一次
	}
)

// New 创建一个游戏服务器接口实例
func New(configPath string, addr string, instancdID int64) *App {
	// 加载配置
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	app := new(App)
	app.Config = v
	app.Addr = addr
	app.IntanceID = instancdID
	return app
}

// Init 初始化服务
func (a *App) Init() error {
	if err := a.BaseInstance.Init(); err != nil {
		return err
	}

	return nil
}

// Start 启动服务
func (a *App) Start() error {
	a.waitGroup.Add(1)
	defer func() {
		a.waitGroup.Done()
	}()

	for {
		select {
		case <-a.exitChan:
			return nil
		default:
		}

		conn, err := a.listener.Accept()
		if err != nil {
			return err
		}

		a.waitGroup.Add(1)
		go func() {
			// TODO 处理连接数据
			a.waitGroup.Done()
		}()
	}
	return nil
}

// Stop 停止服务
func (a *App) Stop() error {
	// 关闭资源
	a.closeOnce.Do(func() {
		close(a.exitChan)
		a.listener.Close()
	})

	// 等待所有携程执行完
	a.waitGroup.Wait()
	return nil
}
