package models

import (
	"errors"
	"fmt"
	"time"
)

type VolunteerContract struct {
	ID         uint
	UserId     uint
	User       User
	Fio        string
	Document   VolunteerContractDocument
	DocumentId uint
	CityId     uint
	City       City
	Status     string
	Reason     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

const status_wait = "waiting"

type VolunteerContractDocument struct {
	ID                uint
	BackImage         File
	BackImageId       uint
	FrontImageId      uint
	FrontImage        File
	BackWithAvatar    File
	BackWithAvatarId  uint
	FrontWithAvatarId uint
	FrontWithAvatar   File
}

func (c *VolunteerContract) Create() (*VolunteerContract, error) {

	var err error
	if c.Status == "" {
		fmt.Println("status nil")
		c.Status = status_wait
	}
	fmt.Println(c)
	err = DB.Create(&c).Error
	if err != nil {
		return &VolunteerContract{}, err
	}
	return c, nil
}

func (c *VolunteerContract) Save() (*VolunteerContract, error) {

	var err error
	fmt.Println(c)
	err = DB.Save(&c).Error
	if err != nil {
		return &VolunteerContract{}, err
	}
	return c, nil
}

func (c *VolunteerContractDocument) Save() (*VolunteerContractDocument, error) {

	var err error
	fmt.Println(c)
	err = DB.Create(&c).Error
	if err != nil {
		return &VolunteerContractDocument{}, err
	}
	return c, nil
}

func GetVolunteerContractByUserId(userId uint) (VolunteerContract, error) {

	var contract VolunteerContract

	if err := DB.Where("user_id=?", userId).
		Preload("Document.FrontImage").
		Preload("Document.BackImage").
		Preload("Document.FrontWithAvatar").
		Preload("Document.BackWithAvatar").
		First(&contract).Error; err != nil {
		return contract, errors.New("Contract not found!")
	}

	return contract, nil

}

func GetVolunteerContractById(id uint) (VolunteerContract, error) {

	var contract VolunteerContract

	if err := DB.
		Preload("Document.FrontImage").
		Preload("Document.BackImage").
		Preload("Document.FrontWithAvatar").
		Preload("Document.BackWithAvatar").
		First(&contract, id).Error; err != nil {
		return contract, errors.New("Contract not found!")
	}

	return contract, nil

}
