package repositories

import (
	"context"
	modelentities "expense-tracker-api/models/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepository interface {
	Create(tx pgx.Tx, ctx context.Context, expense modelentities.Expense) (lasInsertedId int, err error)
	FindByIdAndUserId(tx pgx.Tx, ctx context.Context, id int, userId int) (expense modelentities.Expense, err error)
	Update(tx pgx.Tx, ctx context.Context, expense modelentities.Expense) (rowsAffected int64, err error)
	Delete(tx pgx.Tx, ctx context.Context, id int) (rowsAffected int64, err error)
	FindByStartDateUnixMilliAndEndDateUnixMilli(pool *pgxpool.Pool, ctx context.Context, startDate int64, endDate int64) (expenses []modelentities.Expense, err error)
}

type ExpenseRepositoryImplementation struct {
}

func NewExpenseRepository() ExpenseRepository {
	return &ExpenseRepositoryImplementation{}
}

func (repository *ExpenseRepositoryImplementation) Create(tx pgx.Tx, ctx context.Context, expense modelentities.Expense) (lasInsertedId int, err error) {
	query := `INSERT INTO expenses(user_id,category_id,expense,total,created_at) VALUES($1,$2,$3,$4,$5) RETURNING id;`
	err = tx.QueryRow(ctx, query, expense.UserId, expense.CategoryId, expense.Expense, expense.Total, expense.CreatedAt).Scan(&lasInsertedId)
	return
}

func (repository *ExpenseRepositoryImplementation) FindByIdAndUserId(tx pgx.Tx, ctx context.Context, id int, userId int) (expense modelentities.Expense, err error) {
	query := `SELECT id,user_id,category_id,expense,total,created_at FROM expenses WHERE id = $1 AND user_id = $2;`
	err = tx.QueryRow(ctx, query, id, userId).Scan(&expense.Id, &expense.UserId, &expense.CategoryId, &expense.Expense, &expense.Total, &expense.CreatedAt)
	return
}

func (repository *ExpenseRepositoryImplementation) Update(tx pgx.Tx, ctx context.Context, expense modelentities.Expense) (rowsAffected int64, err error) {
	query := `UPDATE expenses SET category_id = $1, expense = $2, total = $3 WHERE id = $4 AND user_id = $5;`
	result, err := tx.Exec(ctx, query, expense.CategoryId, expense.Expense, expense.Total, expense.Id, expense.UserId)
	if err != nil {
		return
	}
	rowsAffected = result.RowsAffected()
	return
}

func (repository *ExpenseRepositoryImplementation) Delete(tx pgx.Tx, ctx context.Context, id int) (rowsAffected int64, err error) {
	query := `DELETE FROM expenses WHERE id = $1;`
	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return
	}
	rowsAffected = result.RowsAffected()
	return
}

func (repository *ExpenseRepositoryImplementation) FindByStartDateUnixMilliAndEndDateUnixMilli(pool *pgxpool.Pool, ctx context.Context, startDate int64, endDate int64) (expenses []modelentities.Expense, err error) {
	query := `SELECT id,user_id,category_id,expense,total,created_at FROM expenses WHERE created_at >= $1 AND created_at <= $2;`
	rows, err := pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var expense modelentities.Expense
		err = rows.Scan(&expense.Id, &expense.UserId, &expense.CategoryId, &expense.Expense, &expense.Total, &expense.CreatedAt)
		if err != nil {
			expenses = []modelentities.Expense{}
			return
		}
		expenses = append(expenses, expense)
	}

	if rows.Err() != nil {
		expenses = []modelentities.Expense{}
		return
	}
	return
}
