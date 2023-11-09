package httpapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ResponseCode 定义响应码应该具有的格式
type ResponseCode interface {
	Code() int32
	Msg() string
}

// ResponseJSON 响应JSON
type ResponseJSON struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var successCode ResponseCode

// InitSuccessCode 初始化成功响应码
func InitSuccessCode(resCode ResponseCode) {
	successCode = resCode
}

// SuccessJson 成功响应JSON
func SuccessJson(c echo.Context, data interface{}) error {
	if successCode == nil {
		return c.JSON(http.StatusOK, ResponseJSON{
			Code: 0,
			Msg:  "ok",
			Data: data,
		})
	}

	return c.JSON(http.StatusOK, ResponseJSON{
		Code: successCode.Code(),
		Msg:  successCode.Msg(),
		Data: data,
	})
}

// ErrorJson 错误返回JSON
func ErrorJson(c echo.Context, resCode ResponseCode) error {
	return c.JSON(http.StatusOK, ResponseJSON{
		Code: resCode.Code(),
		Msg:  resCode.Msg(),
		Data: nil,
	})
}
