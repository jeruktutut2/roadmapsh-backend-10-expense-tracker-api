package main

import (
	"context"
	"expense-tracker-api/controllers"
	"expense-tracker-api/helpers"
	"expense-tracker-api/repositories"
	"expense-tracker-api/routes"
	"expense-tracker-api/services"
	"expense-tracker-api/utils"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	postgresUtil := utils.NewPostgresConnection()
	e := echo.New()
	validate := validator.New()
	bcryptHelper := helpers.NewBcryptHelper()
	jwtHelper := helpers.NewJwtHelper()

	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(postgresUtil, validate, userRepository, bcryptHelper, jwtHelper)
	userController := controllers.NewUserController(userService)
	routes.UserRoute(e, userController)

	expenseRepository := repositories.NewExpenseRepository()
	expenseService := services.NewExpenseService(postgresUtil, validate, expenseRepository)
	expenseController := controllers.NewExpenseController(expenseService)
	routes.ExpenseRoute(e, expenseController)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(os.Getenv("ECHO_HOST")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
