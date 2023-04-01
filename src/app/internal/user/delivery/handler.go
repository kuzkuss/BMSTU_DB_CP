package delivery

import (
	"net/http"

	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	userUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/user/usecase"
)

type Delivery struct {
	UserUC userUsecase.UseCaseI
}

func (delivery *Delivery) CreateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.NickName = c.Param("nickname")

	conflictUsers, err := delivery.UserUC.CreateUser(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConflict):
			c.Logger().Error(err)
			return c.JSON(http.StatusConflict, conflictUsers)
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) SelectUser(c echo.Context) error {
	user, err := delivery.UserUC.SelectUser(c.Param("nickname"))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) UpdateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user.NickName = c.Param("nickname")

	err = delivery.UserUC.UpdateUser(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrConflict):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusConflict, models.ErrConflict.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, user)
}

func NewDelivery(e *echo.Echo, userUC userUsecase.UseCaseI) {
	handler := &Delivery{
		UserUC: userUC,
	}

	e.POST("/api/user/:nickname/create", handler.CreateUser)
	e.GET("/api/user/:nickname/profile", handler.SelectUser)
	e.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

