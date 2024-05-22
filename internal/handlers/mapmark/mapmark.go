package mapmark

import (
	"context"
	"echoFramework/internal/domain"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type MarkService interface {
	Get(ctx context.Context) ([]domain.Mark, error)
	Save(ctx context.Context, mark domain.Mark) (int64, error)
	Update(ctx context.Context, mark domain.Mark) error
	GetById(ctx context.Context, id int) (domain.Mark, error)
}

type MarkRouter struct {
	ms MarkService
}

func SetRouter(ms MarkService, e *echo.Group) *MarkRouter {

	r := &MarkRouter{
		ms: ms,
	}

	e.GET("", r.Get)
	e.POST("", r.Post)
	e.GET("/:id", r.GetById)
	e.PUT("/:id", r.Put)

	return r
}

func (m *MarkRouter) Get(c echo.Context) error {
	marks, err := m.ms.Get(context.Background())
	if err != nil {
		err = c.JSON(http.StatusInternalServerError, err)
		return err
	}

	err = c.JSON(http.StatusOK, marks)
	return err
}

func (m *MarkRouter) Post(c echo.Context) error {

	mark := domain.Mark{}
	err := c.Bind(&mark)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := m.ms.Save(context.Background(), mark)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidParams) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = c.JSON(http.StatusInternalServerError, err)
		return err
	}

	err = c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	return err
}

func (m *MarkRouter) GetById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	mark, err := m.ms.GetById(context.Background(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, mark)
}

func (m *MarkRouter) Put(c echo.Context) error {
	mark := domain.Mark{}
	err := c.Bind(&mark)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	mark.Id = id

	if mark.Lng == nil || mark.Lat == nil {
		return echo.NewHTTPError(http.StatusBadRequest, domain.ErrInvalidParams)
	}

	err = m.ms.Update(context.Background(), mark)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, mark)
}
