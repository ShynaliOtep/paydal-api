package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID          uint
	WalletId    uint    `json:"user_id"`
	Amount      float64 `json:"amount"`
	Type        string
	BalanceType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Transaction) Create() (*Transaction, error) {

	var err error
	fmt.Println(t.WalletId)
	err = DB.Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}
