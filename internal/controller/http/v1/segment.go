package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/passionde/user-segmentation-service/internal/service"
	"net/http"
)

type segmentRoutes struct {
	segmentService service.Segment
}

func newSegmentRoutes(g *echo.Group, segmentService service.Segment) {
	r := segmentRoutes{
		segmentService: segmentService,
	}
	g.POST("/create", r.create)
	g.DELETE("/delete", r.delete)
}

type createSegmentInput struct {
	Slug            string `json:"slug" validate:"required,max=256"`
	PercentageUsers int    `json:"percentageUsers" validate:"omitempty,min=1,max=10000"`
}

type createSegmentResponse struct {
	Slug string `json:"slug"`
}

// @Summary Создание сегмента
// @Description Этот эндпоинт позволяет создать новый сегмент
// @Tags Segments
// @ID createSegment
// @Accept json
// @Produce json
// @Param Authorization header string true "API KEY для аутентификации"
// @Param input body createSegmentInput true "Данные для создания сегмента"
// @Success 201 {object} createSegmentResponse "Успешное выполнение"
// @Failure 400 {object} echo.HTTPError "Некорректный запрос или данные"
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервера"
// @Router /api/v1/segments/create [post]
func (s *segmentRoutes) create(c echo.Context) error {
	var input createSegmentInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := s.segmentService.CreateSegment(c.Request().Context(), service.CreateSegmentInput{
		Slug:            input.Slug,
		PercentageUsers: input.PercentageUsers,
	})

	if err != nil {
		if errors.Is(err, service.ErrSegmentAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusCreated, createSegmentResponse{
		Slug: input.Slug,
	})
}

type deleteSegmentInput struct {
	Slug string `json:"slug" validate:"required,max=256"`
}

// @Summary Удаление сегмента
// @Description Этот эндпоинт позволяет удалить существующий сегмент.
// @Tags Segments
// @ID deleteSegment
// @Accept json
// @Produce json
// @Param Authorization header string true "API KEY для аутентификации"
// @Param input body deleteSegmentInput true "Данные для удаления сегмента"
// @Success 204 "Успешное удаление"
// @Failure 400 {object} echo.HTTPError "Некорректный запрос или данные"
// @Failure 404 {object} echo.HTTPError "Сегмент не найден"
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервера"
// @Router /api/v1/segments/delete [delete]
func (s *segmentRoutes) delete(c echo.Context) error {
	var input deleteSegmentInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := s.segmentService.DeleteSegment(c.Request().Context(), service.SegmentInput{
		Slug: input.Slug,
	})
	if err != nil {
		if errors.Is(err, service.ErrSegmentNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	return c.NoContent(204)
}
