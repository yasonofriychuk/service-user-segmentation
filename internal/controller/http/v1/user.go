package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/passionde/user-segmentation-service/internal/service"
	"net/http"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(g *echo.Group, userService service.User) {
	r := &userRoutes{
		userService: userService,
	}
	g.POST("/segments", r.setSegments)
	g.GET("/active-segments", r.getSegments)
}

type setSegmentsUserInput struct {
	UserID      string   `json:"user_id" validate:"required,max=40"`
	SegmentsAdd []string `json:"segments_add" validate:"required"`
	SegmentsDel []string `json:"segments_del" validate:"required"`
	TTL         uint64   `json:"ttl" validate:"omitempty,min=1,max=18446744073709551615"`
}

// @Summary Обновление сегментов пользователя
// @Description Этот эндпоинт позволяет обновить сегменты, к которым принадлежит пользователь.
// @Tags Users
// @ID setSegments
// @Accept json
// @Produce json
// @Param Authorization header string true "API KEY для аутентификации"
// @Param input body setSegmentsUserInput true "Данные для обновления сегментов пользователя"
// @Success 200 "Успешная операция"
// @Failure 400 {object} echo.HTTPError "Некорректный запрос или данные"
// @Failure 404 {object} echo.HTTPError "Сегмент не найден"
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервера"
// @Router /api/v1/users/segments [post]
func (u *userRoutes) setSegments(c echo.Context) error {
	var input setSegmentsUserInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	err := u.userService.SetSegments(c.Request().Context(), service.SetSegmentsUserInput{
		UserID:      input.UserID,
		SegmentsAdd: input.SegmentsAdd,
		SegmentsDel: input.SegmentsDel,
		TTL:         input.TTL,
	})
	if err != nil {
		if errors.Is(err, service.ErrSegmentNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	return c.NoContent(200)
}

type getSegmentsUserInput struct {
	UserID string `json:"user_id" validate:"required,max=40"`
}

type getSegmentsUserResponse struct {
	UserID   string   `json:"user_id"`
	Segments []string `json:"segments"`
}

// @Summary Получение активных сегментов пользователя
// @Description Этот эндпоинт позволяет получить список сегментов, к которым принадлежит пользователь.
// @Tags Users
// @ID getSegments
// @Accept json
// @Produce json
// @Param Authorization header string true "API KEY для аутентификации"
// @Param user_id query string true "Идентификатор пользователя"
// @Success 200 {object} getSegmentsUserResponse "Успешное выполнение"
// @Failure 400 {object} echo.HTTPError "Некорректный запрос или данные"
// @Failure 404 {object} echo.HTTPError "Пользователь не найден"
// @Failure 500 {object} echo.HTTPError "Внутренняя ошибка сервера"
// @Router /api/v1/users/active-segments [get]
func (u *userRoutes) getSegments(c echo.Context) error {
	input := getSegmentsUserInput{
		UserID: c.QueryParams().Get("user_id"),
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	segments, err := u.userService.GetSegments(
		c.Request().Context(),
		service.GetSegmentsUserInput{UserID: input.UserID},
	)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return err
	}

	return c.JSON(http.StatusOK, getSegmentsUserResponse{
		UserID:   input.UserID,
		Segments: segments,
	})
}
