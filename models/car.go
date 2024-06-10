package models

import (
	"fmt"
	"time"
)

type Car struct {
	ID        uint
	Number    string
	Passport  string
	Mark      string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Car) Create() (*Car, error) {

	var err error

	fmt.Println(c)
	err = DB.Create(&c).Error
	if err != nil {
		return &Car{}, err
	}
	return c, nil
}
