package models

import (
	"errors"
	"fmt"
	"time"
)

type Volunteer struct {
	ID         uint
	UserId     uint
	User       User
	Contract   VolunteerContract
	ContractId uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (v *Volunteer) Create() (*Volunteer, error) {

	var err error
	fmt.Println(v)
	err = DB.Create(&v).Error
	if err != nil {
		return &Volunteer{}, err
	}
	return v, nil
}

func GetVolunteerByUserId(userId uint) (Volunteer, error) {

	var volunteer Volunteer

	if err := DB.Where("user_id=?", userId).
		Preload("Contract").
		First(&volunteer).Error; err != nil {
		return volunteer, errors.New("Волентер не найден!")
	}

	return volunteer, nil

}
