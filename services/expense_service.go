package services

import (
	"context"
	"errors"
	"expense-tracker-api/helpers"
	modelentities "expense-tracker-api/models/entities"
	modelrequests "expense-tracker-api/models/requests"
	modelresponses "expense-tracker-api/models/responses"
	"expense-tracker-api/repositories"
	"expense-tracker-api/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type ExpenseService interface {
	Add(ctx context.Context, addExpenseRequest modelrequests.AddExpenseRequest) (httpCode int, response interface{})
	Update(ctx context.Context, updateExpenseRequest modelrequests.UpdateExpenseRequest, id int) (httpCode int, response interface{})
	Delete(ctx context.Context, id int) (httpCode int, response interface{})
	FindByFilter(ctx context.Context, now time.Time, filter string, startDate string, endDate string) (httpCode int, response interface{})
}

type ExpenseServiceImplementation struct {
	PostgresUtil      utils.PostgresUtil
	Validate          *validator.Validate
	ExpenseRepository repositories.ExpenseRepository
}

func NewExpenseService(postgresUtil utils.PostgresUtil, validate *validator.Validate, expenseRepository repositories.ExpenseRepository) ExpenseService {
	return &ExpenseServiceImplementation{
		PostgresUtil:      postgresUtil,
		Validate:          validate,
		ExpenseRepository: expenseRepository,
	}
}

func (service *ExpenseServiceImplementation) Add(ctx context.Context, addExpenseRequest modelrequests.AddExpenseRequest) (httpCode int, response interface{}) {
	err := service.Validate.Struct(addExpenseRequest)
	if err != nil {
		httpCode = http.StatusBadRequest
		response = helpers.ToResponse(err.Error())
		return
	}

	tx, err := service.PostgresUtil.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	defer func() {
		errCommitOrRollback := service.PostgresUtil.CommitOrRollback(tx, ctx, err)
		if errCommitOrRollback != nil {
			httpCode = http.StatusInternalServerError
			response = helpers.ToResponse(errCommitOrRollback.Error())
		}
	}()

	userId, ok := ctx.Value("userId").(int)
	if !ok {
		err = errors.New("cannot find user id")
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	var expense modelentities.Expense
	expense.UserId = pgtype.Int4{Valid: true, Int32: int32(userId)}
	expense.CategoryId = pgtype.Int4{Valid: true, Int32: int32(addExpenseRequest.CategoryId)}
	expense.Expense = pgtype.Text{Valid: true, String: addExpenseRequest.Expense}
	expense.Total = decimal.NullDecimal{Valid: true, Decimal: addExpenseRequest.Total}
	expense.CreatedAt = pgtype.Int8{Valid: true, Int64: time.Now().UnixMilli()}
	_, err = service.ExpenseRepository.Create(tx, ctx, expense)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	httpCode = http.StatusCreated
	response = helpers.ToResponse("successfully add new expense")
	return
}

func (service *ExpenseServiceImplementation) Update(ctx context.Context, updateExpenseRequest modelrequests.UpdateExpenseRequest, id int) (httpCode int, response interface{}) {
	err := service.Validate.Struct(updateExpenseRequest)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	tx, err := service.PostgresUtil.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	defer func() {
		errCommitOrRollback := service.PostgresUtil.CommitOrRollback(tx, ctx, err)
		if errCommitOrRollback != nil {
			httpCode = http.StatusInternalServerError
			response = helpers.ToResponse(errCommitOrRollback.Error())
		}
	}()

	userId, ok := ctx.Value("userId").(int)
	if !ok {
		err = errors.New("cannot find user id")
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}
	expense, err := service.ExpenseRepository.FindByIdAndUserId(tx, ctx, id, userId)
	if err == pgx.ErrNoRows {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	} else if err != nil && err != pgx.ErrNoRows {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	expense.CategoryId = pgtype.Int4{Valid: true, Int32: int32(updateExpenseRequest.CategoryId)}
	expense.Expense = pgtype.Text{Valid: true, String: updateExpenseRequest.Expense}
	expense.Total = decimal.NullDecimal{Valid: true, Decimal: updateExpenseRequest.Total}
	rowsAffected, err := service.ExpenseRepository.Update(tx, ctx, expense)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not one")
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	httpCode = http.StatusOK
	response = helpers.ToResponse("successfully updated expense")
	return
}

func (service *ExpenseServiceImplementation) Delete(ctx context.Context, id int) (httpCode int, response interface{}) {
	tx, err := service.PostgresUtil.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	defer func() {
		errCommitOrRollback := service.PostgresUtil.CommitOrRollback(tx, ctx, err)
		if errCommitOrRollback != nil {
			httpCode = http.StatusInternalServerError
			response = helpers.ToResponse(errCommitOrRollback.Error())
		}
	}()

	rowsAffected, err := service.ExpenseRepository.Delete(tx, ctx, id)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not one")
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	httpCode = http.StatusNoContent
	response = helpers.ToResponse("successfully deleted expense")
	return
}

func (service *ExpenseServiceImplementation) FindByFilter(ctx context.Context, now time.Time, filter string, startDate string, endDate string) (httpCode int, response interface{}) {
	var startDateUnixMilli int64
	var endDateUnixMilli int64
	var expenses []modelentities.Expense
	var err error
	var total decimal.Decimal
	var totalExpense float64
	endDateUnixMilli = now.UnixMilli()
	if filter == "pastWeek" {
		startDateUnixMilli = now.AddDate(0, 0, -1).UnixMilli()
	} else if filter == "pastMonth" {
		startDateUnixMilli = now.AddDate(0, -1, 0).UnixMilli()
	} else if filter == "last3Months" {
		startDateUnixMilli = now.AddDate(0, -3, 0).UnixMilli()
	} else if filter == "custom" {
		layout := "2006-01-02"
		startDateTime, err := time.Parse(layout, startDate)
		if err != nil {
			httpCode = http.StatusInternalServerError
			response = helpers.ToResponse(err.Error())
			return
		}
		startDateUnixMilli = startDateTime.UnixMilli()

		endDateTime, err := time.Parse(layout, endDate)
		if err != nil {
			httpCode = http.StatusInternalServerError
			response = helpers.ToResponse(err.Error())
			return
		}
		endDateUnixMilli = endDateTime.UnixMilli()
	} else {
		err = errors.New("cannot find filter")
		httpCode = http.StatusBadRequest
		response = helpers.ToResponse(err.Error())
		return
	}
	expenses, err = service.ExpenseRepository.FindByStartDateUnixMilliAndEndDateUnixMilli(service.PostgresUtil.GetPool(), ctx, startDateUnixMilli, endDateUnixMilli)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	for _, expense := range expenses {
		total = total.Add(expense.Total.Decimal)
	}
	totalExpense, ok := total.Float64()
	if !ok {
		err = errors.New("cannot convert decimal to float64")
		httpCode = http.StatusInternalServerError
		response = helpers.ToResponse(err.Error())
		return
	}

	var filterExpenseResponse modelresponses.FilterExpenseResponse
	filterExpenseResponse.Filter = filter
	filterExpenseResponse.Total = totalExpense
	httpCode = http.StatusOK
	response = filterExpenseResponse
	return
}
