package controllers

import (
	modelrequests "expense-tracker-api/models/requests"
	"expense-tracker-api/services"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ExpenseController interface {
	Add(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	FindByFilter(c echo.Context) error
}

type ExpenseControllerImplementation struct {
	ExpenseService services.ExpenseService
}

func NewExpenseController(expenseService services.ExpenseService) ExpenseController {
	return &ExpenseControllerImplementation{
		ExpenseService: expenseService,
	}
}

func (controller *ExpenseControllerImplementation) Add(c echo.Context) error {
	var addExpenseRequest modelrequests.AddExpenseRequest
	err := c.Bind(&addExpenseRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.ExpenseService.Add(c.Request().Context(), addExpenseRequest)
	return c.JSON(httpCode, response)
}

func (controller *ExpenseControllerImplementation) Update(c echo.Context) error {
	var updateExpenseRequest modelrequests.UpdateExpenseRequest
	err := c.Bind(&updateExpenseRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.ExpenseService.Update(c.Request().Context(), updateExpenseRequest, id)
	return c.JSON(httpCode, response)
}

func (controller *ExpenseControllerImplementation) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.ExpenseService.Delete(c.Request().Context(), id)
	return c.JSON(httpCode, response)
}

func (controller *ExpenseControllerImplementation) FindByFilter(c echo.Context) error {
	filter := c.QueryParam("filter")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	now := time.Now()
	httpCode, response := controller.ExpenseService.FindByFilter(c.Request().Context(), now, filter, startDate, endDate)
	return c.JSON(httpCode, response)
}
