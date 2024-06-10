package models

import (
	"fmt"
	"time"
)

type Violator struct {
	ID        uint
	UserId    *uint
	User      User
	Iin       string
	Fio       string
	Address   string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (v *Violator) Create() (*Violator, error) {

	var err error

	fmt.Println(v)
	err = DB.Create(&v).Error
	if err != nil {
		return &Violator{}, err
	}
	return v, nil
}
