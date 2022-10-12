package httpapi

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trainking/goboot/pkg/jwt"
)

// MiddlewareBearerJwt 使用Bearer Token的一个中间件
type MiddlewareBearerJwt struct {
	jt     *jwt.JwtToken // jwt实现
	expire time.Duration // token的过期时间
}

// NewMiddlewareBearerJwt 使用Bearer Token的中间件生成器
func NewMiddlewareBearerJwt(jwtSecert string, expire time.Duration) Middleware {
	return &MiddlewareBearerJwt{
		jt:     jwt.NewJwtToken(jwtSecert),
		expire: expire,
	}
}

// MiddlewareFunc 对Middleware的实现
func (m *MiddlewareBearerJwt) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			if token == "" {
				return echo.ErrUnauthorized
			}

			// 截取前面的Bearer
			if err := m.jt.ParseToken(token[7:], nil); err != nil {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
