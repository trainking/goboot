package httpapi

import "github.com/labstack/echo/v4"

// BindValidate 绑定并且校验
func BindValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	err := c.Validate(req)
	return err
}

// GetRequestID 获取RequestID
func GetRequestID(c echo.Context) string {
	res := c.Response()
	return res.Header().Get(echo.HeaderXRequestID)
}
