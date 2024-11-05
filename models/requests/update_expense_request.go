package modelrequests

import "github.com/shopspring/decimal"

type UpdateExpenseRequest struct {
	CategoryId int             `json:"categoryId"`
	Expense    string          `json:"expense"`
	Total      decimal.Decimal `json:"total"`
}
