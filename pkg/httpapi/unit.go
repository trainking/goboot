package httapi

import "github.com/labstack/echo/v4"

// BindValidate 绑定并且校验
func BindValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	err := c.Validate(req)
	return err
}
