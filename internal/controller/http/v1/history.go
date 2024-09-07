package v1

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/passionde/user-segmentation-service/internal/service"
	"net/http"
)

type historyRoutes struct {
	historyService service.History
}

func newHistoryRoutes(g *echo.Group, historyService service.History) {
	r := historyRoutes{
		historyService: historyService,
	}
	g.POST("/report-link", r.reportLink) // POST as the file is created when called
}

type getHistoryInput struct {
	UserID string `json:"user_id" validate:"required,max=40"`
	Year   int    `json:"year" validate:"required"`
	Month  int    `json:"month" validate:"required"`
}

type getHistoryResponse struct {
	UserID     string `json:"user_id"`
	ReportLink string `json:"report_link"`
}

// @Summary Получение ссылки на CSV отчет
// @Description Этот эндпоинт позволяет получить ссылку на CSV отчет по пользователю за определенный месяц и год.
// @Tags History
// @ID getReportLink
// @Accept json
// @Produce json
// @Param Authorization header string true "API KEY для аутентификации"
// @Param input body getHistoryInput true "Данные для получения ссылки на отчет"
// @Success 200 {object} getHistoryResponse "Успешное выполнение"
// @Failure 400 {object} echo.HTTPError "Некорректный запрос или у пользователя отсутствует история за указанный период"
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервера"
// @Router /api/v1/history/report-link [post]
func (h *historyRoutes) reportLink(c echo.Context) error {
	var input getHistoryInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	filename, err := h.historyService.GetNotes(c.Request().Context(), service.GetHistoryInput{
		UserID: input.UserID,
		Year:   input.Year,
		Month:  input.Month,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNoData) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return err
	}

	return c.JSON(http.StatusOK, getHistoryResponse{
		UserID:     input.UserID,
		ReportLink: fmt.Sprintf("http://%s/reports/%s", c.Request().Host, filename),
	})
}
