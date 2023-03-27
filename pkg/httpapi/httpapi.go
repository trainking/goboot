package httpapi

import (
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// App is the application entrypoint.
	App struct {
		boot.BaseInstance

		modeules []Module

		e         *echo.Echo
		validator *StructValidator // 验证器
	}

	// Router 路由
	Router struct {
		Method      string           // 方法
		Path        string           // 路径
		Name        string           // 名称
		Handle      echo.HandlerFunc // 处理函数
		Middlewares []Middleware     // 中间件函数
	}

	// Group 路由的分组
	Group struct {
		Path        string       // 路径
		Middlewares []Middleware // 中间件函数
		Routers     []Router     // 路由
	}

	// Middleware 中间件抽象
	Middleware interface {
		MiddlewareFunc() echo.MiddlewareFunc
	}

	// StructValidator 结构体验证器
	StructValidator struct {
		validator *validator.Validate
	}

	// Module 按模块组合
	Module interface {
		// 初始化模块
		Init(app *App)
		// 模块的分组路由
		Group() Group
	}
)

// New creates a new application.
func New(configPath string, addr string, instancdID int64) *App {
	// 加载配置
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	app := new(App)
	app.e = echo.New()
	app.validator = NewStructValidator()
	app.Config = v
	app.Addr = addr
	app.IntanceID = instancdID
	return app
}

// Validate 实现echo.Validator
func (s *StructValidator) Validate(i interface{}) error {
	return s.validator.Struct(i)
}

// AddValidator 增加自定义验证器
func (s *StructValidator) AddValidator(tag string, v validator.Func) error {
	return s.validator.RegisterValidation(tag, v)
}

// transEchoValidator 转换为echo接口
func (s *StructValidator) transEchoValidator() echo.Validator {
	return s
}

// NewStructValidator 构建新结构体
func NewStructValidator() *StructValidator {
	return &StructValidator{validator: validator.New()}
}

// AddRouter adds a router to the application.
func (a *App) AddRouter(r Router) {
	a.e.Add(r.Method, r.Path, r.Handle, r.GetMiddlewares()...)
}

// AddGroup adds a group to the application.
func (a *App) AddGroup(g Group) {
	_g := a.e.Group(g.Path, g.GetMiddlewares()...)
	for _, r := range g.Routers {
		_g.Add(r.Method, r.Path, r.Handle, r.GetMiddlewares()...)
	}
}

// AddModule adds a module to the application.
func (a *App) AddModule(m Module) {
	a.modeules = append(a.modeules, m)
	a.AddGroup(m.Group())
}

// Use adds a middleware to the application.
func (a *App) Use(m ...echo.MiddlewareFunc) {
	a.e.Use(m...)
}

// AddValidator 增加自定义验证器
func (a *App) AddValidator(key string, v validator.Func) {
	a.validator.AddValidator(key, v)
}

// start starts the application.
func (a *App) Start() error {
	// 增加验证器
	a.e.Validator = a.validator.transEchoValidator()

	// 全局中间件
	a.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           `{"level": "ACCESS", "ts":"${time_rfc3339}", "id": "${id}", "remote_ip":"${remote_ip}", "host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}", "status":${status},"latency_human":"${latency_human}"}` + "\n",
		CustomTimeFormat: "2006-01-02T15:04:05.000Z",
		Output:           log.GetWriter(),
	}))
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.CORS())
	a.e.Use(middleware.Gzip())

	return a.e.Start(a.Addr)
}

// Init 初始化阶段
func (a *App) Init() error {
	a.BaseInstance.Init()

	// Init各个模块
	for _, m := range a.modeules {
		m.Init(a)
	}

	return nil
}

// Stop 停止服务
func (a *App) Stop() {
	a.e.Close()
}

// GetMiddlewares 获取所有注册中间件
func (r *Router) GetMiddlewares() []echo.MiddlewareFunc {
	var middlewares []echo.MiddlewareFunc
	for _, m := range r.Middlewares {
		middlewares = append(middlewares, m.MiddlewareFunc())
	}
	return middlewares
}

// GetMiddlewares 获取所有注册中间件
func (g *Group) GetMiddlewares() []echo.MiddlewareFunc {
	var middlewares []echo.MiddlewareFunc
	for _, m := range g.Middlewares {
		middlewares = append(middlewares, m.MiddlewareFunc())
	}
	return middlewares
}
