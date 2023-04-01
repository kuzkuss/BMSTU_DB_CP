package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"

	serviceUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/service/usecase"
)

type Delivery struct {
	ServiceUC serviceUsecase.UseCaseI
}

func (delivery *Delivery) ClearData(c echo.Context) error {
	err := delivery.ServiceUC.ClearData()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (delivery *Delivery) SelectStatus(c echo.Context) error {
	status, err := delivery.ServiceUC.SelectStatus()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func NewDelivery(e *echo.Echo, serviceUC serviceUsecase.UseCaseI) {
	handler := &Delivery{
		ServiceUC: serviceUC,
	}

	e.POST("/api/service/clear", handler.ClearData)
	e.GET("/api/service/status", handler.SelectStatus)
}
