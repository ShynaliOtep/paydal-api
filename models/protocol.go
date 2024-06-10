package models

import (
	"errors"
	"fmt"
	"time"
)

type Protocol struct {
	ID          uint
	PoliceId    uint
	Police      User
	ViolationId uint
	Violation   Violation
	ViolatorId  uint
	Violator    Violator
	CarId       *uint
	Car         Car
	ArticleId   uint
	Article     Article
	Status      string
	Mrp         int
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Protocol) Create() (*Protocol, error) {

	var err error

	fmt.Println(p)
	err = DB.Create(&p).Error
	if err != nil {
		return &Protocol{}, err
	}
	return p, nil
}

func GetProtocolsByIin(iin string) ([]Protocol, error) {

	var protocols []Protocol
	if err := DB.
		Preload("Car").
		Preload("Violation").
		Preload("Article").
		Preload("Violation.Medias").
		Preload("Violation.Category").
		Preload("Violation.Location").
		Preload("Violation.City").
		Joins("JOIN violators ON violator_id = violators.id").
		Where("violators.iin = ?", iin).
		Find(&protocols).Error; err != nil {
		return protocols, errors.New("Штраф не найден!")
	}

	return protocols, nil

}

func GetPoliceProtocols(policeId uint) ([]Protocol, error) {
	var protocols []Protocol
	if err := DB.
		Preload("Car").
		Preload("Violator").
		Preload("Violation").
		Preload("Article").
		Preload("Violation.Medias").
		Preload("Violation.Category").
		Preload("Violation.Location").
		Preload("Violation.City").
		Where("police_id = ?", policeId).
		Find(&protocols).Error; err != nil {
		return protocols, errors.New("Штраф не найден!")
	}

	return protocols, nil
}
