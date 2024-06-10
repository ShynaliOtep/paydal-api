package models

import (
	"fmt"
	"time"
)

type Location struct {
	ID        uint      `json:"id"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (l *Location) Create() (*Location, error) {

	var err error
	fmt.Println(l)
	err = DB.Create(&l).Error
	if err != nil {
		return &Location{}, err
	}
	return l, nil
}
