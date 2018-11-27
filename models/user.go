package models

import (
	"github.com/shopspring/decimal"
)

type User struct {
	UID    int
	GID    int
	TCC    decimal.Decimal
	ETH    decimal.Decimal
	NASH   decimal.Decimal
	Desc   string
	Worker map[int]bool
	I1     int
}

func (u *User) GetID() int {
	return u.UID
}

func (u *User) Index() [][]string {
	return nil
}
