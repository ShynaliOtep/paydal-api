package models

import (
	"errors"
	"fmt"
	"time"
)

type Wallet struct {
	ID           uint
	UserId       uint    `json:"user_id"`
	Balance      float64 `json:"balance"`
	BonusBalance float64 `json:"bonus_balance"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (w *Wallet) Create() (*Wallet, error) {

	var err error
	fmt.Println(w)
	err = DB.Create(&w).Error
	if err != nil {
		return &Wallet{}, err
	}
	return w, nil
}

func (w *Wallet) Save() (*Wallet, error) {

	var err error
	fmt.Println(w)
	err = DB.Save(&w).Error
	if err != nil {
		return &Wallet{}, err
	}
	return w, nil
}

func GetWalletByUserId(userId uint) (Wallet, error) {

	var wallet Wallet

	if err := DB.Where("user_id=?", userId).First(&wallet).Error; err != nil {
		return wallet, errors.New("Wallet not found!")
	}

	return wallet, nil

}
