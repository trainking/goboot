package httpapi

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/trainking/goboot/pkg/utils"
)

// Hhttp2App http2协议的实现
type Http2App struct {
	App

	CertPemPath string // cert.pem所在文件的路径
	KeyPemPath  string // key.pem所在文件的路径
}

// NewHttp2App 创建一个http2协议的实现
func NewHttp2App(configPath string, addr string, instancdID int64) *Http2App {
	// 加载配置
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	app := new(Http2App)
	app.e = echo.New()
	app.validator = NewStructValidator()
	app.Config = v
	app.Addr = addr
	app.IntanceID = instancdID
	return app
}

// Init 初始化
func (a *Http2App) Init() error {
	a.App.Init()

	a.CertPemPath = a.Config.GetString("Http2Conf.CertPem")
	a.KeyPemPath = a.Config.GetString("Http2Conf.KeyPem")

	return nil
}

// Start 开始
func (a *Http2App) Start() error {
	// 增加验证器
	a.e.Validator = a.validator.transEchoValidator()

	// 全局中间件
	a.e.Use(middleware.Logger())
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.CORS())
	a.e.Use(middleware.Gzip())

	if a.CertPemPath != "" && a.KeyPemPath != "" {
		return a.e.StartTLS(a.Addr, a.CertPemPath, a.KeyPemPath)
	}

	return a.e.StartAutoTLS(a.Addr)
}
