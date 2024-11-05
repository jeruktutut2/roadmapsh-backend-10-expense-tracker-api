package modelentities

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type Expense struct {
	Id         pgtype.Int4
	UserId     pgtype.Int4
	CategoryId pgtype.Int4
	Expense    pgtype.Text
	Total      decimal.NullDecimal
	CreatedAt  pgtype.Int8
}
