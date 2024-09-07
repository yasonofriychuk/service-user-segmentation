package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/passionde/user-segmentation-service/internal/service"
	"net/http"
	"os"
	"path"
)

type reportsRoutes struct {
	historyService service.History
}

func newReportsRoutes(g *echo.Group, historyService service.History) {
	r := reportsRoutes{
		historyService: historyService,
	}
	g.GET("/:fileName", r.downloadFile)
}

func (r *reportsRoutes) downloadFile(c echo.Context) error {
	fileName := path.Clean(c.Param("fileName"))
	filePath := fmt.Sprintf("reports/%s", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return c.String(http.StatusNotFound, "File not found")
	}
	defer file.Close()
	return c.Stream(http.StatusOK, "text/csv", file)
}
