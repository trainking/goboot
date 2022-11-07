package httpapi

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/trainking/goboot/pkg/jwt"
)

// MiddlewareBearerJwt 使用Bearer Token的一个中间件
type MiddlewareBearerJwt struct {
	JT *jwt.JwtToken // jwt实现
}

// NewMiddlewareBearerJwt 使用Bearer Token的中间件生成器
func NewMiddlewareBearerJwt(jwtSecert string) Middleware {
	return &MiddlewareBearerJwt{
		JT: jwt.NewJwtToken(jwtSecert),
	}
}

// GetHeaderToken 从请求头获取Token
func (m *MiddlewareBearerJwt) GetHeaderToken(c echo.Context) (string, error) {
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	if token == "" {
		return "", echo.ErrUnauthorized
	}

	if len(token) <= 7 {
		return "", echo.ErrUnauthorized
	}

	if strings.ToLower(token[0:6]) != "bearer" {
		return "", echo.ErrUnauthorized
	}

	return token[7:], nil
}

// MiddlewareFunc 对Middleware的实现
func (m *MiddlewareBearerJwt) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := m.GetHeaderToken(c)
			if err != nil {
				return err
			}

			// 截取前面的Bearer
			if err := m.JT.ParseToken(token, nil); err != nil {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
