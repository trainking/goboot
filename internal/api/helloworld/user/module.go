package user

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/httpapi"
	"github.com/trainking/goboot/pkg/log"
	"google.golang.org/grpc/metadata"

	userClient "github.com/trainking/goboot/internal/service/user/client"
)

var Module = func() httpapi.Module {
	return new(M)
}

type M struct {
	Config *viper.Viper

	UserService pb.UserServiceClient
}

func (m *M) Init(app *httpapi.App) {
	m.Config = app.Config

	var err error
	m.UserService, err = userClient.NewUserServiceByMap(app.Config.GetStringMap("UserService"))
	if err != nil {
		log.Errorf("NewUserServiceByMap failed: %v", err)
		return
	}
}

func (m *M) Group() httpapi.Group {
	return httpapi.Group{
		Path:        "/api/helloworld",
		Middlewares: []httpapi.Middleware{},
		Routers: []httpapi.Router{
			{
				Method: http.MethodGet,
				Path:   "/userinfo",
				Handle: m.GetUserInfo,
			},
		},
	}
}

func (m *M) GetUserInfo(c echo.Context) error {
	requestID := httpapi.GetRequestID(c)
	fmt.Printf("reqeustID: %v\n", requestID)
	ctx := c.Request().Context()
	// ctx = context.WithValue(ctx, "RequestID", requestID)
	ctx = metadata.AppendToOutgoingContext(ctx, echo.HeaderXRequestID, requestID)
	reply, err := m.UserService.GetUserInfo(ctx, &pb.GetUserInfoArgs{
		UserId: 1,
	})
	if err != nil {
		log.Errorf("GetUserInfo error: %v", err)
	}
	return c.JSON(http.StatusOK, reply.UserName)
}
