package modelrequests

import (
	"github.com/shopspring/decimal"
)

type AddExpenseRequest struct {
	CategoryId int             `json:"categoryId"`
	Expense    string          `json:"expense"`
	Total      decimal.Decimal `json:"total"`
}
