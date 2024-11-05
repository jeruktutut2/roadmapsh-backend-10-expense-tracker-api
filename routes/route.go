package routes

import (
	"expense-tracker-api/controllers"
	"expense-tracker-api/middlewares"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo, controller controllers.UserController) {
	e.POST("/register", controller.Register)
	e.POST("/login", controller.Login)
}

func ExpenseRoute(e *echo.Echo, controller controllers.ExpenseController) {
	e.POST("/expense", controller.Add, middlewares.Authenticate)
	e.PUT("/expense/:id", controller.Update, middlewares.Authenticate)
	e.DELETE("/expense/:id", controller.Delete, middlewares.Authenticate)
	e.GET("/expense", controller.FindByFilter, middlewares.Authenticate)
}
