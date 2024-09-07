package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/passionde/user-segmentation-service/docs"
	"github.com/passionde/user-segmentation-service/internal/service"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	reportsGroup := handler.Group("/reports")
	{
		newReportsRoutes(reportsGroup, services.History)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newUserRoutes(v1.Group("/users"), services.User)
		newSegmentRoutes(v1.Group("/segments"), services.Segment)
		newHistoryRoutes(v1.Group("/history"), services.History)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
